package repository

import (
	"time"

	"github.com/gabutlabs/godevopin/internal/model"
	"gorm.io/gorm"
)

type PostgreSQLTargetRepository interface {
	ListPostgreSQLTargets() ([]model.PostgreSQLTarget, error)
	GetPostgreSQLTargetByID(id uint) (*model.PostgreSQLTarget, error)
	CreatePostgreSQLTarget(target *model.PostgreSQLTarget) error
	UpdatePostgreSQLTarget(target *model.PostgreSQLTarget) error
	DeletePostgreSQLTarget(id uint) error
	UpdatePostgreSQLTargetHealth(id uint, checkedAt *time.Time, lastError string) error
}

type postgresqlTargetRepository struct {
	db *gorm.DB
}

func NewPostgreSQLTargetRepository(db *gorm.DB) PostgreSQLTargetRepository {
	return &postgresqlTargetRepository{db: db}
}

func (r *postgresqlTargetRepository) ListPostgreSQLTargets() ([]model.PostgreSQLTarget, error) {
	var targets []model.PostgreSQLTarget
	err := r.db.Order("name ASC").Find(&targets).Error
	return targets, err
}

func (r *postgresqlTargetRepository) GetPostgreSQLTargetByID(id uint) (*model.PostgreSQLTarget, error) {
	var target model.PostgreSQLTarget
	if err := r.db.First(&target, id).Error; err != nil {
		return nil, err
	}
	return &target, nil
}

func (r *postgresqlTargetRepository) CreatePostgreSQLTarget(target *model.PostgreSQLTarget) error {
	return r.db.Create(target).Error
}

func (r *postgresqlTargetRepository) UpdatePostgreSQLTarget(target *model.PostgreSQLTarget) error {
	return r.db.Save(target).Error
}

func (r *postgresqlTargetRepository) DeletePostgreSQLTarget(id uint) error {
	return r.db.Delete(&model.PostgreSQLTarget{}, id).Error
}

func (r *postgresqlTargetRepository) UpdatePostgreSQLTargetHealth(id uint, checkedAt *time.Time, lastError string) error {
	return r.db.Model(&model.PostgreSQLTarget{}).Where("id = ?", id).Updates(map[string]any{
		"last_checked_at": checkedAt,
		"last_error":      lastError,
	}).Error
}
