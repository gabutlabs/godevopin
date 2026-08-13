package monitoring

import (
	"fmt"
	"runtime"
	"sort"
	"sync"
	"time"

	"github.com/gabutlabs/godevopin/internal/model"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/process"
)

type processCPUState struct {
	user       float64
	system     float64
	observedAt time.Time
}

// ProcessCollector gathers a best-effort snapshot of processes running on the
// host. Individual permission or process-exit errors are ignored so one
// inaccessible process cannot make the complete snapshot fail.
type ProcessCollector struct {
	mu       sync.Mutex
	previous map[string]processCPUState
}

func NewProcessCollector() *ProcessCollector {
	return &ProcessCollector{previous: make(map[string]processCPUState)}
}

func (c *ProcessCollector) Collect() ([]model.ProcessMetric, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	processes, err := process.Processes()
	if err != nil {
		return nil, fmt.Errorf("list processes: %w", err)
	}

	virtualMemory, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("read host memory: %w", err)
	}

	now := time.Now().UTC()
	current := make(map[string]processCPUState, len(processes))
	snapshot := make([]model.ProcessMetric, 0, len(processes))
	for _, proc := range processes {
		metric, cpuState, ok := c.collectProcess(proc, now, virtualMemory.Total)
		if !ok {
			continue
		}
		current[metric.ProcessKey] = cpuState
		snapshot = append(snapshot, metric)
	}
	c.previous = current

	sort.SliceStable(snapshot, func(i, j int) bool {
		return snapshot[i].CPUPercent > snapshot[j].CPUPercent
	})
	return snapshot, nil
}

func (c *ProcessCollector) collectProcess(proc *process.Process, now time.Time, totalMemory uint64) (model.ProcessMetric, processCPUState, bool) {
	createdAt, err := proc.CreateTime()
	if err != nil {
		return model.ProcessMetric{}, processCPUState{}, false
	}
	times, err := proc.Times()
	if err != nil {
		return model.ProcessMetric{}, processCPUState{}, false
	}

	processKey := fmt.Sprintf("%d-%d", proc.Pid, createdAt)
	cpuPercent := 0.0
	if previous, exists := c.previous[processKey]; exists {
		elapsed := now.Sub(previous.observedAt).Seconds()
		cpuDelta := (times.User - previous.user) + (times.System - previous.system)
		if elapsed > 0 && cpuDelta > 0 {
			// Match htop's per-logical-CPU presentation: a process can use
			// more than 100% when it runs on multiple CPUs.
			cpuPercent = (cpuDelta / elapsed) * 100
		}
	}

	metric := model.ProcessMetric{
		ProcessKey:       processKey,
		PID:              proc.Pid,
		ProcessStartTime: createdAt,
		CPUPercent:       cpuPercent,
		ObservedAt:       now,
	}

	if parentPID, err := proc.Ppid(); err == nil {
		metric.ParentPID = parentPID
	}
	if name, err := proc.Name(); err == nil {
		metric.Name = name
	} else {
		metric.Name = "unknown"
	}
	if executable, err := proc.Exe(); err == nil {
		metric.Executable = executable
	}
	if username, err := proc.Username(); err == nil {
		metric.Username = username
	}
	if statuses, err := proc.Status(); err == nil && len(statuses) > 0 {
		metric.Status = statuses[0]
	}
	if memoryInfo, err := proc.MemoryInfo(); err == nil && memoryInfo != nil {
		metric.MemoryBytes = memoryInfo.RSS
		metric.VirtualMemoryBytes = memoryInfo.VMS
		if totalMemory > 0 {
			metric.MemoryPercent = float64(memoryInfo.RSS) / float64(totalMemory) * 100
		}
	}
	if threadCount, err := proc.NumThreads(); err == nil {
		metric.ThreadCount = threadCount
	}

	return metric, processCPUState{
		user:       times.User,
		system:     times.System,
		observedAt: now,
	}, true
}

// LogicalCPUCount is exposed for tests and documentation of the CPU scale.
func LogicalCPUCount() int {
	return runtime.NumCPU()
}
