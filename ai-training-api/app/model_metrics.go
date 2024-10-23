package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/gorm"

	"github.com/grafana/ai-training-o11y/ai-training-api/middleware"
	"github.com/grafana/ai-training-o11y/ai-training-api/model"
)

// Incoming format is an array of these
type AddModelMetricsPayload struct {
	StepName  string                 `json:"step_name"`
	StepValue uint32                 `json:"step_value"`
	Metrics   map[string]json.Number `json:"metrics"`
}

type AddModelMetricsResponse struct {
	Message        string `json:"message"`
	MetricsCreated uint32 `json:"metricsCreated"`
}

// Result struct to hold our query results
type Result struct {
	TenantID    string
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
	Title  string    `json:"title"`
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

	// Save the metrics and get the count of created metrics
	createdCount, err := a.saveModelMetrics(req.Context(), tenantID, processID, metricsData)
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

func parseAndValidateModelMetricsRequest(req *http.Request) ([]model.ModelMetrics, error) {
	var metricsData []AddModelMetricsPayload

	decoder := json.NewDecoder(req.Body)

	if err := decoder.Decode(&metricsData); err != nil {
		return nil, middleware.ErrBadRequest(fmt.Errorf("invalid JSON: %v", err))
	}

	var metrics []model.ModelMetrics

	for _, item := range metricsData {
		for metricName, metricValue := range item.Metrics {
			metric := model.ModelMetrics{
				MetricName:  metricName,
				StepName:    item.StepName,
				Step:        item.StepValue,
				MetricValue: metricValue.String(),
			}

			if err := validateModelMetric(&metric); err != nil {
				return nil, middleware.ErrBadRequest(fmt.Errorf("invalid metric: %v", err))
			}

			metrics = append(metrics, metric)
		}
	}

	return metrics, nil
}

func validateModelMetric(m *model.ModelMetrics) error {
	if len(m.MetricName) == 0 || len(m.MetricName) > 32 {
		return fmt.Errorf("metric name must be between 1 and 32 characters")
	}
	if len(m.StepName) == 0 || len(m.StepName) > 32 {
		return fmt.Errorf("step name must be between 1 and 32 characters")
	}

	if m.Step == 0 {
		return fmt.Errorf("step must be a positive number")
	}
	if m.MetricValue == "" {
		return fmt.Errorf("metric value cannot be empty")
	}
	return nil
}

func (a *App) saveModelMetrics(ctx context.Context, tenantID string, processID uuid.UUID, metricsData []model.ModelMetrics) (int, error) {
	var createdCount int

	// Start a transaction
	tx := a.db(ctx).Begin()
	if tx.Error != nil {
		return 0, fmt.Errorf("error starting transaction: %w", tx.Error)
	}

	for _, metric := range metricsData {
		metric.TenantID = tenantID
		metric.ProcessID = processID

		// Save to database
		if err := tx.Create(&metric).Error; err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("error creating model metric: %w", err)
		}
		createdCount++
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return 0, fmt.Errorf("error committing transaction: %w", err)
	}

	return createdCount, nil
}

func getCompleteMetrics(ctx context.Context, db *gorm.DB, tenantID string, processes []string) ([]Result, error) {
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
			WHERE tenant_id = ? AND process_id IN ?
		),
		metric_steps AS (
			SELECT metric_name, step_name, step
			FROM model_metrics
			WHERE tenant_id = ? AND process_id IN ?
		),
		all_combinations AS (
			SELECT DISTINCT
				pm.process_id,
				pm.metric_name,
				pm.step_name,
				ms.step
			FROM 
				process_metrics pm
			JOIN
				metric_steps ms ON pm.metric_name = ms.metric_name AND pm.step_name = ms.step_name
		)
		SELECT DISTINCT
			ac.process_id,
			ac.metric_name,
			ac.step_name,
			ac.step,
			d.metric_value
		FROM 
			all_combinations ac
		LEFT JOIN
			model_metrics d ON d.tenant_id = ? 
							AND d.process_id = ac.process_id 
							AND d.metric_name = ac.metric_name 
							AND d.step_name = ac.step_name
							AND d.step = ac.step
		ORDER BY ac.step ASC
    `

	err := db.WithContext(ctx).Raw(query, tenantID, uuidProcesses, tenantID, uuidProcesses, tenantID).Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	return results, nil
}

func transformMetricsData(results []Result) GetModelMetricsResponse {
	// Group results by metric_name and step_name
	// This makes it easy to separate panels: each []Result is a panel
	// This "only" leaves turning it into a DataFrame to send to the frontend
	groupedData := make(map[string]map[string][]Result)
	for _, r := range results {
		if _, ok := groupedData[r.MetricName]; !ok {
			groupedData[r.MetricName] = make(map[string][]Result)
		}

		if _, ok := groupedData[r.MetricName][r.StepName]; !ok {
			groupedData[r.MetricName][r.StepName] = make([]Result, 0)
		}
		groupedData[r.MetricName][r.StepName] = append(groupedData[r.MetricName][r.StepName], r)
	}

	fmt.Println(groupedData)

	response := GetModelMetricsResponse{
		Sections: make(map[string][]Panel),
	}

	for metricName, stepData := range groupedData {
		sectionName := "default"
		displayName := metricName
		first, second, hasSectionName := strings.Cut(metricName, "/")
		if hasSectionName {
			sectionName = first
			displayName = second
		}

		panels := []Panel{}
		// Construct an individual panel
		for stepName, metricRows := range stepData {
			// Create empty panel

			newPanel := Panel{
				Title:  displayName,
				Series: make([]Field, 0),
			}

			// Create step field
			stepField := Field{
				Name:   stepName,
				Type:   "number",
				Values: make([]interface{}, 0),
			}

			steps := make([]interface{}, 0)
			lastStep := uint32(0)
			processFields := make(map[string]*Field)
			processOrder := make([]string, 0)
			for _, row := range metricRows {
				// Add new steps
				if row.Step != lastStep {
					steps = append(steps, row.Step)
					lastStep = row.Step
				}

				// Check if the field already exists
				if _, ok := processFields[row.ProcessID.String()]; !ok {
					processFields[row.ProcessID.String()] = &Field{
						Name:   row.ProcessID.String(),
						Type:   "number",
						Values: make([]interface{}, 0),
					}

					processOrder = append(processOrder, row.ProcessID.String())
				}
				// Append the value to the field
				processFields[row.ProcessID.String()].Values = append(processFields[row.ProcessID.String()].Values, row.MetricValue)
			}

			stepField.Values = steps
			// turn processFields into a slice
			processFieldsSlice := make([]Field, 0)
			for _, proc_id := range processOrder {
				processFieldsSlice = append(processFieldsSlice, *processFields[proc_id])
			}

			newPanel.Series = append(newPanel.Series, stepField)
			// unpack the processFieldsSlice and append it into stepField.values too

			// This is broken and I am going to ask claude how to fix it
			newPanel.Series = append(newPanel.Series, processFieldsSlice...)
			panels = append(panels, newPanel)
		}

		response.Sections[sectionName] = panels
	}

	return response
}

func (a *App) getModelMetrics(tenantID string, req *http.Request) (interface{}, error) {
	// parse request body into an array
	var processes []string

	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&processes); err != nil {
		return nil, middleware.ErrBadRequest(fmt.Errorf("invalid JSON: %v", err))
	}

	results, err := getCompleteMetrics(req.Context(), a.db(req.Context()), tenantID, processes)
	if err != nil {
		return nil, fmt.Errorf("error getting complete metrics: %w", err)
	}

	transformedMetricsData := transformMetricsData(results)

	return transformedMetricsData, nil // Return results instead of nil
}
