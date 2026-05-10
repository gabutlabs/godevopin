package service

import (
	"github.com/gabutlabs/devopin/internal/model"
	"github.com/gabutlabs/devopin/internal/repository"
)

type LogHistoryService interface {
	BatchInsert(logs []model.LogHistory) error
	GetLogsByProjectID(projectID uint, limit int, offset int) ([]model.LogHistory, int64, error)
}

type logHistoryService struct {
	repo repository.LogHistoryRepository
}

func NewLogHistoryService(repo repository.LogHistoryRepository) LogHistoryService {
	return &logHistoryService{repo: repo}
}

func (s *logHistoryService) BatchInsert(logs []model.LogHistory) error {
	return s.repo.BatchInsert(logs)
}

func (s *logHistoryService) GetLogsByProjectID(projectID uint, limit int, offset int) ([]model.LogHistory, int64, error) {
	return s.repo.GetLogsByProjectID(projectID, limit, offset)
}
