package model

import (
	"time"
)

// Definisikan tipe kustom untuk status agar lebih type-safe
type DesiredState string
type CurrentStatus string
type HealthStatus string

// Definisikan konstanta untuk setiap status
const (
	StateEnabled  DesiredState = "enabled"
	StateDisabled DesiredState = "disabled"

	StatusStarting CurrentStatus = "starting"
	StatusRunning  CurrentStatus = "running"
	StatusStopped  CurrentStatus = "stopped"
	StatusRestart  CurrentStatus = "restart"
	StatusFailed   CurrentStatus = "failed"
	StatusDegraded CurrentStatus = "degraded"

	HealthHealthy   HealthStatus = "healthy"
	HealthUnhealthy HealthStatus = "unhealthy"
	HealthUnknown   HealthStatus = "unknown"
)

// WorkerService adalah model GORM yang merepresentasikan tabel worker_services
type WorkerService struct {
	// Kolom Identitas
	ID          uint   `gorm:"primarykey" json:"id"`
	Name        string `gorm:"type:varchar(100);not null;unique" json:"name"`
	Description string `gorm:"type:text" json:"description"`

	// Kolom Status & State (menggunakan tipe kustom)
	DesiredState  DesiredState  `gorm:"type:varchar(20);not null;default:'enabled'" json:"desired_state"`
	CurrentStatus CurrentStatus `gorm:"type:varchar(20);not null;default:'stopped'" json:"current_status"`
	HealthStatus  HealthStatus  `gorm:"type:varchar(20);not null;default:'unknown'" json:"health_status"`

	// Kolom Data Operasional (menggunakan pointer untuk nullable)
	LastHeartbeatAt  *time.Time `json:"last_heartbeat_at,omitempty"`
	LastSuccessAt    *time.Time `json:"last_success_at,omitempty"`
	LastFailureAt    *time.Time `json:"last_failure_at,omitempty"`
	LastErrorMessage *string    `gorm:"type:text" json:"last_error_message,omitempty"` // Typo diperbaiki

	// Kolom Metadata (menggunakan pointer untuk nullable)
	PID     *int    `json:"pid,omitempty"`
	Version *string `gorm:"type:varchar(50)" json:"version,omitempty"`

	// Timestamps (GORM akan mengelolanya secara otomatis)
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
