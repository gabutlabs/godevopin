package service

import (
	"strconv"
	"time"

	"github.com/gabutlabs/devopin/internal/model"
	"github.com/gabutlabs/devopin/internal/monitoring"
	"github.com/gabutlabs/devopin/internal/repository"
)

type SystemMetricService interface {
	GetSystemMetrics() (map[string]any, error)
	GetFilteredSystemMetrics(filter string) ([]repository.ResultSystemMetric, error)
	GetDiskUsage() (repository.ResultDiskUsage, error)
}

type systemMetricService struct {
	// Tambahkan dependensi yang diperlukan, misalnya database connection
	repository repository.SystemMetricRepository
}

func NewSystemMetricService(repo repository.SystemMetricRepository) SystemMetricService {
	return &systemMetricService{repository: repo}
}

// GetDiskUsage implements SystemMetricService.
func (s *systemMetricService) GetDiskUsage() (repository.ResultDiskUsage, error) {
	diskUsage, err := s.repository.GetDiskUsage()
	if err != nil {
		return diskUsage, err
	}
	return diskUsage, nil
}

// GetSystemMetrics implements SystemMetricService.
func (s *systemMetricService) GetSystemMetrics() (map[string]any, error) {
	cpuUsage, err := monitoring.NewCollector().GetCPUUsage()
	if err != nil {
		return nil, err
	}
	memoryUsage, err := monitoring.NewCollector().GetMemoryUsage()
	if err != nil {
		return nil, err
	}
	net, err := monitoring.NewCollector().GetNetworkUsage()
	if err != nil {
		return nil, err
	}
	disk, err := monitoring.NewCollector().GetDiskUsage("/")
	if err != nil {
		return nil, err
	}
	s.repository.CreateMetric(&model.SystemMetric{
		ID:            strconv.FormatInt(time.Now().UTC().Unix(), 10),
		CPUUsage:      cpuUsage.UsagePercent,
		MemUsageByte:  float64(memoryUsage.Used),
		MemTotalByte:  float64(memoryUsage.Total),
		DiskUsageByte: float64(disk.Used),
		DiskTotalByte: float64(disk.Total),
	})
	return map[string]any{
		"cpu_usage":     cpuUsage,
		"memory_usage":  memoryUsage,
		"network_usage": net,
		"disk_usage":    disk,
	}, nil
}

func (s *systemMetricService) GetFilteredSystemMetrics(filter string) ([]repository.ResultSystemMetric, error) {
	return s.repository.FilterSystemMetrics(filter)
}
