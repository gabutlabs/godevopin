package model

import "time"

// ProcessMetric is both the live process representation and the historical
// process sample stored in the metrics database. ID is only used for the
// persisted record and is intentionally omitted from API responses.
type ProcessMetric struct {
	ID                 string    `gorm:"primaryKey" json:"-"`
	ProcessKey         string    `gorm:"not null;index" json:"process_key"`
	PID                int32     `gorm:"not null;index" json:"pid"`
	ProcessStartTime   int64     `gorm:"not null" json:"process_start_time"`
	ParentPID          int32     `json:"parent_pid"`
	Name               string    `gorm:"not null" json:"name"`
	Executable         string    `json:"executable,omitempty"`
	Username           string    `json:"username,omitempty"`
	Status             string    `json:"status"`
	CPUPercent         float64   `json:"cpu_percent"`
	MemoryPercent      float64   `json:"memory_percent"`
	MemoryBytes        uint64    `json:"memory_bytes"`
	VirtualMemoryBytes uint64    `json:"virtual_memory_bytes"`
	ThreadCount        int32     `json:"thread_count"`
	ObservedAt         time.Time `gorm:"not null;index" json:"observed_at"`
}
