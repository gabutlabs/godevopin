package repository

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/gabutlabs/godevopin/internal/model"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestProcessMetricRepositoryAggregatesHistory(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(filepath.Join(t.TempDir(), "metrics.db")), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&model.ProcessMetric{}); err != nil {
		t.Fatalf("migrate process metrics: %v", err)
	}

	now := time.Now().UTC().Truncate(time.Second)
	metrics := []model.ProcessMetric{
		{ID: "proc-1", ProcessKey: "42-1000", PID: 42, CPUPercent: 20, MemoryPercent: 10, MemoryBytes: 100, ObservedAt: now.Add(-2 * time.Minute)},
		{ID: "proc-2", ProcessKey: "42-1000", PID: 42, CPUPercent: 40, MemoryPercent: 30, MemoryBytes: 300, ObservedAt: now.Add(-time.Minute)},
		{ID: "proc-3", ProcessKey: "99-1000", PID: 99, CPUPercent: 90, ObservedAt: now.Add(-time.Minute)},
	}
	if err := db.Create(&metrics).Error; err != nil {
		t.Fatalf("insert process metrics: %v", err)
	}

	repo := NewProcessMetricRepository(db)
	history, err := repo.GetProcessHistory("42-1000", time.Hour, time.Minute)
	if err != nil {
		t.Fatalf("GetProcessHistory() error = %v", err)
	}
	if len(history) != 2 {
		t.Fatalf("GetProcessHistory() returned %d buckets, want 2", len(history))
	}
	if history[0].AvgCPUPercent != 20 || history[1].AvgCPUPercent != 40 {
		t.Fatalf("unexpected history: %+v", history)
	}
}

func TestProcessMetricRepositoryDeletesOnlyExpiredRecords(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(filepath.Join(t.TempDir(), "metrics.db")), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&model.ProcessMetric{}); err != nil {
		t.Fatalf("migrate process metrics: %v", err)
	}

	now := time.Now().UTC()
	metrics := []model.ProcessMetric{
		{ID: "expired", ProcessKey: "1-1000", PID: 1, ObservedAt: now.Add(-31 * 24 * time.Hour)},
		{ID: "retained", ProcessKey: "2-1000", PID: 2, ObservedAt: now.Add(-29 * 24 * time.Hour)},
	}
	if err := db.Create(&metrics).Error; err != nil {
		t.Fatalf("insert process metrics: %v", err)
	}

	repo := NewProcessMetricRepository(db)
	deleted, err := repo.DeleteOlderProcessMetrics(now.Add(-30 * 24 * time.Hour))
	if err != nil {
		t.Fatalf("DeleteOlderProcessMetrics() error = %v", err)
	}
	if deleted != 1 {
		t.Fatalf("DeleteOlderProcessMetrics() deleted %d records, want 1", deleted)
	}

	var remaining int64
	if err := db.Model(&model.ProcessMetric{}).Count(&remaining).Error; err != nil {
		t.Fatalf("count remaining process metrics: %v", err)
	}
	if remaining != 1 {
		t.Fatalf("remaining process metrics = %d, want 1", remaining)
	}
}
