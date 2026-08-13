package repository

import (
	"testing"
	"time"

	"github.com/gabutlabs/godevopin/internal/model"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestSystemMetricRepositoryAggregatesSQLiteData(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&model.SystemMetric{}); err != nil {
		t.Fatalf("migrate metrics: %v", err)
	}

	now := time.Now().UTC().Truncate(time.Second)
	metrics := []model.SystemMetric{
		{ID: "metric-1", CPUUsage: 20, MemUsageByte: 100, DiskUsageByte: 40, DiskTotalByte: 100, CreatedAt: now.Add(-2 * time.Minute)},
		{ID: "metric-2", CPUUsage: 40, MemUsageByte: 200, DiskUsageByte: 60, DiskTotalByte: 100, CreatedAt: now.Add(-time.Minute)},
		{ID: "metric-3", CPUUsage: 80, MemUsageByte: 300, DiskUsageByte: 70, DiskTotalByte: 100, CreatedAt: now},
	}
	if err := db.Create(&metrics).Error; err != nil {
		t.Fatalf("insert metrics: %v", err)
	}

	repo := NewSystemMetricRepository(db)
	filtered, err := repo.FilterSystemMetrics("1h")
	if err != nil {
		t.Fatalf("FilterSystemMetrics() error = %v", err)
	}
	if len(filtered) == 0 {
		t.Fatal("FilterSystemMetrics() returned no buckets")
	}

	last, err := repo.LastSystemMetricsFilter("1 hour")
	if err != nil {
		t.Fatalf("LastSystemMetricsFilter() error = %v", err)
	}
	if last.AvgCPUUsage != 80 || last.AvgDiskUsage != 70 || last.AvgDiskTotal != 100 {
		t.Fatalf("unexpected last metric: %+v", last)
	}
}
