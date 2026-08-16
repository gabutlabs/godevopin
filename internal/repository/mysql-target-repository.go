package repository

import (
	"time"

	"github.com/gabutlabs/godevopin/internal/model"
	"gorm.io/gorm"
)

type MySQLTargetRepository interface {
	ListMySQLTargets() ([]model.MySQLTarget, error)
	GetMySQLTargetByID(id uint) (*model.MySQLTarget, error)
	CreateMySQLTarget(target *model.MySQLTarget) error
	UpdateMySQLTarget(target *model.MySQLTarget) error
	DeleteMySQLTarget(id uint) error
	UpdateMySQLTargetHealth(id uint, checkedAt *time.Time, lastError string) error
}

type mysqlTargetRepository struct{ db *gorm.DB }

func NewMySQLTargetRepository(db *gorm.DB) MySQLTargetRepository {
	return &mysqlTargetRepository{db: db}
}

func (r *mysqlTargetRepository) ListMySQLTargets() ([]model.MySQLTarget, error) {
	var targets []model.MySQLTarget
	err := r.db.Order("name ASC").Find(&targets).Error
	return targets, err
}

func (r *mysqlTargetRepository) GetMySQLTargetByID(id uint) (*model.MySQLTarget, error) {
	var target model.MySQLTarget
	if err := r.db.First(&target, id).Error; err != nil {
		return nil, err
	}
	return &target, nil
}

func (r *mysqlTargetRepository) CreateMySQLTarget(target *model.MySQLTarget) error {
	return r.db.Create(target).Error
}

func (r *mysqlTargetRepository) UpdateMySQLTarget(target *model.MySQLTarget) error {
	return r.db.Save(target).Error
}

func (r *mysqlTargetRepository) DeleteMySQLTarget(id uint) error {
	return r.db.Delete(&model.MySQLTarget{}, id).Error
}

func (r *mysqlTargetRepository) UpdateMySQLTargetHealth(id uint, checkedAt *time.Time, lastError string) error {
	return r.db.Model(&model.MySQLTarget{}).Where("id = ?", id).Updates(map[string]any{
		"last_checked_at": checkedAt,
		"last_error":      lastError,
	}).Error
}
