package worker

import (
	"log"

	service "github.com/gabutlabs/godevopin/internal/services"
)

type ProcessMonitoringWorker struct {
	service service.ProcessMonitoringService
}

func NewProcessMonitoringWorker(svc service.ProcessMonitoringService) *ProcessMonitoringWorker {
	return &ProcessMonitoringWorker{service: svc}
}

func (w *ProcessMonitoringWorker) CollectAndPersist() {
	if err := w.service.CollectAndPersist(); err != nil {
		log.Printf("process monitoring collection failed: %v", err)
	}
}

func (w *ProcessMonitoringWorker) CleanupRetention() {
	deleted, err := w.service.DeleteExpiredProcessMetrics()
	if err != nil {
		log.Printf("process metrics retention cleanup failed: %v", err)
		return
	}
	if deleted > 0 {
		log.Printf("deleted %d process metric records older than 30 days", deleted)
	}
}
