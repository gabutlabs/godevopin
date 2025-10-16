package repository

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gabutlabs/devopin/internal/model"
	"gorm.io/gorm"
)

type ResultSystemMetric struct {
	TimeInterval string  `json:"time_interval"`
	AvgCPUUsage  float64 `json:"avg_cpu_usage"`
	AvgMemUsage  float64 `json:"avg_mem_usage"`
}

type ResultDiskUsage struct {
	DiskUsageByte float64 `json:"disk_usage_byte"`
	DiskTotalByte float64 `json:"disk_total_byte"`
}

type SystemMetricRepository interface {
	CreateMetric(metric *model.SystemMetric) error
	GetMetricByID(id uint) (*model.SystemMetric, error)
	UpdateMetric(metric *model.SystemMetric) error
	DeleteMetric(id uint) error
	FilterSystemMetrics(filter string) ([]ResultSystemMetric, error)
	GetDiskUsage() (ResultDiskUsage, error)
}
type systemMetricRepository struct {
	// Tambahkan field yang diperlukan, misalnya koneksi database
	db *gorm.DB
}

// NewSystemMetricRepository membuat instance baru dari systemMetricRepository
func NewSystemMetricRepository(db *gorm.DB) SystemMetricRepository {
	return &systemMetricRepository{db: db}
}

// Implementasikan metode CRUD sesuai kebutuhan
// Contoh: CreateMetric, GetMetricByID, UpdateMetric, DeleteMetric, dll.
// Contoh metode untuk mendapatkan metric berdasarkan ID
func (r *systemMetricRepository) GetMetricByID(id uint) (*model.SystemMetric, error) {
	var metric model.SystemMetric
	if err := r.db.First(&metric, id).Error; err != nil {
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

func (r *systemMetricRepository) DeleteMetric(id uint) error {
	return r.db.Delete(&model.SystemMetric{}, id).Error
}

func (r *systemMetricRepository) GetDiskUsage() (ResultDiskUsage, error) {
	var diskUsage ResultDiskUsage
	if err := r.db.Raw("select disk_usage_byte,disk_total_byte from system_metrics order by created_at desc").Scan(&diskUsage).Error; err != nil {
		return ResultDiskUsage{DiskUsageByte: 0, DiskTotalByte: 0}, err
	}
	return diskUsage, nil
}

func (r *systemMetricRepository) FilterSystemMetrics(filter string) ([]ResultSystemMetric, error) {
	var metrics []ResultSystemMetric
	query := ""
	if strings.HasSuffix(filter, "h") {
		filterHour := "1 hour"
		switch filter {
		case "1h":
			filterHour = "1 hour"
		case "6h":
			filterHour = "6 hours"
		case "12h":
			filterHour = "12 hours"
		default:
			return nil, errors.New("invalid filter format")
		}
		query = fmt.Sprintf(`
				SELECT
					-- Kelompokkan waktu ke dalam interval 15 detik
					time_bucket('15 seconds', created_at) AS time_interval,
					AVG(cpu_usage) AS avg_cpu_usage,
					AVG(mem_usage_byte) AS avg_mem_usage
					FROM
					system_metrics -- <--- Mengambil dari tabel data MENTAH
					WHERE
					created_at > NOW() - interval '%s' -- <--- Filter waktu Anda
					GROUP BY
					time_interval
					ORDER BY
					time_interval;
				`, filterHour)
	} else if strings.HasSuffix(filter, "d") {
		filterDay := "1 day"
		switch filter {
		case "1d":
			filterDay = "1 day"
		case "7d":
			filterDay = "7 days"
		case "30d":
			filterDay = "30 days"
		default:
			return nil, errors.New("invalid filter format")
		}
		if filter == "1d" {
			query = fmt.Sprintf(`
				SELECT
					hour as time_interval,
					avg_cpu_usage,
					avg_mem_usage
					FROM
					system_metrics_hourly -- <--- Mengambil dari tabel data MENTAH
					WHERE
					hour > NOW() - interval '%s' -- <--- Filter waktu Anda
					GROUP BY
					hour,avg_cpu_usage,avg_mem_usage
					ORDER BY
					hour;
				`, filterDay)
		} else {
			query = fmt.Sprintf(`
				SELECT days as time_interval,
					avg_cpu_usage,
					avg_mem_usage from system_metrics_daily
				WHERE
				days > NOW() - interval '%s' -- <--- Filter waktu Anda
				ORDER BY
				days;
				`, filterDay)
		}

	}
	if query == "" {
		return nil, errors.New("invalid filter format")
	}
	if err := r.db.Raw(query).Scan(&metrics).Error; err != nil {
		return nil, err
	}
	return metrics, nil
}
