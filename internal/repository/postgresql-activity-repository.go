package repository

import (
	"errors"
	"strings"
	"time"

	"github.com/gabutlabs/godevopin/internal/model"
	"gorm.io/gorm"
)

type PostgreSQLActivityHistoryFilter struct {
	TargetID        uint
	Cutoff          time.Time
	ActivityKey     string
	Search          string
	State           string
	BlockedOnly     bool
	MinimumDuration int64
	Limit           int
}

type PostgreSQLActivityRepository interface {
	CreatePostgreSQLActivities(activities []model.PostgreSQLActivity) error
	DeleteOlderPostgreSQLActivities(cutoff time.Time) (int64, error)
	ListPostgreSQLActivityHistory(filter PostgreSQLActivityHistoryFilter) ([]model.PostgreSQLActivity, error)
}

type postgresqlActivityRepository struct {
	db *gorm.DB
}

func NewPostgreSQLActivityRepository(db *gorm.DB) PostgreSQLActivityRepository {
	return &postgresqlActivityRepository{db: db}
}

func (r *postgresqlActivityRepository) CreatePostgreSQLActivities(activities []model.PostgreSQLActivity) error {
	if len(activities) == 0 {
		return nil
	}
	return r.db.Create(&activities).Error
}

func (r *postgresqlActivityRepository) DeleteOlderPostgreSQLActivities(cutoff time.Time) (int64, error) {
	result := r.db.Where("observed_at < ?", cutoff.UTC()).Delete(&model.PostgreSQLActivity{})
	return result.RowsAffected, result.Error
}

func (r *postgresqlActivityRepository) ListPostgreSQLActivityHistory(filter PostgreSQLActivityHistoryFilter) ([]model.PostgreSQLActivity, error) {
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
	if filter.BlockedOnly {
		query = query.Where("blocking_p_id > 0")
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

	var activities []model.PostgreSQLActivity
	if err := query.Order("observed_at DESC").Limit(filter.Limit).Find(&activities).Error; err != nil {
		return nil, err
	}
	return activities, nil
}
