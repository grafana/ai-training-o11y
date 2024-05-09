package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-kit/log/level"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	flatten "github.com/jeremywohl/flatten/v2"
	"gorm.io/gorm"

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
	router.HandleFunc("/process/{id}/delete", requestMiddleware(app.deleteProcess)).Methods("POST")
	router.HandleFunc("/processes", requestMiddleware(app.listProcess)).Methods("GET")
	router.HandleFunc("/process/{id}/update-metadata", requestMiddleware(app.updateProcessMetadata)).Methods("POST")
	router.HandleFunc("/process/{id}/model-metrics", requestMiddleware(app.addModelMetrics)).Methods("POST")

	router.HandleFunc("/group/new", requestMiddleware(app.registerNewGroup)).Methods("POST")
	router.HandleFunc("/group/{id}", requestMiddleware(app.getGroup)).Methods("GET")
	router.HandleFunc("/group/{id}/delete", requestMiddleware(app.deleteGroup)).Methods("POST")
}

// registerNewProcess registers a new Process and returns a UUID.
func (a *App) registerNewProcess(tenantID string, req *http.Request) (interface{}, error) {
	level.Info(a.logger).Log("msg", "request received to register new process")

	// Register a new process.
	process := &model.Process{}
	process.ID = uuid.New()
	process.TenantID = tenantID
	process.StartTime = time.Now()

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
		case "user_metadata":
			// Store metadata information in the Metadata table.
			metadata := value.(map[string]interface{})
			// Flatten JSON body into key-value pairs and store in Metadata table.
			dataMap, err := flatten.Flatten(metadata, "", flatten.DotStyle)
			if err != nil {
				return nil, fmt.Errorf("error flattening metadata: %w", err)
			}
			for mk, mv := range dataMap {
				valueType, valueBytes := model.MarshalMetadataValue(mv)
				err = a.db(req.Context()).
					Model(&model.MetadataKV{}).
					Create(&model.MetadataKV{
						TenantID:  tenantID,
						Key:       mk,
						Value:     valueBytes,
						Type:      valueType,
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

// deleteProcess deletes a process by ID.
func (a *App) deleteProcess(tenantID string, req *http.Request) (interface{}, error) {
	processID := namedParam(req, "id")
	parsed, err := uuid.Parse(processID)
	if err != nil {
		return nil, middleware.ErrBadRequest(err)
	}

	// Delete the process.
	err = a.db(req.Context()).
		Where(&model.Process{
			TenantID: tenantID,
			ID:       parsed,
		}).Delete(&model.Process{}).Error
	if err != nil {
		return nil, middleware.ErrNotFound(err)
	}

	level.Info(a.logger).Log("msg", "deleted process", "process_id", processID)
	return nil, err
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

// updateProcessMetadata updates the metadata of a process.
func (a *App) updateProcessMetadata(tenantID string, req *http.Request) (interface{}, error) {
	processID := namedParam(req, "id")
	parsed, err := uuid.Parse(processID)
	if err != nil {
		return nil, middleware.ErrBadRequest(err)
	}

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

	// Only look for metadata in the request body.
	for key, value := range data {
		switch key {
		case "user_metadata":
			metadata := value.(map[string]interface{})
			// Flatten JSON body into key-value pairs and store in Metadata table.
			dataMap, err := flatten.Flatten(metadata, "", flatten.DotStyle)
			if err != nil {
				return nil, fmt.Errorf("error flattening metadata: %w", err)
			}

			// Check if these keys already exist in the Metadata table.
			for mk, mv := range dataMap {
				var metadata model.MetadataKV
				err = a.db(req.Context()).
					Where(&model.MetadataKV{
						TenantID:  tenantID,
						Key:       mk,
						ProcessID: parsed,
					}).First(&metadata).Error
				if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
					// If the key does not exist, create a new entry.
					valueType, valueBytes := model.MarshalMetadataValue(mv)
					err = a.db(req.Context()).
						Model(&model.MetadataKV{}).
						Create(&model.MetadataKV{
							TenantID:  tenantID,
							Key:       mk,
							Value:     valueBytes,
							Type:      valueType,
							ProcessID: parsed,
						}).Error
					if err != nil {
						return nil, fmt.Errorf("error creating metadata: %w", err)
					}
				} else {
					// If the key exists, update the value.
					valueType, valueBytes := model.MarshalMetadataValue(mv)
					err = a.db(req.Context()).
						Model(&model.MetadataKV{}).
						Where(&model.MetadataKV{
							TenantID:  tenantID,
							Key:       mk,
							ProcessID: parsed,
						}).Update("value", valueBytes).
						Update("type", valueType).Error
					if err != nil {
						return nil, fmt.Errorf("error updating metadata: %w", err)
					}
				}
			}
			continue
		default:
			level.Error(a.logger).Log("msg", "unknown key in request body", "key", key)
		}
	}

	level.Info(a.logger).Log("msg", "updated metadata", "process_id", processID)

	// Return the process ID.
	return model.Process{ID: parsed}, err
}

// registerNewGroup registers a new Group and returns a UUID.
func (a *App) registerNewGroup(tenantID string, req *http.Request) (interface{}, error) {
	level.Info(a.logger).Log("msg", "request received to register new group")

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
		Preload("Processes").
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

// deleteGroup deletes a group by ID.
func (a *App) deleteGroup(tenantID string, req *http.Request) (interface{}, error) {
	groupId := namedParam(req, "id")
	parsed, err := uuid.Parse(groupId)
	if err != nil {
		return nil, middleware.ErrBadRequest(err)
	}

	// Delete the group.
	err = a.db(req.Context()).
		Where(&model.Group{
			TenantID: tenantID,
			ID:       parsed,
		}).Delete(&model.Group{}).Error
	if err != nil {
		return nil, middleware.ErrNotFound(err)
	}

	level.Info(a.logger).Log("msg", "deleted group", "group_id", groupId)
	return nil, err
}

// addModelMetrics proxies logs related model-metrics to Loki.
func (a *App) addModelMetrics(tenantID string, req *http.Request) (interface{}, error) {
	// TODO: Integrate with GCom API to find the corresponding Loki TenantID associated
	// with the tenantID.

	// For now, we can just forward the request body as is, to the Loki endpoint.
	// Read the request body.
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, middleware.ErrBadRequest(err)
	}
	defer req.Body.Close()

	level.Debug(a.logger).Log("msg", "forwarding model-metrics to Loki", "body", string(body))

	// Forward the request to the Loki endpoint.
	httpClient := &http.Client{}
	lokiEndpoint := a.lokiAddress
	lokiReq, err := http.NewRequest("POST", lokiEndpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, middleware.ErrBadRequest(err)
	}
	lokiReq.Header.Set("Content-Type", "application/json")
	lokiResp, err := httpClient.Do(lokiReq)
	if err != nil {
		return nil, middleware.ErrBadRequest(err)
	}
	defer lokiResp.Body.Close()

	// Read the response body.
	lokiRespBody, err := io.ReadAll(lokiResp.Body)
	if err != nil {
		return nil, middleware.ErrBadRequest(err)
	}

	// Return the response body.
	return string(lokiRespBody), nil
}

func namedParam(req *http.Request, name string) string {
	return mux.Vars(req)[name]
}
