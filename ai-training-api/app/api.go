package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-kit/log/level"
	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/grafana/ai-training-o11y/ai-training-api/middleware"
	"github.com/grafana/ai-training-o11y/ai-training-api/model"
)

// RegisterAPI registers all routes to the router.
func (app *App) registerAPI(router *mux.Router) {
	requestMiddleware := middleware.RequestResponseMiddleware(app.logger)

	router.HandleFunc("/process/new", requestMiddleware(app.registerNewProcess)).Methods("POST")
	router.HandleFunc("/process/{id}", requestMiddleware(app.getProcess)).Methods("GET")
	// router.HandleFunc("/process/{id}/update-metadata", requestMiddleware(app.updateProcessMetadata)).Methods("POST")
	// router.HandleFunc("/process/{id}/proxy/logs", requestMiddleware(app.proxyProcessLogs)).Methods("POST")
	// router.HandleFunc("/process/{id}/proxy/traces", requestMiddleware(app.proxyProcessTraces)).Methods("POST")
	// router.HandleFunc("/process/{id}/model-metrics", requestMiddleware(app.addModelMetrics)).Methods("POST")
	// router.HandleFunc("/process/{id}/state", requestMiddleware(app.updateProcessState)).Methods("POST")

	router.HandleFunc("/trainings/new", requestMiddleware(app.registerNewTraining)).Methods("POST")
}

// CreateProcessResponse is the response for the CreateProcess API.
// It contains a subset of fields from model.Process that we want to return to the user.
type CreateProcessResponse struct {
	ID string `json:"process_uuid"`
}

// GetProcessResponse is the response for the GetProcess API.
type GetProcessResponse struct {
	ID       string             `json:"process_uuid"`
	Metadata []model.MetadataKV `json:"metadata"`
}

// CreateTrainingResponse is the response for the CreateTraining API.
// It contains a subset of fields from model.Training that we want to return to the user.
type CreateTrainingResponse struct {
	ID string `json:"training_uuid"`
}

// registerNewProcess registers a new Process and returns a UUID.
func (a *App) registerNewProcess(tenantID string, req *http.Request) (interface{}, error) {
	level.Info(a.logger).Log("msg", "request received to register new process")

	// Register a new process.
	processID := uuid.New()

	// Read and parse request body.
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, middleware.ErrBadRequest(err)
	}
	defer req.Body.Close()
	var data = map[string]any{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, middleware.ErrBadRequest(err)
	}

	// Store process in DB.
	err = a.db(req.Context()).
		Model(&model.Process{}).
		Create(&model.Process{
			TenantID:  tenantID,
			ID:        processID,
			StartTime: time.Now(),
		}).Error
	if err != nil {
		return nil, fmt.Errorf("error creating process: %w", err)
	}

	// Flatten JSON body into key-value pairs and store in Metadata table.
	for key, value := range data {
		err = a.db(req.Context()).
			Model(&model.MetadataKV{}).
			Create(&model.MetadataKV{
				TenantID:  tenantID,
				Key:       key,
				Value:     value,
				ProcessID: processID,
			}).Error
		if err != nil {
			return nil, fmt.Errorf("error creating metadata: %w", err)
		}
	}

	level.Info(a.logger).Log("msg", "registered new process", "process_id", processID)
	// Return the process ID.
	return CreateProcessResponse{ID: processID.String()}, err
}

// registerNewProcess registers a new Process and returns a UUID.
func (a *App) getProcess(tenantID string, req *http.Request) (interface{}, error) {
	processID := namedParam(req, "id")
	parsed, err := uuid.Parse(processID)
	if err != nil {
		return nil, middleware.ErrBadRequest(err)
	}

	process := model.Process{}
	err = a.db(req.Context()).
		Where(&model.Process{
			TenantID: tenantID,
			ID:       parsed,
		}).First(&process).Error
	if err != nil {
		return nil, middleware.ErrNotFound(err)
	}

	level.Info(a.logger).Log("msg", "found process", "process_id", processID)

	err = a.db(req.Context()).
		Where(&model.MetadataKV{
			ProcessID: parsed,
			TenantID:  tenantID,
		}).Find(&process.Metadata).Error

	return GetProcessResponse{ID: process.ID.String(), Metadata: process.Metadata}, err
}

// registerNewProcess registers a new Process and returns a UUID.
func (a *App) registerNewTraining(tenantID string, req *http.Request) (interface{}, error) {
	level.Info(a.logger).Log("msg", "request received to register new training")

	// Register a new training.
	trainingID := uuid.New()

	err := a.db(req.Context()).Create(&model.Training{
		TenantID: tenantID,
		ID:       trainingID,
	}).Error

	level.Info(a.logger).Log("msg", "registered new training", "training_id", trainingID)
	// Return the training ID.
	return CreateTrainingResponse{ID: trainingID.String()}, err
}

func namedParam(req *http.Request, name string) string {
	return mux.Vars(req)[name]
}
