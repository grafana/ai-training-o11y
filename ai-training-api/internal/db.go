// Based on https://github.com/grafana/machine-learning/blob/main/mlapi/internal/db/db.go
package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/prometheus"
)

const (
	SQLite = "sqlite3"
	MySQL  = "mysql"
)

func New(logger log.Logger, addr, dbType string) (*gorm.DB, error) {
	if logger == nil {
		logger = log.NewNopLogger()
	}

	var db *gorm.DB
	err := retry.Do(
		func() error {
			var err error
			db, err = newDB(logger, addr, dbType)
			return err
		},
		retry.Delay(5*time.Second),
		retry.Attempts(20),
	)
	if err != nil {
		return nil, err
	}

	err = db.Use(prometheus.New(prometheus.Config{}))
	if err != nil {
		return nil, err
	}
	err = db.Use(&tracer{})
	if err != nil {
		return nil, err
	}

	// Cleanup idle connections after 5 minutes so we don't hold unnecessary
	// resources for a long time. Closing the connections may also help with
	// invalid connection issues coming from Azure.
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)

	return db, nil
}

func NewFromEnvironment(logger log.Logger, defaultAddr, defaultType string) (*gorm.DB, error) {
	addr := os.Getenv("DB_ADDR")
	if addr == "" {
		addr = defaultAddr
	}

	dbType := os.Getenv("DB_TYPE")
	if dbType == "" {
		dbType = defaultType
	}
	return New(logger, addr, dbType)
}

func newDB(logger log.Logger, addr, dbType string) (*gorm.DB, error) {
	cfg := &gorm.Config{
		Logger: &gormLogger{logger},
	}
	switch dbType {
	case MySQL:
		return gorm.Open(mysql.Open(addr), cfg)
	case SQLite:
		db, err := gorm.Open(sqlite.Open(addr), cfg)
		if err != nil {
			return nil, fmt.Errorf("connecting to sqlite: %w", err)
		}
		// Foreign key support is disabled by default in sqlite, but
		// we use it for cascading deletes of holidays.
		err = db.Exec("PRAGMA foreign_keys = ON", nil).Error
		if err != nil {
			return nil, fmt.Errorf("enabling foreign keys: %w", err)
		}
		return db, nil
	}
	return nil, fmt.Errorf("unknown database type: `%s`", dbType)
}

type gormLogger struct {
	log.Logger
}

// LogMode is a noop as the level is set by the go-kit logger we are wrapping.
func (logger *gormLogger) LogMode(logger.LogLevel) logger.Interface {
	return logger
}

func (logger *gormLogger) Info(_ context.Context, msg string, args ...interface{}) {
	level.Info(logger).Log("msg", fmt.Sprintf(msg, args...))
}

func (logger *gormLogger) Warn(_ context.Context, msg string, args ...interface{}) {
	level.Warn(logger).Log("msg", fmt.Sprintf(msg, args...))
}

func (logger *gormLogger) Error(_ context.Context, msg string, args ...interface{}) {
	level.Error(logger).Log("msg", fmt.Sprintf(msg, args...))
}

func (logger *gormLogger) Trace(_ context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	if err != nil {
		level.Error(logger).Log("msg", "error running database transaction", "err", err, "elapsed", elapsed, "sql", sql, "rows", rows)
		return
	}

	if elapsed > time.Second {
		level.Warn(logger).Log("msg", "slow database query", "elapsed", elapsed, "sql", sql, "rows", rows)
		return
	}
	level.Debug(logger).Log("msg", "database query", "elapsed", elapsed, "sql", sql, "rows", rows)
}
