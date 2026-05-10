package service

import (
	"github.com/gabutlabs/devopin/internal/model"
	"github.com/gabutlabs/devopin/internal/repository"
)

type UpdateSettingRequest struct {
	MonitoringIntervalSeconds     int     `json:"monitoring_interval_seconds"`
	AlarmCheckIntervalSeconds     int     `json:"alarm_check_interval_seconds"`
	AlarmRepeatIntervalMinutes    int     `json:"alarm_repeat_interval_minutes"`
	SystemCPUCriticalPercent      float64 `json:"system_cpu_critical_percent"`
	SystemDiskCriticalPercent     float64 `json:"system_disk_critical_percent"`
	SystemMemCriticalPercent      float64 `json:"system_mem_critical_percent"`
	WorkerHeartbeatTimeoutSeconds int     `json:"worker_heartbeat_timeout_seconds"`
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

	return s.repo.UpdateSettings(setting)
}
