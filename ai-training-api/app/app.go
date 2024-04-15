package api

import (
	"context"
	"net/http"
	"sync"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	dskit_log "github.com/grafana/dskit/log"
	"github.com/grafana/dskit/server"
	"github.com/prometheus/common/promlog"
	"gorm.io/gorm"

	db "github.com/grafana/ai-o11y/metadata-service/internal"
	"github.com/grafana/ai-o11y/metadata-service/middleware"
	"github.com/grafana/ai-o11y/metadata-service/model"
)

const (
	metricsNamespace = "ai_o11y_metadata_service"
)

// App is the main application struct.
type App struct {
	_db   *gorm.DB
	dbMux *sync.Mutex

	// The server instance.
	server *server.Server

	logger log.Logger
}

func New(listenAddress *string, listenPort *int, databaseAddress *string, databaseType *string, promlogConfig *promlog.Config) (*App, error) {
	// Initialize observability constructs.
	logger := promlog.New(promlogConfig)

	if logger == nil {
		logger = log.NewNopLogger()
	}

	// Initialize the database connection.
	db, err := db.New(logger, *databaseAddress, *databaseType)
	if err != nil {
		level.Error(logger).Log("msg", "error connecting to database", "err", err)
		return nil, err
	}

	// Create server and router.
	serverLogLevel := dskit_log.Level{}
	serverLogLevel.Set(promlogConfig.Level.String())
	s, err := server.New(server.Config{
		MetricsNamespace:  metricsNamespace,
		HTTPListenAddress: *listenAddress,
		HTTPListenPort:    *listenPort,
		LogLevel:          serverLogLevel,
	})
	if err != nil {
		level.Error(logger).Log("msg", "error creating server", "err", err)
		return nil, err
	}

	// Create the App.
	a := &App{
		_db:    db,
		dbMux:  &sync.Mutex{},
		server: s,
		logger: logger,
	}

	// Register all API routes.
	router := a.server.HTTP.PathPrefix("/api/v1").Subrouter()
	a.registerAPI(router)

	// Start the server.
	level.Info(logger).Log("msg", "starting server")
	err = a.server.Run()
	if err != nil {
		level.Error(logger).Log("msg", "error running server", "err", err)
		return nil, err
	}

	return a, nil
}

func (a *App) db(ctx context.Context) *gorm.DB {
	if a.dbMux != nil {
		a.dbMux.Lock()
		defer a.dbMux.Unlock()
	}
	return a._db.WithContext(ctx)
}

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

// registerNewProcess registers a new Process and returns a UUID.
func (a *App) registerNewProcess(tenantID string, req *http.Request) (interface{}, error) {
	// Register a new process.
	processID := uuid.New()

	// TODO: read and parse request body

	err := a.db(req.Context()).Create(&model.Process{
		TenantID: tenantID,
		ID:       processID,
	}).Error

	level.Info(a.logger).Log("msg", "registered new process", "process_id", processID)
	// Return the process ID.
	return processID, err
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
