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
	"github.com/google/go-cmp/cmp"
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

func TestTransformMetricsData(t *testing.T) {
	// Helper function to create a string pointer
	strPtr := func(s string) *string {
		return &s
	}

	// Test case
	testCase := struct {
		name     string
		input    []Result
		expected GetModelMetricsResponse
	}{
		name: "Mixed metrics with different sections and step names",
		input: []Result{
			{TenantID: "1", ProcessID: uuid.MustParse("11111111-1111-1111-1111-111111111111"), MetricName: "training/accuracy", StepName: "Epoch", Step: 1, MetricValue: strPtr("0.75")},
			{TenantID: "1", ProcessID: uuid.MustParse("22222222-2222-2222-2222-222222222222"), MetricName: "training/accuracy", StepName: "Epoch", Step: 1, MetricValue: strPtr("0.70")},
			{TenantID: "1", ProcessID: uuid.MustParse("11111111-1111-1111-1111-111111111111"), MetricName: "training/accuracy", StepName: "Epoch", Step: 2, MetricValue: strPtr("0.80")},
			{TenantID: "1", ProcessID: uuid.MustParse("22222222-2222-2222-2222-222222222222"), MetricName: "training/accuracy", StepName: "Epoch", Step: 2, MetricValue: strPtr("0.78")},
			{TenantID: "1", ProcessID: uuid.MustParse("33333333-3333-3333-3333-333333333333"), MetricName: "evaluation/f1_score", StepName: "Step", Step: 1, MetricValue: strPtr("0.65")},
			{TenantID: "1", ProcessID: uuid.MustParse("33333333-3333-3333-3333-333333333333"), MetricName: "evaluation/f1_score", StepName: "Step", Step: 2, MetricValue: strPtr("0.70")},
			{TenantID: "1", ProcessID: uuid.MustParse("44444444-4444-4444-4444-444444444444"), MetricName: "custom_metric", StepName: "Iteration", Step: 1, MetricValue: strPtr("10")},
			{TenantID: "1", ProcessID: uuid.MustParse("44444444-4444-4444-4444-444444444444"), MetricName: "custom_metric", StepName: "Iteration", Step: 2, MetricValue: strPtr("15")},
			{TenantID: "1", ProcessID: uuid.MustParse("55555555-5555-5555-5555-555555555555"), MetricName: "test/accuracy", StepName: "Step", Step: 1, MetricValue: strPtr("0.9")},
		},
		expected: GetModelMetricsResponse{
			Sections: map[string][]Panel{
				"training": {
					{
						Title: "accuracy",
						Series: DataFrame{
							{Name: "Epoch", Type: "number", Values: []interface{}{uint32(1), uint32(2)}},
							{Name: "11111111-1111-1111-1111-111111111111", Type: "number", Values: []interface{}{strPtr("0.75"), strPtr("0.80")}},
							{Name: "22222222-2222-2222-2222-222222222222", Type: "number", Values: []interface{}{strPtr("0.70"), strPtr("0.78")}},
						},
					},
				},
				"evaluation": {
					{
						Title: "f1_score",
						Series: DataFrame{
							{Name: "Step", Type: "number", Values: []interface{}{uint32(1), uint32(2)}},
							{Name: "33333333-3333-3333-3333-333333333333", Type: "number", Values: []interface{}{strPtr("0.65"), strPtr("0.70")}},
						},
					},
				},
				"test": {
					{
						Title: "accuracy",
						Series: DataFrame{
							{
								Name: "Step",
								Type: "number",
								Values: []interface{}{uint32(1)},
							},
							{
								Name: "55555555-5555-5555-5555-555555555555",
								Type: "number",
								Values: []interface{}{strPtr("0.9")},
							},
						},
					},
				},
				"default": {
					{
						Title: "custom_metric",
						Series: DataFrame{
							{Name: "Iteration", Type: "number", Values: []interface{}{uint32(1), uint32(2)}},
							{Name: "44444444-4444-4444-4444-444444444444", Type: "number", Values: []interface{}{strPtr("10"), strPtr("15")}},
						},
					},
				},
			},
		},
	}

    t.Run(testCase.name, func(t *testing.T) {
        result := transformMetricsData(testCase.input)
        if diff := cmp.Diff(testCase.expected, result); diff != "" {
            t.Errorf("transformMetricsData() mismatch (-want +got):\n%s", diff)
        }
    })
}
