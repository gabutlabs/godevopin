package repository

import (
	"errors"
	"strings"
	"time"

	"github.com/gabutlabs/godevopin/internal/model"
	"gorm.io/gorm"
)

type MySQLActivityHistoryFilter struct {
	TargetID        uint
	Cutoff          time.Time
	ActivityKey     string
	Search          string
	State           string
	MinimumDuration int64
	Limit           int
}

type MySQLActivityRepository interface {
	CreateMySQLActivities(activities []model.MySQLActivity) error
	DeleteOlderMySQLActivities(cutoff time.Time) (int64, error)
	ListMySQLActivityHistory(filter MySQLActivityHistoryFilter) ([]model.MySQLActivity, error)
}

type mysqlActivityRepository struct{ db *gorm.DB }

func NewMySQLActivityRepository(db *gorm.DB) MySQLActivityRepository {
	return &mysqlActivityRepository{db: db}
}

func (r *mysqlActivityRepository) CreateMySQLActivities(activities []model.MySQLActivity) error {
	if len(activities) == 0 {
		return nil
	}
	return r.db.Create(&activities).Error
}

func (r *mysqlActivityRepository) DeleteOlderMySQLActivities(cutoff time.Time) (int64, error) {
	result := r.db.Where("observed_at < ?", cutoff.UTC()).Delete(&model.MySQLActivity{})
	return result.RowsAffected, result.Error
}

func (r *mysqlActivityRepository) ListMySQLActivityHistory(filter MySQLActivityHistoryFilter) ([]model.MySQLActivity, error) {
	if filter.Cutoff.IsZero() {
		return nil, errors.New("activity history cutoff is required")
	}
	query := r.db.Where("observed_at >= ?", filter.Cutoff.UTC())
	if filter.TargetID > 0 {
		query = query.Where("target_id = ?", filter.TargetID)
	}
	if strings.TrimSpace(filter.ActivityKey) != "" {
		query = query.Where("activity_key = ?", strings.TrimSpace(filter.ActivityKey))
	}
	if strings.TrimSpace(filter.State) != "" {
		query = query.Where("state = ?", strings.TrimSpace(filter.State))
	}
	if filter.MinimumDuration > 0 {
		query = query.Where("query_duration_ms >= ?", filter.MinimumDuration)
	}
	if search := strings.TrimSpace(filter.Search); search != "" {
		pattern := "%" + search + "%"
		query = query.Where(
			"target_name LIKE ? OR database_name LIKE ? OR username LIKE ? OR query LIKE ? OR activity_key LIKE ?",
			pattern, pattern, pattern, pattern, pattern,
		)
	}
	if filter.Limit <= 0 {
		filter.Limit = 500
	}
	if filter.Limit > 2000 {
		filter.Limit = 2000
	}

	var activities []model.MySQLActivity
	if err := query.Order("observed_at DESC").Limit(filter.Limit).Find(&activities).Error; err != nil {
		return nil, err
	}
	return activities, nil
}
