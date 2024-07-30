package plugin

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/backend/resource/httpadapter"
)

// Make sure App implements required interfaces. This is important to do
// since otherwise we will only get a not implemented error response from plugin in
// runtime. Plugin should not implement all these interfaces - only those which are
// required for a particular task.
var (
	_ backend.CallResourceHandler   = (*App)(nil)
	_ instancemgmt.InstanceDisposer = (*App)(nil)
	_ backend.CheckHealthHandler    = (*App)(nil)
)

// App is an example app backend plugin which can respond to data queries.
type App struct {
	backend.CallResourceHandler
	stackId string
	metadataUrl string
	metadataToken string
}

// NewApp creates a new example *App instance.
func NewApp(_ context.Context, appSettings backend.AppInstanceSettings) (instancemgmt.Instance, error) {
	log.DefaultLogger.Info("Creating new App instance")
	var app App

	var settings map[string]interface{}
    err := json.Unmarshal(appSettings.JSONData, &settings)
    if err != nil {
        log.DefaultLogger.Error("Failed to unmarshal app settings", "error", err)
        return nil, err
    }

	// Check if metadataUrl exists in settings
	if value, ok := settings["metadataUrl"]; ok {
		switch url := value.(type) {
		case string:
			app.metadataUrl = url
			log.DefaultLogger.Info("Metadata URL set", "url", app.metadataUrl)
		default:
			app.metadataUrl = "" 
			log.DefaultLogger.Warn("Metadata URL in settings but is not a string, using empty string")
		}
	} else {
		// If metadataUrl is not found, set it to an empty string or a default value
		app.metadataUrl = "" 
		log.DefaultLogger.Warn("Metadata URL not found in settings, using empty string")
	}

	// Use a httpadapter (provided by the SDK) for resource calls. This allows us
	// to use a *http.ServeMux for resource calls, so we can map multiple routes
	// to CallResource without having to implement extra logic.
	mux := http.NewServeMux()
	app.registerRoutes(mux)
	app.CallResourceHandler = httpadapter.New(mux)

	log.DefaultLogger.Info("App instance created successfully")
	return &app, nil
}

// Dispose here tells plugin SDK that plugin wants to clean up resources when a new instance
// created.
func (a *App) Dispose() {
	log.DefaultLogger.Info("Disposing App instance")
	// cleanup
}

// CheckHealth handles health checks sent from Grafana to the plugin.
func (a *App) CheckHealth(_ context.Context, _ *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	log.DefaultLogger.Info("Performing health check")
	result := &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: "ok",
	}
	log.DefaultLogger.Info("Health check completed", "status", result.Status, "message", result.Message)
	return result, nil
}
