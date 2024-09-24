package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
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
	process.TenantID = tenantID
	process.StartTime = time.Now()
	process.Status = "running" // Set default status to "running"

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
			// Check if the group already exists.
			groupName := value.(string)
			var group model.Group
			err = a.db(req.Context()).
				Where(&model.Group{
					TenantID: tenantID,
					Name:     groupName,
				}).First(&group).Error
			if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
				// If the group does not exist, create a new group.
				groupID := uuid.New()
				err = a.db(req.Context()).
					Model(&model.Group{}).
					Create(&model.Group{
						TenantID: tenantID,
						ID:       groupID,
						Name:     value.(string),
					}).Error
				if err != nil {
					return nil, fmt.Errorf("error creating group: %w", err)
				}
				process.GroupID = &groupID
			} else {
				process.GroupID = &group.ID
			}
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
	processes := []model.Process{}
	err := a.db(req.Context()).
		Where(&model.Process{
			TenantID: tenantID,
		}).Find(&processes).Limit(listProcessLimit).Error
	if err != nil {
		return nil, middleware.ErrNotFound(err)
	}

	level.Info(a.logger).Log("msg", "found processes", "tenantID", tenantID, "len_processes", len(processes))
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

func (a *App) addModelMetrics(tenantID string, req *http.Request) (interface{}, error) {
    // Extract ProcessID from URL path
    vars := mux.Vars(req)
    processIDStr, ok := vars["id"]
    if !ok {
        return nil, middleware.ErrBadRequest(fmt.Errorf("process ID not provided in URL"))
    }

    // Parse ProcessID to UUID
    processID, err := uuid.Parse(processIDStr)
    if err != nil {
        return nil, middleware.ErrBadRequest(fmt.Errorf("invalid process ID: %w", err))
    }

    // Validate ProcessID exists
    var process model.Process
    if err := a.db(req.Context()).First(&process, "id = ?", processID).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, middleware.ErrNotFound(fmt.Errorf("process not found"))
        }
        return nil, fmt.Errorf("error checking process: %w", err)
    }

    // Parse request body
    var metricsData []struct {
        MetricName string `json:"metric_name"`
        StepName   string `json:"step_name"`
        Points     []struct {
            Step  uint32 `json:"step"`
            Value string `json:"value"`
        } `json:"points"`
    }

    if err := json.NewDecoder(req.Body).Decode(&metricsData); err != nil {
        return nil, middleware.ErrBadRequest(err)
    }

    // Convert tenantID to uint64 for StackID
    stackID, err := strconv.ParseUint(tenantID, 10, 64)
    if err != nil {
        return nil, middleware.ErrBadRequest(fmt.Errorf("invalid tenant ID: %w", err))
    }

    var createdMetrics []model.ModelMetrics

    // Start a transaction
    tx := a.db(req.Context()).Begin()
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

            // Validate fields
            if err := validateModelMetric(&metric); err != nil {
                tx.Rollback()
                return nil, middleware.ErrBadRequest(err)
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
    if len(m.MetricValue) == 0 || len(m.MetricValue) > 64 {
        return fmt.Errorf("metric value must be between 1 and 64 characters")
    }
    return nil
}

func namedParam(req *http.Request, name string) string {
	return mux.Vars(req)[name]
}
