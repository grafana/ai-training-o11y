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

type modelMetricsRequest struct {
	MetricName string `json:"metric_name"`
	StepName   string `json:"step_name"`
	Points     []struct {
		Step  uint32 `json:"step"`
		Value string `json:"value"`
	} `json:"points"`
}

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

	createdMetrics, err := a.saveModelMetrics(req.Context(), stackID, processID, metricsData)
	if err != nil {
		return nil, err
	}

	return createdMetrics, nil
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

func parseAndValidateModelMetricsRequest(req *http.Request) ([]modelMetricsRequest, error) {
	var metricsData []modelMetricsRequest

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

func validateModelMetricRequest(m *modelMetricsRequest) error {
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

func (a *App) saveModelMetrics(ctx context.Context, stackID uint64, processID uuid.UUID, metricsData []modelMetricsRequest) ([]model.ModelMetrics, error) {
	var createdMetrics []model.ModelMetrics

	// Start a transaction
	tx := a.db(ctx).Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("error starting transaction: %w", tx.Error)
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
				return nil, fmt.Errorf("error creating model metric: %w", err)
			}

			createdMetrics = append(createdMetrics, metric)
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	return createdMetrics, nil
}
