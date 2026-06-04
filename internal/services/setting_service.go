package service

import (
	"github.com/gabutlabs/godevopin/internal/model"
	"github.com/gabutlabs/godevopin/internal/repository"
)

type UpdateSettingRequest struct {
	MonitoringIntervalSeconds     int     `json:"monitoring_interval_seconds"`
	AlarmCheckIntervalSeconds     int     `json:"alarm_check_interval_seconds"`
	AlarmRepeatIntervalMinutes    int     `json:"alarm_repeat_interval_minutes"`
	SystemCPUCriticalPercent      float64 `json:"system_cpu_critical_percent"`
	SystemDiskCriticalPercent     float64 `json:"system_disk_critical_percent"`
	SystemMemCriticalPercent      float64 `json:"system_mem_critical_percent"`
	WorkerHeartbeatTimeoutSeconds int     `json:"worker_heartbeat_timeout_seconds"`
	AiProvider                    string  `json:"ai_provider"`
	AiModelName                   string  `json:"ai_model_name"`
	AiApiKey                      string  `json:"ai_api_key"`
	AiBaseURL                     string  `json:"ai_base_url"`
	TelegramBotToken              string  `json:"telegram_bot_token"`
}

type SettingService interface {
	GetSettings() (*model.AppSetting, error)
	UpdateSettings(req *UpdateSettingRequest) error
}

type settingService struct {
	repo repository.SettingRepository
}

func NewSettingService(repo repository.SettingRepository) SettingService {
	return &settingService{repo: repo}
}

func (s *settingService) GetSettings() (*model.AppSetting, error) {
	return s.repo.GetSettings()
}

func (s *settingService) UpdateSettings(req *UpdateSettingRequest) error {
	setting, err := s.repo.GetSettings()
	if err != nil {
		return err
	}

	setting.MonitoringIntervalSeconds = req.MonitoringIntervalSeconds
	setting.AlarmCheckIntervalSeconds = req.AlarmCheckIntervalSeconds
	setting.AlarmRepeatIntervalMinutes = req.AlarmRepeatIntervalMinutes
	setting.SystemCPUCriticalPercent = req.SystemCPUCriticalPercent
	setting.SystemDiskCriticalPercent = req.SystemDiskCriticalPercent
	setting.SystemMemCriticalPercent = req.SystemMemCriticalPercent
	setting.WorkerHeartbeatTimeoutSeconds = req.WorkerHeartbeatTimeoutSeconds
	setting.AIProvider = req.AiProvider
	setting.AIModelName = req.AiModelName
	setting.AIApiKey = req.AiApiKey
	setting.AIBaseURL = req.AiBaseURL
	setting.TelegramBotToken = req.TelegramBotToken

	return s.repo.UpdateSettings(setting)
}
