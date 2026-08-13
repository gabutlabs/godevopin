package repository

import (
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/gabutlabs/godevopin/internal/model"
	"gorm.io/gorm"
)

type ProcessHistoryPoint struct {
	TimeInterval     string  `json:"time_interval"`
	AvgCPUPercent    float64 `json:"avg_cpu_percent"`
	AvgMemoryPercent float64 `json:"avg_memory_percent"`
	AvgMemoryBytes   float64 `json:"avg_memory_bytes"`
	MaxCPUPercent    float64 `json:"max_cpu_percent"`
}

type ProcessMetricRepository interface {
	CreateProcessMetrics(metrics []model.ProcessMetric) error
	DeleteOlderProcessMetrics(cutoff time.Time) (int64, error)
	GetProcessHistory(processKey string, duration, bucket time.Duration) ([]ProcessHistoryPoint, error)
}

type processMetricRepository struct {
	db *gorm.DB
}

func NewProcessMetricRepository(db *gorm.DB) ProcessMetricRepository {
	return &processMetricRepository{db: db}
}

func (r *processMetricRepository) CreateProcessMetrics(metrics []model.ProcessMetric) error {
	if len(metrics) == 0 {
		return nil
	}
	return r.db.Create(&metrics).Error
}

func (r *processMetricRepository) DeleteOlderProcessMetrics(cutoff time.Time) (int64, error) {
	result := r.db.Where("observed_at < ?", cutoff.UTC()).Delete(&model.ProcessMetric{})
	return result.RowsAffected, result.Error
}

func (r *processMetricRepository) GetProcessHistory(processKey string, duration, bucket time.Duration) ([]ProcessHistoryPoint, error) {
	if strings.TrimSpace(processKey) == "" {
		return nil, errors.New("process key is required")
	}

	cutoff := time.Now().UTC().Add(-duration)
	var metrics []model.ProcessMetric
	if err := r.db.Where("process_key = ? AND observed_at >= ?", processKey, cutoff).
		Order("observed_at ASC").Find(&metrics).Error; err != nil {
		return nil, err
	}

	type aggregate struct {
		time       time.Time
		cpu        float64
		memory     float64
		memoryByte float64
		maxCPU     float64
		count      int
	}
	aggregates := make(map[int64]*aggregate)
	for _, metric := range metrics {
		bucketTime := metric.ObservedAt.UTC().Truncate(bucket)
		key := bucketTime.UnixNano()
		if aggregates[key] == nil {
			aggregates[key] = &aggregate{time: bucketTime}
		}
		item := aggregates[key]
		item.cpu += metric.CPUPercent
		item.memory += metric.MemoryPercent
		item.memoryByte += float64(metric.MemoryBytes)
		if metric.CPUPercent > item.maxCPU {
			item.maxCPU = metric.CPUPercent
		}
		item.count++
	}

	ordered := make([]aggregate, 0, len(aggregates))
	for _, item := range aggregates {
		ordered = append(ordered, *item)
	}
	sort.Slice(ordered, func(i, j int) bool {
		return ordered[i].time.Before(ordered[j].time)
	})

	result := make([]ProcessHistoryPoint, 0, len(ordered))
	for _, item := range ordered {
		count := float64(item.count)
		result = append(result, ProcessHistoryPoint{
			TimeInterval:     item.time.UTC().Format(time.RFC3339),
			AvgCPUPercent:    item.cpu / count,
			AvgMemoryPercent: item.memory / count,
			AvgMemoryBytes:   item.memoryByte / count,
			MaxCPUPercent:    item.maxCPU,
		})
	}
	return result, nil
}

func ProcessMetricRange(filter string) (time.Duration, time.Duration, error) {
	switch strings.ToLower(strings.TrimSpace(filter)) {
	case "1h", "1 hour":
		return time.Hour, 15 * time.Second, nil
	case "6h", "6 hours":
		return 6 * time.Hour, time.Minute, nil
	case "12h", "12 hours":
		return 12 * time.Hour, 2 * time.Minute, nil
	case "1d", "1 day":
		return 24 * time.Hour, time.Hour, nil
	case "7d", "7 days":
		return 7 * 24 * time.Hour, 6 * time.Hour, nil
	case "30d", "30 days":
		return 30 * 24 * time.Hour, 24 * time.Hour, nil
	default:
		return 0, 0, errors.New("invalid filter format")
	}
}
