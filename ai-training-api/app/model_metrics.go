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

// Incoming format is an array of these
type ModelMetricsSeries struct {
	MetricName string `json:"metric_name"`
	StepName   string `json:"step_name"`
	Points     []struct {
		Step  uint32 `json:"step"`
		Value json.Number `json:"value"`
	} `json:"points"`
}

type AddModelMetricsResponse struct {
    Message        string `json:"message"`
    MetricsCreated uint32    `json:"metricsCreated"`
}

// Result struct to hold our query results
type Result struct {
    StackID     uint64
    ProcessID   uuid.UUID
    MetricName  string
    StepName    string
    Step        uint32
    MetricValue *string // Pointer to allow for NULL values
}

// This is for return
// We want an array of objects that contain grafana dataframes
// For visualizing
type Field struct {
    Name   string        `json:"name"`
    Type   string        `json:"type"`
    Values []interface{} `json:"values"`
}

type DataFrame []Field

type Panel struct {
	Title string `json:"title"`
	Series DataFrame `json:"series"`
}

type GetModelMetricsResponse struct {
    Sections map[string][]Panel `json:"sections"`
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

func parseAndValidateModelMetricsRequest(req *http.Request) ([]ModelMetricsSeries, error) {
    var metricsData []ModelMetricsSeries

    decoder := json.NewDecoder(req.Body)

    if err := decoder.Decode(&metricsData); err != nil {
        return nil, middleware.ErrBadRequest(fmt.Errorf("invalid JSON: %v", err))
    }

	fmt.Println(metricsData)

    for _, metric := range metricsData {
        if err := validateModelMetricRequest(&metric); err != nil {
            return nil, middleware.ErrBadRequest(err)
        }
    }

    return metricsData, nil
}

func validateModelMetricRequest(m *ModelMetricsSeries) error {
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
        if point.Value.String() == "" {
            return fmt.Errorf("metric value cannot be empty")
        }
        // Validate that Value is a valid number
        if _, err := point.Value.Float64(); err != nil {
            return fmt.Errorf("invalid numeric value: %v", err)
        }
    }
    return nil
}

func (a *App) saveModelMetrics(ctx context.Context, stackID uint64, processID uuid.UUID, metricsData []ModelMetricsSeries) (int, error) {
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
				MetricValue: point.Value.String(),
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

func getCompleteMetrics(ctx context.Context, db *gorm.DB, stackID uint64, processes []string) ([]Result, error) {
    // Convert []string to []uuid.UUID
    uuidProcesses := make([]uuid.UUID, 0, len(processes))
    for _, p := range processes {
        uid, err := uuid.Parse(p)
        if err != nil {
            return nil, fmt.Errorf("invalid UUID string: %s", p)
        }
        uuidProcesses = append(uuidProcesses, uid)
    }

    var results []Result

    query := `
			WITH process_metrics AS (
			SELECT DISTINCT process_id, metric_name, step_name
			FROM model_metrics
			WHERE stack_id = ? AND process_id IN ?
		),
		metric_steps AS (
			SELECT metric_name, step_name, step
			FROM model_metrics
			WHERE stack_id = ? AND process_id IN ?
		),
		all_combinations AS (
			SELECT 
				pm.process_id,
				pm.metric_name,
				pm.step_name,
				ms.step
			FROM 
				process_metrics pm
			JOIN
				metric_steps ms ON pm.metric_name = ms.metric_name AND pm.step_name = ms.step_name
		)
		SELECT 
			ac.process_id,
			ac.metric_name,
			ac.step_name,
			ac.step,
			d.metric_value
		FROM 
			all_combinations ac
		LEFT JOIN
			model_metrics d ON d.stack_id = ? 
							AND d.process_id = ac.process_id 
							AND d.metric_name = ac.metric_name 
							AND d.step_name = ac.step_name
							AND d.step = ac.step
		ORDER BY 
			ac.metric_name ASC, ac.step_name ASC, ac.step ASC, ac.process_id ASC
    `

    err := db.WithContext(ctx).Raw(query, stackID, uuidProcesses, stackID, uuidProcesses, stackID).Scan(&results).Error
    if err != nil {
        return nil, fmt.Errorf("error executing query: %v", err)
    }

    return results, nil
}

func (a *App) getModelMetrics(tenantID string, req *http.Request) (interface{}, error) {
	// parse request body into an array
	var processes []string

	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&processes); err != nil {
		return nil, middleware.ErrBadRequest(fmt.Errorf("invalid JSON: %v", err))
	}

	// Convert tenantID to uint64 for StackID
	stackID, err := strconv.ParseUint(tenantID, 10, 64)
	if err != nil {
		return nil, middleware.ErrBadRequest(fmt.Errorf("invalid tenant ID: %w", err))
	}

    results, err := getCompleteMetrics(req.Context(), a.db(req.Context()), stackID, processes)
    if err != nil {
        return nil, fmt.Errorf("error getting complete metrics: %w", err)
    }

    // Print results to console
    fmt.Println("Results:")
    for _, r := range results {
        metricValue := "NULL"
        if r.MetricValue != nil {
            metricValue = *r.MetricValue
        }
        fmt.Printf("ProcessID: %s, MetricName: %s, StepName: %s, Step: %d, MetricValue: %s\n",
            r.ProcessID, r.MetricName, r.StepName, r.Step, metricValue)
    }

    return results, nil  // Return results instead of nil
}
