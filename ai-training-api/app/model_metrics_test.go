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
			requestBody: []ModelMetricsSeries{
				{
					MetricName: "accuracy",
					StepName:   "training",
					Points: []struct {
						Step  uint32 `json:"step"`
						Value json.Number `json:"value"`
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
			requestBody: []ModelMetricsSeries{
				{
					MetricName: "",
					StepName:   "training",
					Points: []struct {
						Step  uint32 `json:"step"`
						Value json.Number `json:"value"`
					}{{Step: 1, Value: "0.75"}},
				},
			},
			expectedErrMsg: "metric name must be between 1 and 32 characters",
		},
		{
			name: "Invalid step name",
			requestBody: []ModelMetricsSeries{
				{
					MetricName: "accuracy",
					StepName:   "",
					Points: []struct {
						Step  uint32 `json:"step"`
						Value json.Number `json:"value"`
					}{{Step: 1, Value: "0.75"}},
				},
			},
			expectedErrMsg: "step name must be between 1 and 32 characters",
		},
		{
			name: "Invalid step value",
			requestBody: []ModelMetricsSeries{
				{
					MetricName: "accuracy",
					StepName:   "training",
					Points: []struct {
						Step  uint32 `json:"step"`
						Value json.Number `json:"value"`
					}{{Step: 0, Value: "0.75"}},
				},
			},
			expectedErrMsg: "step must be a positive number",
		},
		{
			name: "Invalid metric value (empty string)",
			requestBody: []interface{}{
				map[string]interface{}{
					"metric_name": "accuracy",
					"step_name":   "training",
					"points": []interface{}{
						map[string]interface{}{
							"step":  1,
							"value": "",
						},
					},
				},
			},
			expectedErrMsg: "invalid JSON: json: invalid number literal, trying to unmarshal \"\\\"\\\"\" into Number",
		},
		{
			name: "Invalid metric value (not a number)",
			requestBody: []interface{}{
				map[string]interface{}{
					"metric_name": "accuracy",
					"step_name":   "training",
					"points": []interface{}{
						map[string]interface{}{
							"step":  1,
							"value": "not a number",
						},
					},
				},
			},
			expectedErrMsg: "invalid JSON: json: invalid number literal, trying to unmarshal \"\\\"not a number\\\"\" into Number",
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
				if err == nil {
					t.Errorf("Expected error containing '%s', but got nil error", tt.expectedErrMsg)
				} else {
					assert.Contains(t, err.Error(), tt.expectedErrMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.expectedLen)
			}
	
			// Add this line for debugging
			t.Logf("Test case '%s': error = %v, result = %+v", tt.name, err, result)
		})
	}
}

func TestGetModelMetrics(t *testing.T) {
    db, cleanup := setupTestDB(t)
    defer cleanup()

    app := &testApp{
        App: App{_db: db},
    }

    type testCase struct {
        name    string
        metrics []model.ModelMetrics
        check   func(*testing.T, GetModelMetricsResponse)
    }

    testCases := []testCase{
        {
            name: "Basic case",
            metrics: []model.ModelMetrics{
                {MetricName: "accuracy", StepName: "train", Step: 1, MetricValue: "0.75"},
                {MetricName: "accuracy", StepName: "train", Step: 2, MetricValue: "0.80"},
                {MetricName: "loss", StepName: "train", Step: 1, MetricValue: "0.5"},
                {MetricName: "loss", StepName: "train", Step: 2, MetricValue: "0.4"},
            },
            check: func(t *testing.T, response GetModelMetricsResponse) {
                require.Len(t, response, 2) // Two DataFrameWrappers: one for accuracy, one for loss
                
                // Check accuracy metrics
                require.Equal(t, "accuracy", response[0].MetricName)
                require.Equal(t, "train", response[0].StepName)
                require.Len(t, response[0].Fields, 2)
                require.Equal(t, []interface{}{uint32(1), uint32(2)}, response[0].Fields[0].Values)
                require.Equal(t, []interface{}{"0.75", "0.80"}, response[0].Fields[1].Values)

                // Check loss metrics
                require.Equal(t, "loss", response[1].MetricName)
                require.Equal(t, "train", response[1].StepName)
                require.Len(t, response[1].Fields, 2)
                require.Equal(t, []interface{}{uint32(1), uint32(2)}, response[1].Fields[0].Values)
                require.Equal(t, []interface{}{"0.5", "0.4"}, response[1].Fields[1].Values)
            },
        },
        {
            name:    "No metrics",
            metrics: []model.ModelMetrics{},
            check: func(t *testing.T, response GetModelMetricsResponse) {
                require.Len(t, response, 0)
            },
        },
        {
            name: "Single metric",
            metrics: []model.ModelMetrics{
                {MetricName: "accuracy", StepName: "train", Step: 1, MetricValue: "0.75"},
            },
            check: func(t *testing.T, response GetModelMetricsResponse) {
                require.Len(t, response, 1)
                require.Equal(t, "accuracy", response[0].MetricName)
                require.Equal(t, "train", response[0].StepName)
                require.Len(t, response[0].Fields, 2)
                require.Equal(t, []interface{}{uint32(1)}, response[0].Fields[0].Values)
                require.Equal(t, []interface{}{"0.75"}, response[0].Fields[1].Values)
            },
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // Clear the database
            db.Exec("DELETE FROM model_metrics")

            processID := uuid.New()
            for i := range tc.metrics {
                tc.metrics[i].ProcessID = processID
            }
            insertTestMetrics(t, db, tc.metrics)

            req := setupTestRequest(processID.String())
            result, err := app.getModelMetrics("0", req)
            require.NoError(t, err)
            response, ok := result.(GetModelMetricsResponse)
            require.True(t, ok)

            // Print out the entire response for debugging
            t.Logf("Response: %+v", response)

            if tc.name == "Basic case" {
                // Print out the Values slices for debugging
                t.Logf("Step Values: %+v", response[0].Fields[0].Values)
                t.Logf("Metric Values: %+v", response[0].Fields[1].Values)
            }

            // Run the check function
            tc.check(t, response)
        })
    }
}

func insertTestMetrics(t *testing.T, db *gorm.DB, metrics []model.ModelMetrics) {
    for _, metric := range metrics {
        err := db.Create(&metric).Error
        require.NoError(t, err)
    }
}

func setupTestRequest(processID string) *http.Request {
    req, _ := http.NewRequest("GET", "/process/"+processID+"/model-metrics", nil)
    vars := map[string]string{
        "id": processID,
    }
    return mux.SetURLVars(req, vars)
}