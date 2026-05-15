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
	TelegramBotToken              string    `gorm:"type:varchar(255)" json:"telegram_bot_token"`
	AIProvider                    string    `gorm:"type:varchar(50);default:'gemini'" json:"ai_provider"`
	AIModelName                   string    `gorm:"type:varchar(100)" json:"ai_model_name"`
	AIApiKey                      string    `gorm:"type:varchar(255)" json:"ai_api_key"`
	AIBaseURL                     string    `gorm:"type:varchar(255)" json:"ai_base_url"`
	UpdatedAt                     time.Time `json:"updated_at"`
}
