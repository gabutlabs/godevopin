package repository

import (
	"github.com/gabutlabs/devopin/internal/model"
	"gorm.io/gorm"
)

type ProjectRepository interface {
	CreateProject(project *model.Project) error
	GetProjectByID(id uint) (*model.Project, error)
	ListProjects() ([]model.Project, error)
	UpdateProject(project *model.Project) error
	DeleteProject(id uint) error
}

type projectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) CreateProject(project *model.Project) error {
	return r.db.Create(project).Error
}

func (r *projectRepository) GetProjectByID(id uint) (*model.Project, error) {
	var project model.Project
	if err := r.db.First(&project, id).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *projectRepository) ListProjects() ([]model.Project, error) {
	var projects []model.Project
	if err := r.db.Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *projectRepository) UpdateProject(project *model.Project) error {
	return r.db.Save(project).Error
}

func (r *projectRepository) DeleteProject(id uint) error {
	return r.db.Delete(&model.Project{}, id).Error
}
