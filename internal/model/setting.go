package model

import "time"

type AppSetting struct {
	ID                            uint      `gorm:"primaryKey" json:"id"`
	MonitoringIntervalSeconds     int       `gorm:"default:10" json:"monitoring_interval_seconds"`
	AlarmCheckIntervalSeconds     int       `gorm:"default:60" json:"alarm_check_interval_seconds"`
	AlarmRepeatIntervalMinutes    int       `gorm:"default:15" json:"alarm_repeat_interval_minutes"`
	SystemCPUCriticalPercent      float64   `gorm:"default:90.0" json:"system_cpu_critical_percent"`
	SystemDiskCriticalPercent     float64   `gorm:"default:85.0" json:"system_disk_critical_percent"`
	SystemMemCriticalPercent      float64   `gorm:"default:85.0" json:"system_mem_critical_percent"`
	WorkerHeartbeatTimeoutSeconds int       `gorm:"default:300" json:"worker_heartbeat_timeout_seconds"`
	UpdatedAt                     time.Time `json:"updated_at"`
}
