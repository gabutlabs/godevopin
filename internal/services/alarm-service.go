package service

import (
	"time"

	"github.com/gabutlabs/devopin/internal/model"
	"github.com/gabutlabs/devopin/internal/repository"
	"github.com/gabutlabs/devopin/pkg"
)

type AlarmService interface {
	CreateActiveAlarm(alarmName, target string, status model.ActiveAlarmStatus, message string) (*model.ActiveAlarm, error)
	GetActiveAlarmByNameAndTarget(alarmName, target string) (*model.ActiveAlarm, error)
	GetActiveAlarmByAlarmName(alarmName string) ([]model.ActiveAlarm, error)
	GetAllActiveAlarms() ([]model.ActiveAlarm, error)
	GetAllActiveAlarmsPaginated(page, limit int, sort string) (pkg.PaginatedResponse[model.ActiveAlarm], error)
	UpdateActiveAlarm(alarmName, target string, status model.ActiveAlarmStatus, acknowledgedBy *string, message string) error
	DeleteActiveAlarm(alarmName, target string) error
	AcknowledgeActiveAlarm(alarmName, target, acknowledgedBy string) error
	GetActiveAlarmsByStatus(status model.ActiveAlarmStatus) ([]model.ActiveAlarm, error)

	CreateAlarmHistory(alarmName, target string, status model.AlarmHistoryStatus, message string, metadata []byte) error
	GetAlarmHistoryByID(id uint) (*model.AlarmHistory, error)
	GetAlarmHistoryByAlarmName(alarmName string) ([]model.AlarmHistory, error)
	GetAlarmHistoryByTarget(target string) ([]model.AlarmHistory, error)
	GetAlarmHistoryByStatus(status model.AlarmHistoryStatus) ([]model.AlarmHistory, error)
	GetAllAlarmHistory(limit, offset int) ([]model.AlarmHistory, error)
	GetAllAlarmHistoryPaginated(page, limit int, sort string) (pkg.PaginatedResponse[model.AlarmHistory], error)
	DeleteAlarmHistory(id uint) error
	GetAlarmHistoryByAlarmNameAndTarget(alarmName, target string, limit, offset int) ([]model.AlarmHistory, error)
	GetAlarmHistoryByAlarmNameAndTargetPaginated(alarmName, target string, page, limit int, sort string) (pkg.PaginatedResponse[model.AlarmHistory], error)
}

type alarmService struct {
	repository repository.AlarmRepository
}

func NewAlarmService(repo repository.AlarmRepository) AlarmService {
	return &alarmService{repository: repo}
}

// ActiveAlarm methods
func (s *alarmService) CreateActiveAlarm(alarmName, target string, status model.ActiveAlarmStatus, message string) (*model.ActiveAlarm, error) {
	existingAlarm, err := s.repository.GetActiveAlarmByNameAndTarget(alarmName, target)

	if err == nil && existingAlarm != nil {
		// --- BLOK UPDATE ---
		// 1. Siapkan data alarm yang akan di-update
		alarmToUpdate := &model.ActiveAlarm{
			AlarmName:      alarmName,
			Target:         target,
			Status:         status,
			Message:        message,
			StartedAt:      existingAlarm.StartedAt, // Pertahankan waktu mulai
			LastNotifiedAt: time.Now(),              // Update waktu notifikasi
		}

		// 2. Lakukan update
		err = s.repository.UpdateActiveAlarm(alarmToUpdate)
		if err != nil {
			return nil, err // Gagal update, kembalikan error
		}

		// 3. Kembalikan alarm yang sudah di-update
		return alarmToUpdate, nil
	}

	// --- BLOK CREATE ---
	// 1. Buat record alarm baru
	newAlarm := &model.ActiveAlarm{
		AlarmName:      alarmName,
		Target:         target,
		Status:         status,
		Message:        message,
		StartedAt:      time.Now(),
		LastNotifiedAt: time.Now(),
	}

	// 2. Simpan ke repository
	err = s.repository.CreateActiveAlarm(newAlarm)
	if err != nil {
		return nil, err // Gagal create, kembalikan error
	}

	// 3. Kembalikan alarm yang baru dibuat
	return newAlarm, nil
}

func (s *alarmService) GetActiveAlarmByNameAndTarget(alarmName, target string) (*model.ActiveAlarm, error) {
	return s.repository.GetActiveAlarmByNameAndTarget(alarmName, target)
}

func (s *alarmService) GetActiveAlarmByAlarmName(alarmName string) ([]model.ActiveAlarm, error) {
	return s.repository.GetActiveAlarmByAlarmName(alarmName)
}

func (s *alarmService) GetAllActiveAlarms() ([]model.ActiveAlarm, error) {
	return s.repository.GetAllActiveAlarms()
}

func (s *alarmService) UpdateActiveAlarm(alarmName, target string, status model.ActiveAlarmStatus, acknowledgedBy *string, message string) error {
	currentAlarm, err := s.repository.GetActiveAlarmByNameAndTarget(alarmName, target)
	if err != nil {
		return err
	}

	currentAlarm.Status = status
	currentAlarm.AcknowledgedBy = acknowledgedBy
	currentAlarm.Message = message
	currentAlarm.LastNotifiedAt = time.Now()

	return s.repository.UpdateActiveAlarm(currentAlarm)
}

func (s *alarmService) DeleteActiveAlarm(alarmName, target string) error {
	return s.repository.DeleteActiveAlarm(alarmName, target)
}

func (s *alarmService) AcknowledgeActiveAlarm(alarmName, target, acknowledgedBy string) error {
	currentAlarm, err := s.repository.GetActiveAlarmByNameAndTarget(alarmName, target)
	if err != nil {
		return err
	}

	currentAlarm.Status = model.StatusAcknowledged
	currentAlarm.AcknowledgedBy = &acknowledgedBy
	currentAlarm.LastNotifiedAt = time.Now()

	return s.repository.UpdateActiveAlarm(currentAlarm)
}

func (s *alarmService) GetActiveAlarmsByStatus(status model.ActiveAlarmStatus) ([]model.ActiveAlarm, error) {
	return s.repository.GetActiveAlarmsByStatus(status)
}

// AlarmHistory methods
func (s *alarmService) CreateAlarmHistory(alarmName, target string, status model.AlarmHistoryStatus, message string, metadata []byte) error {
	history := &model.AlarmHistory{
		AlarmName: alarmName,
		Target:    target,
		Status:    status,
		Message:   message,
		Metadata:  metadata,
		CreatedAt: time.Now(),
	}

	return s.repository.CreateAlarmHistory(history)
}

func (s *alarmService) GetAlarmHistoryByID(id uint) (*model.AlarmHistory, error) {
	return s.repository.GetAlarmHistoryByID(id)
}

func (s *alarmService) GetAlarmHistoryByAlarmName(alarmName string) ([]model.AlarmHistory, error) {
	return s.repository.GetAlarmHistoryByAlarmName(alarmName)
}

func (s *alarmService) GetAlarmHistoryByTarget(target string) ([]model.AlarmHistory, error) {
	return s.repository.GetAlarmHistoryByTarget(target)
}

func (s *alarmService) GetAlarmHistoryByStatus(status model.AlarmHistoryStatus) ([]model.AlarmHistory, error) {
	return s.repository.GetAlarmHistoryByStatus(status)
}

func (s *alarmService) GetAllAlarmHistory(limit, offset int) ([]model.AlarmHistory, error) {
	return s.repository.GetAllAlarmHistory(limit, offset)
}

func (s *alarmService) DeleteAlarmHistory(id uint) error {
	return s.repository.DeleteAlarmHistory(id)
}

func (s *alarmService) GetAlarmHistoryByAlarmNameAndTarget(alarmName, target string, limit, offset int) ([]model.AlarmHistory, error) {
	return s.repository.GetAlarmHistoryByAlarmNameAndTarget(alarmName, target, limit, offset)
}

// Pagination methods
func (s *alarmService) GetAllActiveAlarmsPaginated(page, limit int, sort string) (pkg.PaginatedResponse[model.ActiveAlarm], error) {
	return s.repository.GetAllActiveAlarmsPaginated(page, limit, sort)
}

func (s *alarmService) GetAllAlarmHistoryPaginated(page, limit int, sort string) (pkg.PaginatedResponse[model.AlarmHistory], error) {
	return s.repository.GetAllAlarmHistoryPaginated(page, limit, sort)
}

func (s *alarmService) GetAlarmHistoryByAlarmNameAndTargetPaginated(alarmName, target string, page, limit int, sort string) (pkg.PaginatedResponse[model.AlarmHistory], error) {
	return s.repository.GetAlarmHistoryByAlarmNameAndTargetPaginated(alarmName, target, page, limit, sort)
}
