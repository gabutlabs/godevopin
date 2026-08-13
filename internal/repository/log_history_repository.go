package repository

import (
	"github.com/gabutlabs/godevopin/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LogHistoryRepository interface {
	BatchInsert(logs []model.LogHistory) error
	GetLogsByProjectID(projectID uint, limit int, offset int) ([]model.LogHistory, int64, error)
}

type logHistoryRepository struct {
	db *gorm.DB
}

func NewLogHistoryRepository(db *gorm.DB) LogHistoryRepository {
	return &logHistoryRepository{db: db}
}

func (r *logHistoryRepository) BatchInsert(logs []model.LogHistory) error {
	if len(logs) == 0 {
		return nil
	}
	for i := range logs {
		if logs[i].ID == "" {
			logs[i].ID = uuid.NewString()
		}
	}
	return r.db.CreateInBatches(logs, 100).Error
}

func (r *logHistoryRepository) GetLogsByProjectID(projectID uint, limit int, offset int) ([]model.LogHistory, int64, error) {
	var logs []model.LogHistory
	var total int64

	db := r.db.Model(&model.LogHistory{}).Where("project_id = ?", projectID)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Limit(limit).Offset(offset).Order("timestamp desc").Find(&logs).Error
	return logs, total, err
}
