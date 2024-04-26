package api

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	dskit_log "github.com/grafana/dskit/log"
	"github.com/grafana/dskit/server"
	"github.com/mattn/go-sqlite3"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/promlog"
	"gorm.io/gorm"

	db "github.com/grafana/ai-training-o11y/ai-training-api/internal"
	"github.com/grafana/ai-training-o11y/ai-training-api/middleware"
	"github.com/grafana/ai-training-o11y/ai-training-api/model"
)

const (
	metricsNamespace = "ai_o11y_training_api"
)

// App is the main application struct.
type App struct {
	_db   *gorm.DB
	dbMux *sync.Mutex

	// The server instance.
	server *server.Server

	logger log.Logger
}

func New(listenAddress string, listenPort int, databaseAddress string, databaseType string, constTenant string, promlogConfig *promlog.Config) (*App, error) {
	// Initialize observability constructs.
	logger := promlog.New(promlogConfig)

	if logger == nil {
		logger = log.NewNopLogger()
	}

	// Initialize the database connection.
	db, err := db.New(logger, databaseAddress, databaseType)
	if err != nil {
		level.Error(logger).Log("msg", "error connecting to database", "err", err)
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	// Migrate the database.
	err = db.AutoMigrate(&model.Process{})
	if err != nil {
		return nil, fmt.Errorf("error migrating Process table: %w", err)
	}
	level.Debug(logger).Log("msg", "checking tables", "process_table_exists", db.Migrator().HasTable(&model.Process{}))
	err = db.AutoMigrate(&model.Group{})
	if err != nil {
		return nil, fmt.Errorf("error migrating Group table: %w", err)
	}
	level.Debug(logger).Log("msg", "checking tables", "group_table_exists", db.Migrator().HasTable(&model.Group{}))
	err = db.AutoMigrate(&model.MetadataKV{})
	if err != nil {
		return nil, fmt.Errorf("error migrating MetadataKV table: %w", err)
	}
	level.Debug(logger).Log("msg", "checking tables", "metadata_kv_table_exists", db.Migrator().HasTable(&model.MetadataKV{}))

	// Create server and router.
	serverLogLevel := dskit_log.Level{}
	serverLogLevel.Set(promlogConfig.Level.String())
	// Create a prometheus registry to avoid "duplicate metrics collector registration attempted"
	// errors when running tests.
	reg := prometheus.NewRegistry()
	s, err := server.New(server.Config{
		Registerer:        reg,
		MetricsNamespace:  metricsNamespace,
		HTTPListenAddress: listenAddress,
		HTTPListenPort:    listenPort,
		LogLevel:          serverLogLevel,
	})
	if err != nil {
		level.Error(logger).Log("msg", "error creating server", "err", err)
		return nil, err
	}

	// Create the App.
	a := &App{
		_db:    db,
		server: s,
		logger: logger,
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("unable to determine underlying sql.DB: %w", err)
	}

	_, needsLock := sqlDB.Driver().(*sqlite3.SQLiteDriver)
	if needsLock {
		a.dbMux = &sync.Mutex{}
	}

	// Register all API routes.
	router := a.server.HTTP.PathPrefix("/api/v1").Subrouter()
	router.Use(middleware.AuthnMiddleware(constTenant))
	a.registerAPI(router)

	// Register the admin routes.
	adm := NewAdmin(a)
	adm.Register(a.server.HTTP.PathPrefix("/admin").Subrouter())

	// Start the server.
	level.Info(logger).Log("msg", "starting server")
	go func() {
		err = a.server.Run()
		if err != nil {
			level.Error(logger).Log("msg", "error running server", "err", err)
		}
	}()

	return a, nil
}

func (a *App) db(ctx context.Context) *gorm.DB {
	if a.dbMux != nil {
		a.dbMux.Lock()
		defer a.dbMux.Unlock()
	}
	return a._db.WithContext(ctx)
}

func (a *App) Shutdown() {
	a.server.Shutdown()
}
