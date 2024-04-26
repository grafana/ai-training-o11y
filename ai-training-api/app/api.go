package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-kit/log/level"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	flatten "github.com/jeremywohl/flatten/v2"

	"github.com/grafana/ai-training-o11y/ai-training-api/middleware"
	"github.com/grafana/ai-training-o11y/ai-training-api/model"
)

const (
	listProcessLimit = 100
)

// RegisterAPI registers all routes to the router.
func (app *App) registerAPI(router *mux.Router) {
	requestMiddleware := middleware.RequestResponseMiddleware(app.logger)

	router.HandleFunc("/process/new", requestMiddleware(app.registerNewProcess)).Methods("POST")
	router.HandleFunc("/process/{id}", requestMiddleware(app.getProcess)).Methods("GET")
	router.HandleFunc("/processes", requestMiddleware(app.listProcess)).Methods("GET")
	// router.HandleFunc("/process/{id}/update-metadata", requestMiddleware(app.updateProcessMetadata)).Methods("POST")
	// router.HandleFunc("/process/{id}/proxy/logs", requestMiddleware(app.proxyProcessLogs)).Methods("POST")
	// router.HandleFunc("/process/{id}/proxy/traces", requestMiddleware(app.proxyProcessTraces)).Methods("POST")
	// router.HandleFunc("/process/{id}/model-metrics", requestMiddleware(app.addModelMetrics)).Methods("POST")

	router.HandleFunc("/group/new", requestMiddleware(app.registerNewGroup)).Methods("POST")
	router.HandleFunc("/group/{id}", requestMiddleware(app.getGroup)).Methods("GET")
}

// registerNewProcess registers a new Process and returns a UUID.
func (a *App) registerNewProcess(tenantID string, req *http.Request) (interface{}, error) {
	level.Info(a.logger).Log("msg", "request received to register new process")

	// Register a new process.
	process := &model.Process{}
	process.ID = uuid.New()
	process.TenantID = tenantID

	// Store process in DB.
	err := a.db(req.Context()).Model(&model.Process{}).Create(process).Error
	if err != nil {
		return nil, fmt.Errorf("error creating process: %w", err)
	}
	level.Info(a.logger).Log("msg", "registered new process", "process_id", process.ID)

	// Read and parse request body.
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, middleware.ErrBadRequest(err)
	}
	defer req.Body.Close()
	var data = map[string]interface{}{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, middleware.ErrBadRequest(err)
	}

	// There are several fields in the request body, some of which are metadata
	// while others contain project and group information. We need to store the
	// metadata in a separate table.
	for key, value := range data {
		switch key {
		case "project":
			process.Project = value.(string)
			continue
		case "group":
			// Store group information in the Group table.
			groupID := uuid.New()
			err = a.db(req.Context()).
				Model(&model.Group{}).
				Create(&model.Group{
					TenantID:  tenantID,
					ID:        groupID,
					Name:      value.(string),
					Processes: []model.Process{*process},
				}).Error
			if err != nil {
				return nil, fmt.Errorf("error creating group: %w", err)
			}
			process.GroupID = &groupID
			continue
		case "metadata":
			// Store metadata information in the Metadata table.
			metadata := value.(map[string]interface{})
			// Flatten JSON body into key-value pairs and store in Metadata table.
			dataMap, err := flatten.Flatten(metadata, "", flatten.DotStyle)
			if err != nil {
				return nil, fmt.Errorf("error flattening metadata: %w", err)
			}
			for mk, mv := range dataMap {
				err = a.db(req.Context()).
					Model(&model.MetadataKV{}).
					Create(&model.MetadataKV{
						TenantID:  tenantID,
						Key:       mk,
						Value:     mv.(string),
						ProcessID: process.ID,
					}).Error
				if err != nil {
					return nil, fmt.Errorf("error creating metadata: %w", err)
				}
			}
			continue
		default:
			level.Error(a.logger).Log("msg", "unknown key in request body", "key", key)
		}
	}

	// Update the process in the DB.
	err = a.db(req.Context()).Model(&model.Process{ID: process.ID}).Updates(process).Error

	// Return the process ID.
	return process, err
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

	return process, err
}

// listProcess returns a list of all processes.
// It limits the number of processes returned to listProcessLimit.
func (a *App) listProcess(tenantID string, req *http.Request) (interface{}, error) {
	processes := []model.Process{}
	err := a.db(req.Context()).
		Where(&model.Process{
			TenantID: tenantID,
		}).Find(&processes).Limit(listProcessLimit).Error
	if err != nil {
		return nil, middleware.ErrNotFound(err)
	}

	level.Info(a.logger).Log("msg", "found processes", "processes", processes)
	return processes, err
}

// registerNewGroup registers a new Group and returns a UUID.
func (a *App) registerNewGroup(tenantID string, req *http.Request) (interface{}, error) {
	level.Info(a.logger).Log("msg", "request received to register new training")

	// Register a new group.
	groupId := uuid.New()

	err := a.db(req.Context()).Create(&model.Group{
		TenantID: tenantID,
		ID:       groupId,
	}).Error

	level.Info(a.logger).Log("msg", "registered new group", "group_id", groupId)
	// Return the groupId.
	return model.Group{ID: groupId}, err
}

// getGroup returns a group by ID.
func (a *App) getGroup(tenantID string, req *http.Request) (interface{}, error) {
	groupId := namedParam(req, "id")
	parsed, err := uuid.Parse(groupId)
	if err != nil {
		return nil, middleware.ErrBadRequest(err)
	}

	group := model.Group{}
	err = a.db(req.Context()).
		Where(&model.Group{
			TenantID: tenantID,
			ID:       parsed,
		}).First(&group).Error
	if err != nil {
		return nil, middleware.ErrNotFound(err)
	}

	level.Info(a.logger).Log("msg", "found group", "group_id", groupId)
	return group, err
}

func namedParam(req *http.Request, name string) string {
	return mux.Vars(req)[name]
}
