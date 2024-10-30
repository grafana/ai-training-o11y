package api

import (
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
	limitGroupLimit  = 10
)

// RegisterAPI registers all routes to the router.
func (app *App) registerAPI(router *mux.Router) {
	requestMiddleware := middleware.RequestResponseMiddleware(app.logger)

	router.HandleFunc("/process/new", requestMiddleware(app.registerNewProcess)).Methods("POST")
	router.HandleFunc("/process/{id}", requestMiddleware(app.getProcess)).Methods("GET")
	router.HandleFunc("/process/{id}/delete", requestMiddleware(app.deleteProcess)).Methods("POST")
	router.HandleFunc("/processes", requestMiddleware(app.listProcess)).Methods("GET")
	router.HandleFunc("/processes/model-metrics", requestMiddleware(app.getModelMetrics)).Methods("POST")
	router.HandleFunc("/process/{id}/update-metadata", requestMiddleware(app.updateProcessMetadata)).Methods("POST")
	router.HandleFunc("/process/{id}/model-metrics", requestMiddleware(app.addModelMetrics)).Methods("POST")
	router.HandleFunc("/group/new", requestMiddleware(app.registerNewGroup)).Methods("POST")
	router.HandleFunc("/group/{id}", requestMiddleware(app.getGroup)).Methods("GET")
	router.HandleFunc("/groups", requestMiddleware(app.getGroups)).Methods("GET")
	router.HandleFunc("/group/{id}/delete", requestMiddleware(app.deleteGroup)).Methods("POST")
}

// registerNewProcess registers a new Process and returns a UUID.
func (a *App) registerNewProcess(tenantID string, req *http.Request) (interface{}, error) {
	// Register a new process.
	process := &model.Process{}
	process.ID = uuid.New()
	level.Debug(a.logger).Log("msg", "generated new UUID", "process_id", process.ID, "uuid_length", len(process.ID.String()))

	process.TenantID = tenantID
	process.StartTime = time.Now()
	process.Status = "running"

	// Store process in DB.
	level.Debug(a.logger).Log("msg", "attempting to create process", "process_id", process.ID, "tenant_id", tenantID)
	err := a.db(req.Context()).Model(&model.Process{}).Create(process).Error
	if err != nil {
		level.Error(a.logger).Log("msg", "failed to create process", "process_id", process.ID, "error", err)
		return nil, fmt.Errorf("error creating process: %w", err)
	}
	level.Debug(a.logger).Log("msg", "created process in DB", "process_id", process.ID)

	// Read and parse request body.
	body, err := io.ReadAll(req.Body)
	if err != nil {
		level.Error(a.logger).Log("msg", "failed to read request body", "process_id", process.ID, "error", err)
		return nil, middleware.ErrBadRequest(err)
	}
	defer req.Body.Close()

	level.Debug(a.logger).Log("msg", "request body read", "process_id", process.ID, "body_length", len(body))

	var data = map[string]interface{}{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		level.Error(a.logger).Log("msg", "failed to unmarshal request body", "process_id", process.ID, "error", err)
		return nil, middleware.ErrBadRequest(err)
	}

	level.Debug(a.logger).Log("msg", "parsed request body", "process_id", process.ID, "keys", fmt.Sprintf("%v", keys(data)))

	// Process each field
	for key, value := range data {
		level.Debug(a.logger).Log("msg", "processing field", "process_id", process.ID, "key", key, "value_type", fmt.Sprintf("%T", value))

		switch key {
		case "project":
			process.Project = value.(string)
			level.Debug(a.logger).Log("msg", "set project", "process_id", process.ID, "project", process.Project)

		case "group":
			groupName := value.(string)
			level.Debug(a.logger).Log("msg", "processing group", "process_id", process.ID, "group_name", groupName)

			var group model.Group
			err = a.db(req.Context()).
				Where(&model.Group{
					TenantID: tenantID,
					Name:     groupName,
				}).First(&group).Error

			if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
				level.Debug(a.logger).Log("msg", "creating new group", "process_id", process.ID, "group_name", groupName)
				groupID := uuid.New()
				err = a.db(req.Context()).
					Model(&model.Group{}).
					Create(&model.Group{
						TenantID: tenantID,
						ID:       groupID,
						Name:     value.(string),
					}).Error
				if err != nil {
					level.Error(a.logger).Log("msg", "failed to create group", "process_id", process.ID, "error", err)
					return nil, fmt.Errorf("error creating group: %w", err)
				}
				process.GroupID = &groupID
				level.Debug(a.logger).Log("msg", "created new group", "process_id", process.ID, "group_id", groupID)
			} else {
				process.GroupID = &group.ID
				level.Debug(a.logger).Log("msg", "using existing group", "process_id", process.ID, "group_id", group.ID)
			}

		case "user_metadata":
			metadata := value.(map[string]interface{})
			level.Debug(a.logger).Log("msg", "processing metadata", "process_id", process.ID, "metadata_keys", len(metadata))

			dataMap, err := flatten.Flatten(metadata, "", flatten.DotStyle)
			if err != nil {
				level.Error(a.logger).Log("msg", "failed to flatten metadata", "process_id", process.ID, "error", err)
				return nil, fmt.Errorf("error flattening metadata: %w", err)
			}

			level.Debug(a.logger).Log("msg", "flattened metadata", "process_id", process.ID, "flattened_keys", len(dataMap))

			for mk, mv := range dataMap {
				valueType, valueBytes := model.MarshalMetadataValue(mv)
				level.Debug(a.logger).Log("msg", "creating metadata entry", "process_id", process.ID,
					"key", mk, "value_type", valueType)

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
					level.Error(a.logger).Log("msg", "failed to create metadata", "process_id", process.ID,
						"key", mk, "error", err)
					return nil, fmt.Errorf("error creating metadata: %w", err)
				}
			}

		default:
			level.Error(a.logger).Log("msg", "unknown key in request body", "process_id", process.ID, "key", key)
		}
	}

	// Final update
	level.Debug(a.logger).Log("msg", "updating process", "process_id", process.ID,
		"uuid_length", len(process.ID.String()))
	err = a.db(req.Context()).Model(&model.Process{ID: process.ID}).Updates(process).Error
	if err != nil {
		level.Error(a.logger).Log("msg", "failed to update process", "process_id", process.ID, "error", err)
	}

	return process, err
}

// Helper function to get map keys for logging
func keys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
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

	level.Info(a.logger).Log("msg", "found process", "tenantID", tenantID, "process_id", processID)

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

	level.Info(a.logger).Log("msg", "deleted process", "tenantID", tenantID, "process_id", processID)
	return nil, err
}

// listProcess returns a list of all processes.
// It limits the number of processes returned to listProcessLimit.
func (a *App) listProcess(tenantID string, req *http.Request) (interface{}, error) {
	var processes []model.Process

	result := a.db(req.Context()).
		Where("tenant_id = ?", tenantID).
		Order("start_time DESC").
		Limit(listProcessLimit).
		Find(&processes)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, middleware.ErrNotFound(result.Error)
		}
		return nil, result.Error
	}

	level.Info(a.logger).Log("msg", "found processes", "tenantID", tenantID, "len_processes", len(processes))
	return processes, nil
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

	level.Info(a.logger).Log("msg", "updated metadata", "tenantID", tenantID, "process_id", processID)

	// Return the process ID.
	return model.Process{ID: parsed}, err
}

type registerNewGroupRequest struct {
	Name       string      `json:"name"`
	ProcessIDs []uuid.UUID `json:"process_ids"`
}

// registerNewGroup registers a new Group and returns a UUID.
func (a *App) registerNewGroup(tenantID string, req *http.Request) (interface{}, error) {
	// Read and parse request body.
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, middleware.ErrBadRequest(err)
	}
	defer req.Body.Close()
	var data = registerNewGroupRequest{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, middleware.ErrBadRequest(err)
	}

	// Create unique group ID.
	groupId := uuid.New()
	err = a.db(req.Context()).Create(&model.Group{
		TenantID: tenantID,
		ID:       groupId,
		Name:     data.Name,
	}).Error
	if err != nil {
		return nil, fmt.Errorf("error creating group: %w", err)
	}

	// Add processes to the group.
	err = a.db(req.Context()).Model(&model.Process{}).
		Where("id IN ?", data.ProcessIDs).
		Update("group_id", groupId).Error
	if err != nil {
		return nil, fmt.Errorf("error adding processes to group: %w", err)
	}

	level.Info(a.logger).Log("msg", "registered new group", "tenantID", tenantID, "group_id", groupId)
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

	level.Info(a.logger).Log("msg", "found group", "tenantID", tenantID, "group_id", groupId)
	return group, err
}

// getGroups returns a list of all groups.
// It limits the number of groups returned to limitGroupLimit.
func (a *App) getGroups(tenantID string, req *http.Request) (interface{}, error) {
	groups := []model.Group{}
	err := a.db(req.Context()).
		Preload("Processes").
		Where(&model.Group{
			TenantID: tenantID,
		}).Find(&groups).Limit(limitGroupLimit).Error
	if err != nil {
		return nil, middleware.ErrNotFound(err)
	}

	level.Info(a.logger).Log("msg", "found groups", "tenantID", tenantID, "groups", groups)
	return groups, err
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

	level.Info(a.logger).Log("msg", "deleted group", "tenantID", tenantID, "group_id", groupId)
	return nil, err
}

func namedParam(req *http.Request, name string) string {
	return mux.Vars(req)[name]
}
