package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

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

	var rows []model.ModelMetrics

	// Retrieve all relevant metrics from the database
	err = a.db(req.Context()).
		Where("stack_id = ? AND process_id IN ?", stackID, processes).
		Order("metric_name ASC, step_name ASC, process_id ASC, step ASC").
		Find(&rows).Error

	if err != nil {
		return nil, fmt.Errorf("error retrieving model metrics: %w", err)
	}

	// Iterate over the metrics and build the series data
    var response GetModelMetricsResponse
	response.Sections = make(map[string][]Panel)

	var prevKey string
	var currStep uint32

	var stepField *Field
	numFields := 1
    var valueMap map[uuid.UUID]int
	var currentPanel *Panel


    for _, row := range rows {
        panelKey := fmt.Sprintf("%s_%s", row.MetricName, row.StepName)
		// The section name is created by splitting the metricname on / and taking the first part
		// If there is no /, it is "default"
		sectionName := "default"
		panelName := row.MetricName
		parts := strings.Split(row.MetricName, "/");
		if len(parts) > 1 { // There is at least one slash
			sectionName = parts[0]
			panelName = parts[1]
		}

        if panelKey != prevKey {
			// If this section doesn't exist, create it.
			if _, exists := response.Sections[sectionName]; !exists {
				response.Sections[sectionName] = []Panel{}
			}

			// Create a new panel
			currentPanel = &Panel{
				Title: panelName,
				Series: make(DataFrame, 0),
			}

			// Initialize a new step field
            stepField := &Field{
				Name: row.StepName,
				Type: "number",
				Values: make([]interface{}, 0),
			}

			currentPanel.Series = append(currentPanel.Series, *stepField)

			// Zero out the valueMap
			valueMap = make(map[uuid.UUID]int)

			// Append the current panel to the response
			response.Sections[sectionName] = append(response.Sections[sectionName], *currentPanel)
        }

		// If currStep is not defined or is different from the current step, create a new step
		if currStep == 0 || currStep != row.Step {
			stepField.Values = append(stepField.Values, currStep)
			currStep = row.Step
		}

		// If the valueMap doesn't have a key for the current processID, create a new key
		// and append a corresponding series to the current panel
		if _, exists := valueMap[row.ProcessID]; !exists {
			valueMap[row.ProcessID] = numFields;
			currentPanel.Series = append(currentPanel.Series, Field{
				Name: row.ProcessID.String(),
				Type: "number",
				Values: make([]interface{}, 0),
			})
			numFields++;
		}

		currentPanel.Series[valueMap[row.ProcessID]].Values = append(currentPanel.Series[valueMap[row.ProcessID]].Values, row.MetricValue)
    }

	return response, nil
}
