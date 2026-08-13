package worker

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gabutlabs/godevopin/internal/config"
	"github.com/gabutlabs/godevopin/internal/database"
	"github.com/gabutlabs/godevopin/internal/repository"
	service "github.com/gabutlabs/godevopin/internal/services"
)

// StartWorkers is a more descriptive name.
// This function starts all workers and blocks until there is a shutdown signal.
func StartWorkers(cfg *config.Config, dbs *database.Connections) {
	log.Println("Starting background workers...")

	// --- Initialize all services and workers ---
	// Monitoring Worker
	sysmetricRepo := repository.NewSystemMetricRepository(dbs.Metrics)
	sysmetricService := service.NewSystemMetricService(sysmetricRepo)
	monitWorker := NewMonitoringWorker(sysmetricService)

	processMetricRepo := repository.NewProcessMetricRepository(dbs.Metrics)
	processMonitoringService := service.NewProcessMonitoringService(processMetricRepo)
	processMonitoringWorker := NewProcessMonitoringWorker(processMonitoringService)

	// Host Syncer Worker
	workerRepo := repository.NewWorkerServiceRepository(dbs.App)
	workerService := service.NewWorkerServiceService(workerRepo)
	hostSyncer := NewSyncer(workerService)

	// Threshold automation worker
	alarmRepo := repository.NewAlarmRepository(dbs.App)
	alarmService := service.NewAlarmService(alarmRepo)

	settingRepo := repository.NewSettingRepository(dbs.App)
	settingService := service.NewSettingService(settingRepo)

	// Log Parser Worker
	projectRepo := repository.NewProjectRepository(dbs.App)
	projectService := service.NewProjectService(projectRepo)
	logHistoryRepo := repository.NewLogHistoryRepository(dbs.Logs)
	logHistoryService := service.NewLogHistoryService(logHistoryRepo)
	logParserWorker := NewLogParserWorker(projectService, logHistoryService)

	// --- Determine intervals from config ---
	settings, err := settingService.GetSettings()
	monitoringInterval := 10 * time.Second // fallback
	if err == nil {
		monitoringInterval = time.Duration(settings.MonitoringIntervalSeconds) * time.Second
	}

	thresholdAutomation := NewThresholdAutomation(workerService, sysmetricService, alarmService, settingService)

	syncInterval := 5 * time.Minute // For example, sync interval is set differently
	logParseInterval := 2 * time.Minute
	processMonitoringInterval := 30 * time.Second

	// Telegram AI Agent Worker
	dockerService := service.NewDockerService()
	telegramWorker := NewTelegramWorker(settingService, sysmetricService, workerService, dockerService, logHistoryService, projectService)

	// --- Run all workers as periodic goroutines ---
	runPeriodicTask(monitWorker.StartMonitoring, monitoringInterval)
	runPeriodicTask(processMonitoringWorker.CollectAndPersist, processMonitoringInterval)
	runPeriodicTask(processMonitoringWorker.CleanupRetention, time.Hour)
	runPeriodicTask(hostSyncer.SyncHostServices, syncInterval)
	runPeriodicTask(thresholdAutomation.AlarmWorker, monitoringInterval)
	runPeriodicTask(logParserWorker.Run, logParseInterval)
	// runPeriodicTask(dockerManagement.RunDocker, monitoringInterval)

	// Start Telegram Bot (Long Polling)
	go telegramWorker.Start()

	log.Println("All workers are running. Press Ctrl+C to shut down.")

	// --- Wait for shutdown signal ---
	// Create channel to receive signal from the OS
	quit := make(chan os.Signal, 1)
	// Tell the channel to capture SIGINT (Ctrl+C) and SIGTERM signals
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block here until signal is received
	<-quit

	// Program will end after signal is received
	log.Println("Shutting down workers...")
}

// runPeriodicTask is a helper function to run tasks periodically.
// This reduces code duplication.
func runPeriodicTask(task func(), interval time.Duration) {
	go func() {
		// Run once at startup so we don't have to wait for the first interval
		task()

		// Create ticker to run tasks periodically
		ticker := time.NewTicker(interval)
		defer ticker.Stop() // Ensure ticker is cleaned up

		for {
			// Wait for next tick
			<-ticker.C
			// Run task
			task()
		}
	}()
}
