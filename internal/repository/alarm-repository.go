package repository

import (
	"github.com/gabutlabs/godevopin/internal/model"
	"github.com/gabutlabs/godevopin/pkg"
	"gorm.io/gorm"
)

// PaginationResult is a temporary struct to handle the return type from GORM queries
type PaginationResult[T any] struct {
	Items []T
	Total int64
}

type AlarmRepository interface {
	CreateActiveAlarm(alarm *model.ActiveAlarm) error
	GetActiveAlarmByNameAndTarget(alarmName, target string) (*model.ActiveAlarm, error)
	GetActiveAlarmByAlarmName(alarmName string) ([]model.ActiveAlarm, error)
	GetAllActiveAlarms() ([]model.ActiveAlarm, error)
	GetAllActiveAlarmsPaginated(page, limit int, sort string) (pkg.PaginatedResponse[model.ActiveAlarm], error)
	UpdateActiveAlarm(alarm *model.ActiveAlarm) error
	DeleteActiveAlarm(alarmName, target string) error
	GetCountActiveAlarm(status string) (int64, error)
	CreateAlarmHistory(history *model.AlarmHistory) error
	GetAlarmHistoryByID(id uint) (*model.AlarmHistory, error)
	GetAlarmHistoryByAlarmName(alarmName string) ([]model.AlarmHistory, error)
	GetAlarmHistoryByTarget(target string) ([]model.AlarmHistory, error)
	GetAlarmHistoryByStatus(status model.AlarmHistoryStatus) ([]model.AlarmHistory, error)
	GetAllAlarmHistory(limit, offset int) ([]model.AlarmHistory, error)
	GetAllAlarmHistoryPaginated(page, limit int, sort string, filterStatus string) (pkg.PaginatedResponse[model.AlarmHistory], error)
	UpdateAlarmHistory(history *model.AlarmHistory) error
	DeleteAlarmHistory(id uint) error
	GetAlarmHistoryByAlarmNameAndTarget(alarmName, target string, limit, offset int) ([]model.AlarmHistory, error)
	GetAlarmHistoryByAlarmNameAndTargetPaginated(alarmName, target string, page, limit int, sort string) (pkg.PaginatedResponse[model.AlarmHistory], error)
	GetActiveAlarmsByStatus(status model.ActiveAlarmStatus) ([]model.ActiveAlarm, error)
}

type alarmRepository struct {
	db *gorm.DB
}

// NewAlarmRepository creates a new instance of AlarmRepository
func NewAlarmRepository(db *gorm.DB) AlarmRepository {
	return &alarmRepository{db: db}
}

// ActiveAlarm methods
func (r *alarmRepository) CreateActiveAlarm(alarm *model.ActiveAlarm) error {
	return r.db.Create(alarm).Error
}

func (r *alarmRepository) GetActiveAlarmByNameAndTarget(alarmName, target string) (*model.ActiveAlarm, error) {
	var alarm model.ActiveAlarm
	if err := r.db.Where("alarm_name = ? AND target = ?", alarmName, target).First(&alarm).Error; err != nil {
		return nil, err
	}
	return &alarm, nil
}

func (r *alarmRepository) GetActiveAlarmByAlarmName(alarmName string) ([]model.ActiveAlarm, error) {
	var alarms []model.ActiveAlarm
	if err := r.db.Where("alarm_name = ?", alarmName).Find(&alarms).Error; err != nil {
		return nil, err
	}
	return alarms, nil
}

func (r *alarmRepository) GetAllActiveAlarms() ([]model.ActiveAlarm, error) {
	var alarms []model.ActiveAlarm
	if err := r.db.Find(&alarms).Error; err != nil {
		return nil, err
	}
	return alarms, nil
}

func (r *alarmRepository) GetCountActiveAlarm(status string) (int64, error) {
	var total int64
	query := r.db.Model(&model.ActiveAlarm{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if err := query.Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func (r *alarmRepository) UpdateActiveAlarm(alarm *model.ActiveAlarm) error {
	return r.db.Save(alarm).Error
}

func (r *alarmRepository) DeleteActiveAlarm(alarmName, target string) error {
	return r.db.Where("alarm_name = ? AND target = ?", alarmName, target).Delete(&model.ActiveAlarm{}).Error
}

func (r *alarmRepository) GetActiveAlarmsByStatus(status model.ActiveAlarmStatus) ([]model.ActiveAlarm, error) {
	var alarms []model.ActiveAlarm
	if err := r.db.Where("status = ?", string(status)).Find(&alarms).Error; err != nil {
		return nil, err
	}
	return alarms, nil
}

// AlarmHistory methods
func (r *alarmRepository) CreateAlarmHistory(history *model.AlarmHistory) error {
	return r.db.Create(history).Error
}

func (r *alarmRepository) GetAlarmHistoryByID(id uint) (*model.AlarmHistory, error) {
	var history model.AlarmHistory
	if err := r.db.First(&history, id).Error; err != nil {
		return nil, err
	}
	return &history, nil
}

func (r *alarmRepository) GetAlarmHistoryByAlarmName(alarmName string) ([]model.AlarmHistory, error) {
	var history []model.AlarmHistory
	if err := r.db.Where("alarm_name = ?", alarmName).Find(&history).Error; err != nil {
		return nil, err
	}
	return history, nil
}

func (r *alarmRepository) GetAlarmHistoryByTarget(target string) ([]model.AlarmHistory, error) {
	var history []model.AlarmHistory
	if err := r.db.Where("target = ?", target).Find(&history).Error; err != nil {
		return nil, err
	}
	return history, nil
}

func (r *alarmRepository) GetAlarmHistoryByStatus(status model.AlarmHistoryStatus) ([]model.AlarmHistory, error) {
	var history []model.AlarmHistory
	if err := r.db.Where("status = ?", string(status)).Find(&history).Error; err != nil {
		return nil, err
	}
	return history, nil
}

func (r *alarmRepository) GetAllAlarmHistory(limit, offset int) ([]model.AlarmHistory, error) {
	var history []model.AlarmHistory
	if err := r.db.Limit(limit).Offset(offset).Find(&history).Error; err != nil {
		return nil, err
	}
	return history, nil
}

func (r *alarmRepository) UpdateAlarmHistory(history *model.AlarmHistory) error {
	return r.db.Save(history).Error
}

func (r *alarmRepository) DeleteAlarmHistory(id uint) error {
	return r.db.Delete(&model.AlarmHistory{}, id).Error
}

func (r *alarmRepository) GetAlarmHistoryByAlarmNameAndTarget(alarmName, target string, limit, offset int) ([]model.AlarmHistory, error) {
	var history []model.AlarmHistory
	if err := r.db.Where("alarm_name = ? AND target = ?", alarmName, target).Limit(limit).Offset(offset).Find(&history).Error; err != nil {
		return nil, err
	}
	return history, nil
}

// Pagination methods
func (r *alarmRepository) GetAllActiveAlarmsPaginated(page, limit int, sort string) (pkg.PaginatedResponse[model.ActiveAlarm], error) {
	var alarms []model.ActiveAlarm
	var total int64

	offset := (page - 1) * limit

	query := r.db.Model(&model.ActiveAlarm{})

	if err := query.Count(&total).Error; err != nil {
		return pkg.PaginatedResponse[model.ActiveAlarm]{}, err
	}

	if sort != "" {
		query = query.Order(sort)
	} else {
		query = query.Order("created_at DESC") // Default sort
	}

	if err := query.Offset(offset).Limit(limit).Find(&alarms).Error; err != nil {
		return pkg.PaginatedResponse[model.ActiveAlarm]{}, err
	}

	pagination := pkg.Pagination{
		Page:  page,
		Limit: limit,
		Sort:  sort,
		Total: total,
	}

	return pkg.PaginatedResponse[model.ActiveAlarm]{
		Data:       alarms,
		Pagination: pagination,
	}, nil
}

func (r *alarmRepository) GetAllAlarmHistoryPaginated(page, limit int, sort string, filterStatus string) (pkg.PaginatedResponse[model.AlarmHistory], error) {
	var history []model.AlarmHistory
	var total int64

	offset := (page - 1) * limit

	query := r.db.Model(&model.AlarmHistory{})

	if err := query.Count(&total).Error; err != nil {
		return pkg.PaginatedResponse[model.AlarmHistory]{}, err
	}

	if sort != "" {
		query = query.Order(sort)
	} else {
		query = query.Order("created_at DESC") // Default sort
	}

	if filterStatus != "" {
		query = query.Where("status = ?", filterStatus)
	}

	if err := query.Offset(offset).Limit(limit).Find(&history).Error; err != nil {
		return pkg.PaginatedResponse[model.AlarmHistory]{}, err
	}

	pagination := pkg.Pagination{
		Page:  page,
		Limit: limit,
		Sort:  sort,
		Total: total,
	}

	return pkg.PaginatedResponse[model.AlarmHistory]{
		Data:       history,
		Pagination: pagination,
	}, nil
}

func (r *alarmRepository) GetAlarmHistoryByAlarmNameAndTargetPaginated(alarmName, target string, page, limit int, sort string) (pkg.PaginatedResponse[model.AlarmHistory], error) {
	var history []model.AlarmHistory
	var total int64

	offset := (page - 1) * limit

	query := r.db.Model(&model.AlarmHistory{}).Where("alarm_name = ? AND target = ?", alarmName, target)

	if err := query.Count(&total).Error; err != nil {
		return pkg.PaginatedResponse[model.AlarmHistory]{}, err
	}

	if sort != "" {
		query = query.Order(sort)
	} else {
		query = query.Order("created_at DESC") // Default sort
	}

	if err := query.Offset(offset).Limit(limit).Find(&history).Error; err != nil {
		return pkg.PaginatedResponse[model.AlarmHistory]{}, err
	}

	pagination := pkg.Pagination{
		Page:  page,
		Limit: limit,
		Sort:  sort,
		Total: total,
	}

	return pkg.PaginatedResponse[model.AlarmHistory]{
		Data:       history,
		Pagination: pagination,
	}, nil
}
