package service

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gabutlabs/godevopin/internal/model"
	"github.com/gabutlabs/godevopin/internal/monitoring"
	"github.com/gabutlabs/godevopin/internal/repository"
)

const (
	ProcessMetricRetention = 30 * 24 * time.Hour
	processLiveCacheTTL    = 2 * time.Second
	processHistoryLimit    = 100
)

type ProcessMonitoringService interface {
	GetLiveProcesses(sortBy, search string, limit int) ([]model.ProcessMetric, error)
	CollectAndPersist() error
	GetProcessHistory(processKey, filter string) ([]repository.ProcessHistoryPoint, error)
	DeleteExpiredProcessMetrics() (int64, error)
}

type processMonitoringService struct {
	repository repository.ProcessMetricRepository
	collector  *monitoring.ProcessCollector

	mu            sync.Mutex
	lastSnapshot  []model.ProcessMetric
	lastCollected time.Time
}

func NewProcessMonitoringService(repo repository.ProcessMetricRepository) ProcessMonitoringService {
	return &processMonitoringService{
		repository: repo,
		collector:  monitoring.NewProcessCollector(),
	}
}

func (s *processMonitoringService) GetLiveProcesses(sortBy, search string, limit int) ([]model.ProcessMetric, error) {
	snapshot, err := s.snapshot(false)
	if err != nil {
		return nil, err
	}

	search = strings.ToLower(strings.TrimSpace(search))
	filtered := make([]model.ProcessMetric, 0, len(snapshot))
	for _, process := range snapshot {
		if search != "" && !strings.Contains(strings.ToLower(process.Name), search) &&
			!strings.Contains(strings.ToLower(process.Username), search) &&
			!strings.Contains(fmt.Sprint(process.PID), search) {
			continue
		}
		filtered = append(filtered, process)
	}

	sortProcesses(filtered, sortBy)
	if limit <= 0 {
		limit = 200
	}
	if limit > 500 {
		limit = 500
	}
	if len(filtered) > limit {
		filtered = filtered[:limit]
	}
	return filtered, nil
}

func (s *processMonitoringService) CollectAndPersist() error {
	snapshot, err := s.snapshot(true)
	if err != nil {
		return err
	}

	selected := selectHistoricalProcesses(snapshot, processHistoryLimit)
	now := time.Now().UTC()
	for i := range selected {
		selected[i].ID = fmt.Sprintf("%s:%d", selected[i].ProcessKey, now.UnixNano())
		selected[i].ObservedAt = now
	}
	return s.repository.CreateProcessMetrics(selected)
}

func (s *processMonitoringService) GetProcessHistory(processKey, filter string) ([]repository.ProcessHistoryPoint, error) {
	if strings.TrimSpace(processKey) == "" {
		return nil, errors.New("process key is required")
	}
	duration, bucket, err := repository.ProcessMetricRange(filter)
	if err != nil {
		return nil, err
	}
	return s.repository.GetProcessHistory(processKey, duration, bucket)
}

func (s *processMonitoringService) DeleteExpiredProcessMetrics() (int64, error) {
	return s.repository.DeleteOlderProcessMetrics(time.Now().UTC().Add(-ProcessMetricRetention))
}

func (s *processMonitoringService) snapshot(force bool) ([]model.ProcessMetric, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !force && len(s.lastSnapshot) > 0 && time.Since(s.lastCollected) < processLiveCacheTTL {
		return cloneProcessMetrics(s.lastSnapshot), nil
	}

	snapshot, err := s.collector.Collect()
	if err != nil {
		return nil, err
	}
	s.lastSnapshot = snapshot
	s.lastCollected = time.Now()
	return cloneProcessMetrics(snapshot), nil
}

func sortProcesses(processes []model.ProcessMetric, sortBy string) {
	switch strings.ToLower(strings.TrimSpace(sortBy)) {
	case "memory", "memory_percent":
		sort.SliceStable(processes, func(i, j int) bool { return processes[i].MemoryPercent > processes[j].MemoryPercent })
	case "pid":
		sort.SliceStable(processes, func(i, j int) bool { return processes[i].PID < processes[j].PID })
	case "name":
		sort.SliceStable(processes, func(i, j int) bool { return processes[i].Name < processes[j].Name })
	case "threads", "thread_count":
		sort.SliceStable(processes, func(i, j int) bool { return processes[i].ThreadCount > processes[j].ThreadCount })
	default:
		sort.SliceStable(processes, func(i, j int) bool { return processes[i].CPUPercent > processes[j].CPUPercent })
	}
}

func selectHistoricalProcesses(snapshot []model.ProcessMetric, limit int) []model.ProcessMetric {
	if len(snapshot) <= limit {
		return cloneProcessMetrics(snapshot)
	}

	byCPU := cloneProcessMetrics(snapshot)
	byMemory := cloneProcessMetrics(snapshot)
	sortProcesses(byCPU, "cpu")
	sortProcesses(byMemory, "memory")

	selected := make(map[string]model.ProcessMetric, limit)
	half := limit / 2
	for _, process := range byCPU[:half] {
		selected[process.ProcessKey] = process
	}
	for _, process := range byMemory[:half] {
		if len(selected) >= limit {
			break
		}
		selected[process.ProcessKey] = process
	}
	for _, process := range byCPU[half:] {
		if len(selected) >= limit {
			break
		}
		selected[process.ProcessKey] = process
	}

	result := make([]model.ProcessMetric, 0, len(selected))
	for _, process := range selected {
		result = append(result, process)
	}
	return result
}

func cloneProcessMetrics(source []model.ProcessMetric) []model.ProcessMetric {
	return append([]model.ProcessMetric(nil), source...)
}
