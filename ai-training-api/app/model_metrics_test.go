package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"

	"github.com/go-kit/log"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/grafana/ai-training-o11y/ai-training-api/model"
)

type testApp struct {
	App
}

func (a *testApp) db(ctx context.Context) *gorm.DB {
	return a.App._db
}

func setupTestDB(t *testing.T) (*gorm.DB, func()) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&model.Process{}, &model.ModelMetrics{})
	require.NoError(t, err)

	return db, func() {
		sqlDB, err := db.DB()
		require.NoError(t, err)
		sqlDB.Close()
	}
}

func TestExtractAndValidateProcessID(t *testing.T) {
    tests := []struct {
        name           string
        url            string
        expectedID     uuid.UUID
        expectedErrMsg string
    }{
        {
            name:       "Valid UUID",
            url:        "/process/123e4567-e89b-12d3-a456-426614174000/model-metrics",
            expectedID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
        },
        {
            name:           "Invalid UUID",
            url:            "/process/invalid-uuid/model-metrics",
            expectedErrMsg: "invalid process ID",
        },
        {
            name:           "Empty ID",
            url:            "/process//model-metrics",
            expectedErrMsg: "process ID is empty",
        },
        {
            name:           "No ID in URL",
            url:            "/process/model-metrics",
            expectedErrMsg: "process ID not provided in URL",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            router := mux.NewRouter()
            router.HandleFunc("/process/{id}/model-metrics", func(w http.ResponseWriter, r *http.Request) {
                processID, err := extractAndValidateProcessID(r)
                
                if tt.expectedErrMsg != "" {
                    assert.Error(t, err)
                    assert.Contains(t, err.Error(), tt.expectedErrMsg)
                } else {
                    assert.NoError(t, err)
                    assert.Equal(t, tt.expectedID, processID)
                }
            }).Methods("POST")

            req, err := http.NewRequest("POST", tt.url, nil)
            if err != nil {
                t.Fatalf("Failed to create request: %v", err)
            }

            rr := httptest.NewRecorder()
            router.ServeHTTP(rr, req)
        })
    }
}

func TestValidateProcessExists(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	if db == nil {
		t.Fatal("setupTestDB returned a nil database")
	}

	app := &testApp{
		App: App{
			_db:   db,
			dbMux: &sync.Mutex{},
			logger: log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr)),
		},
	}

	tests := []struct {
		name           string
		setupDB        func(*gorm.DB)
		expectedErrMsg string
	}{
		{
			name: "Process exists",
			setupDB: func(db *gorm.DB) {
				process := model.Process{ID: uuid.New()}
				result := db.Create(&process)
				if result.Error != nil {
					t.Fatalf("Failed to create process: %v", result.Error)
				}
			},
		},
		{
			name:           "Process does not exist",
			setupDB:        func(db *gorm.DB) {},
			expectedErrMsg: "process not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupDB(db)

			// Use a fixed UUID for testing to ensure we're looking for the correct process
			testUUID := uuid.New()
			if tt.name == "Process exists" {
				process := model.Process{ID: testUUID}
				result := db.Create(&process)
				if result.Error != nil {
					t.Fatalf("Failed to create process: %v", result.Error)
				}
			}

			err := app.validateProcessExists(context.Background(), testUUID)

			if tt.expectedErrMsg != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErrMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestParseAndValidateModelMetricsRequest(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		expectedLen    int
		expectedErrMsg string
	}{
		{
			name: "Valid request",
			requestBody: []modelMetricsRequest{
				{
					MetricName: "accuracy",
					StepName:   "training",
					Points: []struct {
						Step  uint32 `json:"step"`
						Value string `json:"value"`
					}{
						{Step: 1, Value: "0.75"},
						{Step: 2, Value: "0.85"},
					},
				},
			},
			expectedLen: 1,
		},
		{
			name: "Invalid metric name",
			requestBody: []modelMetricsRequest{
				{
					MetricName: "",
					StepName:   "training",
					Points: []struct {
						Step  uint32 `json:"step"`
						Value string `json:"value"`
					}{{Step: 1, Value: "0.75"}},
				},
			},
			expectedErrMsg: "metric name must be between 1 and 32 characters",
		},
		{
			name: "Invalid step name",
			requestBody: []modelMetricsRequest{
				{
					MetricName: "accuracy",
					StepName:   "",
					Points: []struct {
						Step  uint32 `json:"step"`
						Value string `json:"value"`
					}{{Step: 1, Value: "0.75"}},
				},
			},
			expectedErrMsg: "step name must be between 1 and 32 characters",
		},
		{
			name: "Invalid step value",
			requestBody: []modelMetricsRequest{
				{
					MetricName: "accuracy",
					StepName:   "training",
					Points: []struct {
						Step  uint32 `json:"step"`
						Value string `json:"value"`
					}{{Step: 0, Value: "0.75"}},
				},
			},
			expectedErrMsg: "step must be a positive number",
		},
		{
			name: "Invalid metric value",
			requestBody: []modelMetricsRequest{
				{
					MetricName: "accuracy",
					StepName:   "training",
					Points: []struct {
						Step  uint32 `json:"step"`
						Value string `json:"value"`
					}{{Step: 1, Value: ""}},
				},
			},
			expectedErrMsg: "metric value must be between 1 and 64 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req, err := http.NewRequest("POST", "/process/123/model-metrics", bytes.NewBuffer(body))
			require.NoError(t, err)

			result, err := parseAndValidateModelMetricsRequest(req)

			if tt.expectedErrMsg != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErrMsg)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.expectedLen)
			}
		})
	}
}
