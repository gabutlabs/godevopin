package database

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gabutlabs/godevopin/internal/config"
	"github.com/gabutlabs/godevopin/internal/model"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// Connections contains the independent SQLite stores used by the application.
// App stores configuration and operational data, Metrics stores time-series
// samples, and Logs stores parsed log entries.
type Connections struct {
	App     *gorm.DB
	Metrics *gorm.DB
	Logs    *gorm.DB
}

// ConnectDB opens (and, when necessary, creates) all application databases.
func ConnectDB(cfg config.DBConfig) (*Connections, error) {
	directory := cfg.Directory
	if directory == "" {
		directory = "./data"
	}

	appPath := resolvePath(directory, cfg.MainPath, "godevopin.db")
	metricsPath := resolvePath(directory, cfg.MetricsPath, "metrics.db")
	logsPath := resolvePath(directory, cfg.LogsPath, "logs.db")

	appDB, err := openSQLite(appPath)
	if err != nil {
		return nil, fmt.Errorf("open app database %q: %w", appPath, err)
	}
	metricsDB, err := openSQLite(metricsPath)
	if err != nil {
		return nil, fmt.Errorf("open metrics database %q: %w", metricsPath, err)
	}
	logsDB, err := openSQLite(logsPath)
	if err != nil {
		return nil, fmt.Errorf("open logs database %q: %w", logsPath, err)
	}

	fmt.Printf("SQLite databases ready: %s, %s, %s\n", appPath, metricsPath, logsPath)
	return &Connections{App: appDB, Metrics: metricsDB, Logs: logsDB}, nil
}

// AutoMigrate creates only the tables belonging to each store. SQLite creates
// the database files themselves when they are opened above.
func (c *Connections) AutoMigrate() error {
	if err := c.App.AutoMigrate(
		&model.User{},
		&model.WorkerService{},
		&model.ActiveAlarm{},
		&model.AlarmHistory{},
		&model.Project{},
		&model.AppSetting{},
	); err != nil {
		return fmt.Errorf("migrate app database: %w", err)
	}

	if err := c.Metrics.AutoMigrate(&model.SystemMetric{}, &model.ProcessMetric{}); err != nil {
		return fmt.Errorf("migrate metrics database: %w", err)
	}

	if err := c.Logs.AutoMigrate(&model.LogHistory{}); err != nil {
		return fmt.Errorf("migrate logs database: %w", err)
	}

	return nil
}

func resolvePath(directory, configuredPath, defaultName string) string {
	if configuredPath != "" {
		return configuredPath
	}
	return filepath.Join(directory, defaultName)
}

func openSQLite(path string) (*gorm.DB, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("create database directory: %w", err)
	}

	// glebarez/sqlite uses modernc.org/sqlite, so this connection is pure Go
	// and works with CGO_ENABLED=0.
	db, err := gorm.Open(sqlite.Open(path+"?_busy_timeout=5000&_journal_mode=WAL"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	// A single writer avoids SQLITE_BUSY errors while workers and HTTP handlers
	// write to the same local database.
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)
	return db, nil
}
