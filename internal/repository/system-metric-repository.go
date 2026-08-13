package repository

import (
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/gabutlabs/godevopin/internal/model"
	"gorm.io/gorm"
)

type ResultSystemMetric struct {
	TimeInterval string  `json:"time_interval"`
	AvgCPUUsage  float64 `json:"avg_cpu_usage"`
	AvgMemUsage  float64 `json:"avg_mem_usage"`
}

type LastSystemMetric struct {
	TimeInterval string  `json:"time_interval"`
	AvgCPUUsage  float64 `json:"avg_cpu_usage"`
	AvgMemUsage  float64 `json:"avg_mem_usage"`
	AvgDiskUsage float64 `json:"avg_disk_usage"`
	AvgDiskTotal float64 `json:"avg_disk_total"`
}

type ResultDiskUsage struct {
	DiskUsageByte float64 `json:"disk_usage_byte"`
	DiskTotalByte float64 `json:"disk_total_byte"`
}

type SystemMetricRepository interface {
	CreateMetric(metric *model.SystemMetric) error
	GetMetricByID(id string) (*model.SystemMetric, error)
	UpdateMetric(metric *model.SystemMetric) error
	DeleteMetric(id string) error
	FilterSystemMetrics(filter string) ([]ResultSystemMetric, error)
	LastSystemMetricsFilter(filter string) (LastSystemMetric, error)
	GetDiskUsage() (ResultDiskUsage, error)
}

type systemMetricRepository struct {
	db *gorm.DB
}

func NewSystemMetricRepository(db *gorm.DB) SystemMetricRepository {
	return &systemMetricRepository{db: db}
}

func (r *systemMetricRepository) GetMetricByID(id string) (*model.SystemMetric, error) {
	var metric model.SystemMetric
	if err := r.db.First(&metric, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &metric, nil
}

func (r *systemMetricRepository) CreateMetric(metric *model.SystemMetric) error {
	return r.db.Create(metric).Error
}

func (r *systemMetricRepository) UpdateMetric(metric *model.SystemMetric) error {
	return r.db.Save(metric).Error
}

func (r *systemMetricRepository) DeleteMetric(id string) error {
	return r.db.Delete(&model.SystemMetric{}, "id = ?", id).Error
}

func (r *systemMetricRepository) GetDiskUsage() (ResultDiskUsage, error) {
	var metric model.SystemMetric
	err := r.db.Order("created_at DESC").First(&metric).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ResultDiskUsage{}, nil
	}
	if err != nil {
		return ResultDiskUsage{}, err
	}
	return ResultDiskUsage{
		DiskUsageByte: metric.DiskUsageByte,
		DiskTotalByte: metric.DiskTotalByte,
	}, nil
}

// FilterSystemMetrics performs time-window aggregation in Go so the repository
// remains independent of database-specific date and bucketing functions.
func (r *systemMetricRepository) FilterSystemMetrics(filter string) ([]ResultSystemMetric, error) {
	aggregates, err := r.aggregateMetrics(filter)
	if err != nil {
		return nil, err
	}

	result := make([]ResultSystemMetric, 0, len(aggregates))
	for _, aggregate := range aggregates {
		result = append(result, ResultSystemMetric{
			TimeInterval: aggregate.time.UTC().Format(time.RFC3339),
			AvgCPUUsage:  aggregate.cpu / float64(aggregate.count),
			AvgMemUsage:  aggregate.mem / float64(aggregate.count),
		})
	}
	return result, nil
}

func (r *systemMetricRepository) LastSystemMetricsFilter(filter string) (LastSystemMetric, error) {
	aggregates, err := r.aggregateMetrics(filter)
	if err != nil {
		return LastSystemMetric{}, err
	}
	if len(aggregates) == 0 {
		return LastSystemMetric{}, nil
	}

	aggregate := aggregates[len(aggregates)-1]
	count := float64(aggregate.count)
	return LastSystemMetric{
		TimeInterval: aggregate.time.UTC().Format(time.RFC3339),
		AvgCPUUsage:  aggregate.cpu / count,
		AvgMemUsage:  aggregate.mem / count,
		AvgDiskUsage: aggregate.disk / count,
		AvgDiskTotal: aggregate.diskTotal / count,
	}, nil
}

type metricAggregate struct {
	time      time.Time
	cpu       float64
	mem       float64
	disk      float64
	diskTotal float64
	count     int
}

func (r *systemMetricRepository) aggregateMetrics(filter string) ([]metricAggregate, error) {
	duration, bucket, err := metricRange(filter)
	if err != nil {
		return nil, err
	}

	var metrics []model.SystemMetric
	cutoff := time.Now().UTC().Add(-duration)
	if err := r.db.Where("created_at >= ?", cutoff).Order("created_at ASC").Find(&metrics).Error; err != nil {
		return nil, err
	}

	byBucket := make(map[int64]*metricAggregate)
	for _, metric := range metrics {
		bucketTime := metric.CreatedAt.UTC().Truncate(bucket)
		key := bucketTime.UnixNano()
		if byBucket[key] == nil {
			byBucket[key] = &metricAggregate{time: bucketTime}
		}
		aggregate := byBucket[key]
		aggregate.cpu += metric.CPUUsage
		aggregate.mem += metric.MemUsageByte
		aggregate.disk += metric.DiskUsageByte
		aggregate.diskTotal += metric.DiskTotalByte
		aggregate.count++
	}

	result := make([]metricAggregate, 0, len(byBucket))
	for _, aggregate := range byBucket {
		result = append(result, *aggregate)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].time.Before(result[j].time)
	})
	return result, nil
}

func metricRange(filter string) (time.Duration, time.Duration, error) {
	switch strings.ToLower(strings.TrimSpace(filter)) {
	case "1h", "1 hour":
		return time.Hour, 15 * time.Second, nil
	case "6h", "6 hours":
		return 6 * time.Hour, 15 * time.Second, nil
	case "12h", "12 hours":
		return 12 * time.Hour, 15 * time.Second, nil
	case "1d", "1 day":
		return 24 * time.Hour, time.Hour, nil
	case "7d", "7 days":
		return 7 * 24 * time.Hour, 24 * time.Hour, nil
	case "30d", "30 days":
		return 30 * 24 * time.Hour, 24 * time.Hour, nil
	default:
		return 0, 0, errors.New("invalid filter format")
	}
}
