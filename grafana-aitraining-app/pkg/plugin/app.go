package plugin

import (
	"context"
	"encoding/json"
	"fmt"
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
	lokiDatasourceName string
	mimirDatasourceName string
	metadataUrl string
	metadataToken string
	stackId string
}

func NewApp(_ context.Context, appSettings backend.AppInstanceSettings) (instancemgmt.Instance, error) {
	log.DefaultLogger.Info("Creating new App instance")
	var app App

	var settings map[string]interface{}
	err := json.Unmarshal(appSettings.JSONData, &settings)
	if err != nil {
		log.DefaultLogger.Error("Failed to unmarshal app settings", "error", err)
		return nil, err
	}

	// Helper function to get string value from JSONData
	getStringValue := func(key string) string {
		if value, ok := settings[key]; ok {
			if strValue, ok := value.(string); ok {
				log.DefaultLogger.Info(fmt.Sprintf("%s set from JSONData", key), key, strValue)
				return strValue
			}
			log.DefaultLogger.Warn(fmt.Sprintf("%s in JSONData is not a string", key))
		} else {
			log.DefaultLogger.Warn(fmt.Sprintf("%s not found in settings", key))
		}
		return ""
	}

	// Set values from JSONData
	app.metadataUrl = getStringValue("metadataUrl")
	app.stackId = getStringValue("stackId")
	app.lokiDatasourceName = getStringValue("lokiDatasourceName")
	app.mimirDatasourceName = getStringValue("mimirDatasourceName")

	// Set metadataToken from DecryptedSecureJSONData
	if token, ok := appSettings.DecryptedSecureJSONData["metadataToken"]; ok {
		app.metadataToken = token
		log.DefaultLogger.Info("Metadata token set from DecryptedSecureJSONData")
	} else {
		log.DefaultLogger.Warn("Metadata token not found in DecryptedSecureJSONData")
	}

	// Use a httpadapter (provided by the SDK) for resource calls.
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
