package worker

import service "github.com/gabutlabs/devopin/internal/services"

type MonitoringWorker struct {
	service service.SystemMetricService
}

func NewMonitoringWorker(svc service.SystemMetricService) *MonitoringWorker {
	return &MonitoringWorker{service: svc}
}

func (m *MonitoringWorker) StartMonitoring() {
	m.service.GetSystemMetrics()
}
