package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"time"

	"github.com/coreos/go-systemd/v22/dbus"
	"github.com/gabutlabs/devopin/internal/model"
	"github.com/gabutlabs/devopin/internal/repository"
)

type WorkerServiceService interface {
	CreateWorkerService(name, description string, desiredState model.DesiredState) error
	GetWorkerServiceByID(id uint) (*model.WorkerService, error)
	GetWorkerServiceByName(name string) (*model.WorkerService, error)
	UpdateWorkerService(workerService *model.WorkerService) error
	DeleteWorkerService(id uint) error
	ListWorkerServices() ([]model.WorkerService, error)
	UpdateWorkerServiceStatus(id uint, currentStatus model.CurrentStatus, healthStatus model.HealthStatus) error
	UpdateWorkerServiceHeartbeat(id uint) error
}

type workerServiceService struct {
	repository repository.WorkerServiceRepository
}

func NewWorkerServiceService(repo repository.WorkerServiceRepository) WorkerServiceService {
	return &workerServiceService{repository: repo}
}

// CreateWorkerService creates a new WorkerService
func (s *workerServiceService) CreateWorkerService(name, description string, desiredState model.DesiredState) error {
	// Check if a WorkerService with this name already exists
	existingWorkerService, err := s.repository.GetWorkerServiceByName(name)
	if err == nil && existingWorkerService != nil {
		return errors.New("worker service with this name already exists")
	}

	// Create the new WorkerService
	workerService := &model.WorkerService{
		Name:         name,
		Description:  description,
		DesiredState: desiredState,
		// Initialize status values
		CurrentStatus: model.StatusStopped,
		HealthStatus:  model.HealthUnknown,
	}

	return s.repository.CreateWorkerService(workerService)
}

// GetWorkerServiceByID retrieves a WorkerService by its ID
func (s *workerServiceService) GetWorkerServiceByID(id uint) (*model.WorkerService, error) {
	return s.repository.GetWorkerServiceByID(id)
}

// GetWorkerServiceByName retrieves a WorkerService by its name
func (s *workerServiceService) GetWorkerServiceByName(name string) (*model.WorkerService, error) {
	return s.repository.GetWorkerServiceByName(name)
}

// UpdateWorkerService updates an existing WorkerService
func (s *workerServiceService) UpdateWorkerService(workerService *model.WorkerService) error {
	return s.repository.UpdateWorkerService(workerService)
}

// DeleteWorkerService deletes a WorkerService by its ID
func (s *workerServiceService) DeleteWorkerService(id uint) error {
	return s.repository.DeleteWorkerService(id)
}

// ListWorkerServices retrieves all WorkerServices
func (s *workerServiceService) ListWorkerServices() ([]model.WorkerService, error) {
	return s.repository.ListWorkerServices()
}

// UpdateWorkerServiceStatus updates the status of a WorkerService
func (s *workerServiceService) UpdateWorkerServiceStatus(id uint, currentStatus model.CurrentStatus, healthStatus model.HealthStatus) error {
	workerService, err := s.GetWorkerServiceByID(id)
	if err != nil {
		return err
	}

	workerService.CurrentStatus = currentStatus
	workerService.HealthStatus = healthStatus

	// Update the last status change time based on the status
	switch currentStatus {
	case model.StatusRunning:
		now := time.Now()
		workerService.LastSuccessAt = &now
	case model.StatusFailed:
		now := time.Now()
		workerService.LastFailureAt = &now
	}

	// Note: syncing host services should be handled by a component that does not cause an import cycle.
	// TODO: invoke sync logic via repository callbacks or an external orchestrator to avoid importing the worker package here.
	_, err = s.actionHostServices(workerService.Name, currentStatus)
	if err != nil {
		return err
	}
	return s.UpdateWorkerService(workerService)
}

// UpdateWorkerServiceHeartbeat updates the heartbeat timestamp of a WorkerService
func (s *workerServiceService) UpdateWorkerServiceHeartbeat(id uint) error {
	workerService, err := s.GetWorkerServiceByID(id)
	if err != nil {
		return err
	}

	now := time.Now()
	workerService.LastHeartbeatAt = &now

	return s.UpdateWorkerService(workerService)
}

func (s *workerServiceService) actionHostServices(workerName string, status model.CurrentStatus) (string, error) {
	switch runtime.GOOS {
	case "linux":
		return s.actionLinuxServices(workerName, status)
	case "darwin":
		return s.actionDarwinServices(workerName, status)
	default:
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

func (s *workerServiceService) actionLinuxServices(workerName string, status model.CurrentStatus) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	conn, err := dbus.NewSystemConnectionContext(ctx)
	if err != nil {
		return "", fmt.Errorf("failed connect to D-Bus systemd: %w", err)
	}
	defer conn.Close()

	// Declare jobID and err variables outside the switch
	var jobID int
	// err is already declared above, we will reuse it

	// Channel to receive signals when D-Bus job is complete
	completeChangeStatus := make(chan string, 1)

	// Goroutine to ensure channel is closed and no memory leak occurs
	go func() {
		// Waiting for signal from D-Bus
		result := <-completeChangeStatus
		log.Printf("D-Bus job for unit %s completed with result: %s", workerName, result)
	}()

	// --- Main Logic: Running action based on status ---
	switch status {
	case model.StatusRestart:
		log.Printf("Attempting to restart unit: %s", workerName)
		// Use '=' to assign value to already existing variable
		jobID, err = conn.RestartUnitContext(
			ctx,
			workerName,
			"replace", // Standard mode for restart
			completeChangeStatus,
		)

	case model.StatusStopped:
		log.Printf("Attempting to stop unit: %s", workerName)
		// Use '=' to assign value to already existing variable
		jobID, err = conn.StopUnitContext(
			ctx,
			workerName,
			"replace", // Standard mode for stop
			completeChangeStatus,
		)
	case model.StatusStarting:
		jobID, err = conn.StartUnitContext(ctx, workerName, "replace", completeChangeStatus)

	default:
		// Handling unsupported status
		return "", fmt.Errorf("action '%s' is not supported", status)
	}

	// --- Error Handling & Waiting for Result ---
	// Check error AFTER switch, now 'err' contains D-Bus call result
	if err != nil {
		return "", fmt.Errorf("failed to start job '%s' for unit %s: %w", status, workerName, err)
	}

	// Wait for job to complete or timeout
	select {
	case result := <-completeChangeStatus:
		// D-Bus reports job completion
		return fmt.Sprintf("Action '%s' for unit '%s' (Job ID: %d) completed with result: %s", status, workerName, jobID, result), nil
	case <-ctx.Done():
		// Timeout from context we created at the beginning
		return "", fmt.Errorf("timeout waiting for job '%s' for unit '%s' to complete", status, workerName)
	}
}

// This is the implementation for darwin(macOS)
func (s *workerServiceService) actionDarwinServices(workerName string, status model.CurrentStatus) (string, error) {
	// Use context for timeout, important to prevent process hanging
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute) // Brew can be slow
	defer cancel()

	// Find path to 'brew'
	cmdPath, err := exec.LookPath("brew")
	if err != nil {
		return "", fmt.Errorf("'brew' command not found, make sure Homebrew is installed")
	}

	var args []string
	action := ""

	// Determine command-line arguments based on requested status
	switch status {
	case model.StatusRestart:
		action = "restart"
		args = []string{"services", "restart", workerName}
	case model.StatusStopped:
		action = "stop"
		args = []string{"services", "stop", workerName}
	case model.StatusStarting: // With brew, 'start' and 'run' do the same thing
		action = "start"
		args = []string{"services", "start", workerName}
	default:
		return "", fmt.Errorf("action '%s' is not supported on macOS", status)
	}

	log.Printf("Attempting to %s service '%s' via brew...", action, workerName)

	// Create command with context so it can be timed out
	cmd := exec.CommandContext(ctx, cmdPath, args...)

	// Prepare buffer to capture output
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run command and wait until finished
	err = cmd.Run()

	// Check error after command completes
	if err != nil {
		// Check if error was caused by context timeout
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("timeout waiting for action '%s' for service '%s' to complete", action, workerName)
		}
		// If other error, return message from stderr
		return "", fmt.Errorf("failed to run action '%s' for service '%s': %s", action, workerName, stderr.String())
	}

	// If successful, return message from stdout
	return fmt.Sprintf("Action '%s' for service '%s' successful: %s", action, workerName, stdout.String()), nil
}
