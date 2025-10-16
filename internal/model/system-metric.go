package model

import "time"

type SystemMetric struct {
	ID            string    `gorm:"primaryKey" json:"id"`
	CPUUsage      float64   `json:"cpu_usage"`
	MemUsageByte  float64   `json:"mem_usage_byte"`
	MemTotalByte  float64   `json:"mem_total_byte"`
	DiskUsageByte float64   `json:"disk_usage_byte"`
	DiskTotalByte float64   `json:"disk_total_byte"`
	CreatedAt     time.Time `gorm:"primaryKey;type:timestamp with time zone;not null" json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
