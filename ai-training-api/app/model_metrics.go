package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/gorm"

	"github.com/grafana/ai-training-o11y/ai-training-api/middleware"
	"github.com/grafana/ai-training-o11y/ai-training-api/model"
)

type ModelMetricsRequest struct {
	MetricName string `json:"metric_name"`
	StepName   string `json:"step_name"`
	Points     []struct {
		Step  uint32 `json:"step"`
		Value string `json:"value"`
	} `json:"points"`
}

type AddModelMetricsResponse struct {
    Message        string `json:"message"`
    MetricsCreated int    `json:"metricsCreated"`
}

// This is for return
// We want an array of objects that contain grafana dataframes
// For visualizing
type DataFrame struct {
    Name   string        `json:"name"`
    Type   string        `json:"type"`
    Values []interface{} `json:"values"`
}

// To make it less painful to unmarshal and group them
type DataFrameWrapper struct {
    MetricName string   `json:"MetricName"`
    StepName   string   `json:"StepName"`
    Fields     []DataFrame  `json:"fields"`
}

type GetModelMetricsResponse []DataFrameWrapper


func (a *App) addModelMetrics(tenantID string, req *http.Request) (interface{}, error) {
	// Extract and validate ProcessID
	processID, err := extractAndValidateProcessID(req)
	if err != nil {
		return nil, err
	}

	// Validate ProcessID exists
	if err := a.validateProcessExists(req.Context(), processID); err != nil {
		return nil, err
	}

	// Parse and validate the request body
	metricsData, err := parseAndValidateModelMetricsRequest(req)
	if err != nil {
		return nil, err
	}

	// Convert tenantID to uint64 for StackID
	stackID, err := strconv.ParseUint(tenantID, 10, 64)
	if err != nil {
		return nil, middleware.ErrBadRequest(fmt.Errorf("invalid tenant ID: %w", err))
	}

	// Save the metrics and get the count of created metrics
	createdCount, err := a.saveModelMetrics(req.Context(), stackID, processID, metricsData)
	if err != nil {
		return nil, err
	}

	// Return a JSON response with success message and count of metrics inserted
	response := map[string]interface{}{
		"message":        "Metrics successfully added",
		"metricsCreated": createdCount,
	}

	return response, nil
}

func extractAndValidateProcessID(req *http.Request) (uuid.UUID, error) {
    vars := mux.Vars(req)
    if vars == nil {
        return uuid.Nil, fmt.Errorf("mux.Vars(req) returned nil")
    }

    processIDStr, ok := vars["id"]
    if !ok {
        return uuid.Nil, middleware.ErrBadRequest(fmt.Errorf("process ID not provided in URL"))
    }

    // This case handles when the ID is provided in the URL but is empty
    if processIDStr == "" {
        return uuid.Nil, middleware.ErrBadRequest(fmt.Errorf("process ID is empty"))
    }

    processID, err := uuid.Parse(processIDStr)
    if err != nil {
        return uuid.Nil, middleware.ErrBadRequest(fmt.Errorf("invalid process ID: %w", err))
    }

    return processID, nil
}

func (a *App) validateProcessExists(ctx context.Context, processID uuid.UUID) error {
	var process model.Process
	if err := a.db(ctx).First(&process, "id = ?", processID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return middleware.ErrNotFound(fmt.Errorf("process not found"))
		}
		return fmt.Errorf("error checking process: %w", err)
	}
	return nil
}

func parseAndValidateModelMetricsRequest(req *http.Request) ([]ModelMetricsRequest, error) {
	var metricsData []ModelMetricsRequest

	if err := json.NewDecoder(req.Body).Decode(&metricsData); err != nil {
		return nil, middleware.ErrBadRequest(err)
	}

	for _, metric := range metricsData {
		if err := validateModelMetricRequest(&metric); err != nil {
			return nil, middleware.ErrBadRequest(err)
		}
	}

	return metricsData, nil
}

func validateModelMetricRequest(m *ModelMetricsRequest) error {
	if len(m.MetricName) == 0 || len(m.MetricName) > 32 {
		return fmt.Errorf("metric name must be between 1 and 32 characters")
	}
	if len(m.StepName) == 0 || len(m.StepName) > 32 {
		return fmt.Errorf("step name must be between 1 and 32 characters")
	}
	for _, point := range m.Points {
		if point.Step == 0 {
			return fmt.Errorf("step must be a positive number")
		}
		if len(point.Value) == 0 || len(point.Value) > 64 {
			return fmt.Errorf("metric value must be between 1 and 64 characters")
		}
	}
	return nil
}

func (a *App) saveModelMetrics(ctx context.Context, stackID uint64, processID uuid.UUID, metricsData []ModelMetricsRequest) (int, error) {
	var createdCount int

	// Start a transaction
	tx := a.db(ctx).Begin()
	if tx.Error != nil {
		return 0, fmt.Errorf("error starting transaction: %w", tx.Error)
	}

	for _, metricData := range metricsData {
		for _, point := range metricData.Points {
			metric := model.ModelMetrics{
				StackID:     stackID,
				ProcessID:   processID,
				MetricName:  metricData.MetricName,
				StepName:    metricData.StepName,
				Step:        point.Step,
				MetricValue: point.Value,
			}

			// Save to database
			if err := tx.Create(&metric).Error; err != nil {
				tx.Rollback()
				return 0, fmt.Errorf("error creating model metric: %w", err)
			}
			createdCount++
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return 0, fmt.Errorf("error committing transaction: %w", err)
	}

	return createdCount, nil
}

func (a *App) getModelMetrics(tenantID string, req *http.Request) (interface{}, error) {

	// Extract and validate ProcessID
	processID, err := extractAndValidateProcessID(req)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	// Convert tenantID to uint64 for StackID
	stackID, err := strconv.ParseUint(tenantID, 10, 64)
	if err != nil {
		return nil, middleware.ErrBadRequest(fmt.Errorf("invalid tenant ID: %w", err))
	}

	// Retrieved from DB
	var rows []model.ModelMetrics

	// Retrieve all relevant metrics from the database
	err = a.db(req.Context()).
		Where("stack_id = ? AND process_id = ?", stackID, processID).
		Order("metric_name ASC, step_name ASC, step ASC").
		Find(&rows).Error

	if err != nil {
		return nil, fmt.Errorf("error retrieving model metrics: %w", err)
	}

	// Iterate over the metrics and build the series data
    var response GetModelMetricsResponse
    var currentWrapper *DataFrameWrapper
    var stepSlice, valueSlice *[]interface{}

    for _, row := range rows {
        currSeriesKey := fmt.Sprintf("%s_%s", row.MetricName, row.StepName)

        if currentWrapper == nil || currSeriesKey != fmt.Sprintf("%s_%s", currentWrapper.MetricName, currentWrapper.StepName) {
            // We've encountered a new series, so create a new wrapper
            newStepSlice := make([]interface{}, 0, 100)  // Pre-allocate with a capacity of 100
            newValueSlice := make([]interface{}, 0, 100) // Pre-allocate with a capacity of 100
            stepSlice = &newStepSlice
            valueSlice = &newValueSlice
            
            currentWrapper = &DataFrameWrapper{
                MetricName: row.MetricName,
                StepName:   row.StepName,
                Fields: []DataFrame{
                    {
                        Name:   row.StepName,
                        Type:   "number",
                        Values: *stepSlice,
                    },
                    {
                        Name:   row.MetricName,
                        Type:   "number",
                        Values: *valueSlice,
                    },
                },
            }
            response = append(response, *currentWrapper)
        }

        // Append the step and metricValue to the slices
        *stepSlice = append(*stepSlice, row.Step)
        *valueSlice = append(*valueSlice, row.MetricValue)
    }

    return response, nil
}
