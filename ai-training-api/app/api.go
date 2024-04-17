package api

import (
	"net/http"

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
	// router.HandleFunc("/process/{id}/custom-logs", requestMiddleware(app.addProcessCustomLogs)).Methods("POST")
	// router.HandleFunc("/process/{id}/state", requestMiddleware(app.updateProcessState)).Methods("POST")
}

// CreateProcessResponse is the response for the CreateProcess API.
// It contains a subset of fields from model.Process that we want to return to the user.
type CreateProcessResponse struct {
	ID string `json:"process_uuid"`
}

// registerNewProcess registers a new Process and returns a UUID.
func (a *App) registerNewProcess(tenantID string, req *http.Request) (interface{}, error) {
	level.Info(a.logger).Log("msg", "request received to register new process")

	// Register a new process.
	processID := uuid.New()

	// TODO: read and parse request body

	err := a.db(req.Context()).Create(&model.Process{
		TenantID: tenantID,
		ID:       processID,
	}).Error

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

	level.Info(a.logger).Log("msg", "registered new process", "process_id", processID)
	// Return the process ID.
	return processID, err
}

func namedParam(req *http.Request, name string) string {
	return mux.Vars(req)[name]
}
