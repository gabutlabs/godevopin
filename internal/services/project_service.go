package service

import (
	"github.com/gabutlabs/devopin/internal/model"
	"github.com/gabutlabs/devopin/internal/repository"
)

type ProjectService interface {
	CreateProject(name, path, projectType, logFormat string) error
	GetProjectByID(id uint) (*model.Project, error)
	ListProjects() ([]model.Project, error)
	UpdateProject(project *model.Project) error
	DeleteProject(id uint) error
}

type projectService struct {
	repository repository.ProjectRepository
}

func NewProjectService(repo repository.ProjectRepository) ProjectService {
	return &projectService{repository: repo}
}

func (s *projectService) CreateProject(name, path, projectType, logFormat string) error {
	project := &model.Project{
		Name:        name,
		Path:        path,
		ProjectType: projectType,
		LogFormat:   logFormat,
	}
	return s.repository.CreateProject(project)
}

func (s *projectService) GetProjectByID(id uint) (*model.Project, error) {
	return s.repository.GetProjectByID(id)
}

func (s *projectService) ListProjects() ([]model.Project, error) {
	return s.repository.ListProjects()
}

func (s *projectService) UpdateProject(project *model.Project) error {
	return s.repository.UpdateProject(project)
}

func (s *projectService) DeleteProject(id uint) error {
	return s.repository.DeleteProject(id)
}
