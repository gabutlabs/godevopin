package worker

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strings"

	"github.com/coreos/go-systemd/v22/dbus"
	"github.com/gabutlabs/devopin/internal/model"
	service "github.com/gabutlabs/devopin/internal/services"
)

// DiscoveredService is an internal struct to hold data from the OS.
type DiscoveredService struct {
	Name        string
	Description string
	Status      string // Raw status: "active", "started", "failed", "error", etc.
}

// Syncer is responsible for discovering and synchronizing services from host to DB.
type Syncer struct {
	workerService service.WorkerServiceService
}

// NewSyncer creates a new Syncer instance.
func NewSyncer(ws service.WorkerServiceService) *Syncer {
	return &Syncer{
		workerService: ws,
	}
}

// SyncHostServices is the main function that you will run periodically.
func (s *Syncer) SyncHostServices() {
	log.Println("Starting host service synchronization...")

	// 1. Find all services running on the host (multi-platform logic is here)
	discoveredServices, err := s.discoverHostServices()
	if err != nil {
		log.Printf("Error discovering host services: %v", err)
		return
	}

	log.Printf("Discovered %d services on the host.", len(discoveredServices))

	// 2. Synchronize each discovered service
	for _, discovered := range discoveredServices {
		existingWorker, err := s.workerService.GetWorkerServiceByName(discovered.Name)

		// A. If not in DB yet, create new one.
		if err != nil { // Assume error means "not found"
			log.Printf("Service '%s' not found in DB, creating new entry.", discovered.Name)
			errCreate := s.workerService.CreateWorkerService(
				discovered.Name,
				discovered.Description,
				model.StateEnabled,
			)
			if errCreate != nil {
				log.Printf("Failed to create worker service '%s': %v", discovered.Name, errCreate)
			}
			continue
		}

		// B. If already exists, update its status if needed.
		newStatus := mapStatusToModel(discovered.Status)
		if existingWorker.CurrentStatus != newStatus {
			log.Printf("Status for '%s' changed from '%s' to '%s'. Updating DB.",
				existingWorker.Name, existingWorker.CurrentStatus, newStatus)

			errUpdate := s.workerService.UpdateWorkerServiceStatus(existingWorker.ID, newStatus, model.HealthUnknown)
			if errUpdate != nil {
				log.Printf("Failed to update status for '%s': %v", existingWorker.Name, errUpdate)
			}
		}
	}
	log.Println("Host service synchronization finished.")
}

// discoverHostServices memeriksa OS saat runtime dan menjalankan logika yang sesuai.
func (s *Syncer) discoverHostServices() ([]DiscoveredService, error) {
	switch runtime.GOOS {
	case "linux":
		return s.discoverLinuxServices()
	case "darwin":
		return s.discoverMacServices()
	default:
		return nil, fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

// discoverLinuxServices berisi logika untuk systemd (Linux)
func (s *Syncer) discoverLinuxServices() ([]DiscoveredService, error) {
	commonServices := map[string]bool{"nginx.service": true, "postgresql.service": true}
	ctx := context.Background()
	conn, err := dbus.NewSystemConnectionContext(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	units, err := conn.ListUnitsByPatternsContext(ctx, []string{}, []string{"*.service"})
	if err != nil {
		return nil, err
	}
	var discoveredServices []DiscoveredService
	found := make(map[string]bool)
	for _, unit := range units {
		properties, _ := conn.GetUnitPropertiesContext(ctx, unit.Name)
		fragmentPath, _ := properties["FragmentPath"].(string)
		isCustomService := strings.HasPrefix(fragmentPath, "/etc/systemd/system/")
		isCommonService := commonServices[unit.Name]
		if (isCustomService || isCommonService) && !found[unit.Name] {
			discoveredServices = append(discoveredServices, DiscoveredService{
				Name: unit.Name, Description: unit.Description, Status: unit.ActiveState,
			})
			found[unit.Name] = true
		}
	}
	return discoveredServices, nil
}

// discoverMacServices berisi logika untuk brew services (macOS)
func (s *Syncer) discoverMacServices() ([]DiscoveredService, error) {
	// ... (implementasi 'brew services list' dari jawaban sebelumnya)
	cmdPath, err := exec.LookPath("brew")
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(cmdPath, "services", "list")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	var discoveredServices []DiscoveredService
	lines := strings.Split(out.String(), "\n")
	if len(lines) < 2 {
		return discoveredServices, nil
	}
	for _, line := range lines[1:] {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		discoveredServices = append(discoveredServices, DiscoveredService{
			Name: fields[0], Description: "Managed by Homebrew", Status: fields[1],
		})
	}
	return discoveredServices, nil
}

// mapStatusToModel menerjemahkan status mentah dari OS ke tipe model kita
func mapStatusToModel(rawStatus string) model.CurrentStatus {
	s := strings.ToLower(rawStatus)
	switch s {
	case "active", "started", "running":
		return model.StatusRunning
	case "inactive", "deactivating", "stopped", "none":
		return model.StatusStopped
	case "failed", "error":
		return model.StatusFailed
	case "activating":
		return model.StatusStarting
	default:
		return model.StatusStopped
	}
}
