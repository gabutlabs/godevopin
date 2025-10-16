package repository

import (
	"github.com/gabutlabs/devopin/internal/model"
	"gorm.io/gorm"
)

type WorkerServiceRepository interface {
	CreateWorkerService(WorkerService *model.WorkerService) error
	GetWorkerServiceByID(id uint) (*model.WorkerService, error)
	ListWorkerServices() ([]model.WorkerService, error)
	UpdateWorkerService(WorkerService *model.WorkerService) error
	DeleteWorkerService(id uint) error
	GetWorkerServiceByName(name string) (*model.WorkerService, error)
}

type workerServiceRepository struct {
	// Tambahkan field yang diperlukan, misalnya koneksi database
	db *gorm.DB
}

// NewWorkerServiceRepository membuat instance baru dari WorkerServiceRepository
func NewWorkerServiceRepository(db *gorm.DB) WorkerServiceRepository {
	return &workerServiceRepository{db: db}
}

// Implementasikan metode CRUD sesuai kebutuhan
// Contoh: CreateWorkerService, GetWorkerServiceByID, UpdateWorkerService, DeleteWorkerService, dll.
// Contoh metode untuk mendapatkan WorkerService berdasarkan ID
func (r *workerServiceRepository) GetWorkerServiceByID(id uint) (*model.WorkerService, error) {
	var WorkerService model.WorkerService
	if err := r.db.First(&WorkerService, id).Error; err != nil {
		return nil, err
	}
	return &WorkerService, nil
}

func (r *workerServiceRepository) CreateWorkerService(WorkerService *model.WorkerService) error {
	return r.db.Create(WorkerService).Error
}

func (r *workerServiceRepository) UpdateWorkerService(WorkerService *model.WorkerService) error {
	return r.db.Save(WorkerService).Error
}

func (r *workerServiceRepository) DeleteWorkerService(id uint) error {
	return r.db.Delete(&model.WorkerService{}, id).Error
}
func (r *workerServiceRepository) ListWorkerServices() ([]model.WorkerService, error) {
	var WorkerServices []model.WorkerService
	if err := r.db.Find(&WorkerServices).Error; err != nil {
		return nil, err
	}
	return WorkerServices, nil
}

// GetWorkerServiceByName retrieves a WorkerService by its name
func (r *workerServiceRepository) GetWorkerServiceByName(name string) (*model.WorkerService, error) {
	var WorkerService model.WorkerService
	if err := r.db.Where("name = ?", name).First(&WorkerService).Error; err != nil {
		return nil, err
	}
	return &WorkerService, nil
}
