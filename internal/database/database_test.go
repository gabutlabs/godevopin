package database

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gabutlabs/godevopin/internal/config"
)

func TestConnectDBCreatesAndMigratesIndependentStores(t *testing.T) {
	directory := t.TempDir()
	dbs, err := ConnectDB(config.DBConfig{Directory: directory})
	if err != nil {
		t.Fatalf("ConnectDB() error = %v", err)
	}

	if err := dbs.AutoMigrate(); err != nil {
		t.Fatalf("AutoMigrate() error = %v", err)
	}

	for _, name := range []string{"godevopin.db", "metrics.db", "logs.db"} {
		path, err := filepath.Abs(filepath.Join(directory, name))
		if err != nil {
			t.Fatalf("could not resolve database path %q: %v", name, err)
		}
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("expected database file %q: %v", path, err)
		}
	}

	assertTable := func(dbName string, table string) {
		t.Helper()
		var count int64
		var db = dbs.App
		switch dbName {
		case "metrics":
			db = dbs.Metrics
		case "logs":
			db = dbs.Logs
		}
		if err := db.Raw("SELECT count(*) FROM sqlite_master WHERE type = 'table' AND name = ?", table).Scan(&count).Error; err != nil {
			t.Fatalf("inspect %s.%s: %v", dbName, table, err)
		}
		if count != 1 {
			t.Fatalf("expected %s.%s to exist, got count %d", dbName, table, count)
		}
	}

	assertTable("app", "users")
	assertTable("app", "projects")
	assertTable("metrics", "system_metrics")
	assertTable("metrics", "process_metrics")
	assertTable("logs", "log_histories")

	assertMissingTable := func(dbName string, table string) {
		t.Helper()
		var count int64
		db := dbs.App
		switch dbName {
		case "metrics":
			db = dbs.Metrics
		case "logs":
			db = dbs.Logs
		}
		if err := db.Raw("SELECT count(*) FROM sqlite_master WHERE type = 'table' AND name = ?", table).Scan(&count).Error; err != nil {
			t.Fatalf("inspect %s.%s: %v", dbName, table, err)
		}
		if count != 0 {
			t.Fatalf("expected %s.%s to be isolated, got count %d", dbName, table, count)
		}
	}

	assertMissingTable("app", "system_metrics")
	assertMissingTable("metrics", "users")
	assertMissingTable("logs", "projects")
}
