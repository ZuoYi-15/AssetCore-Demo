package data

import (
	"asset-core/internal/infrastructure/kafka"

	"github.com/google/uuid"
)

type Service struct {
	repo     *Repository
	producer kafka.Producer
}

func NewService(repo *Repository, producer kafka.Producer) *Service {
	return &Service{repo: repo, producer: producer}
}

func (s *Service) CreateImportTask(req CreateImportTaskRequest) (*ImportTask, error) {
	task := &ImportTask{
		TaskNo:     "import-" + uuid.NewString(),
		FileName:   req.FileName,
		FileURL:    req.FileURL,
		Status:     ImportPending,
		OperatorID: req.OperatorID,
	}
	if err := s.repo.CreateTask(task); err != nil {
		return nil, err
	}
	_ = s.producer.Publish("data.import.requested", kafka.NewEvent("data.import.requested", task))
	return task, nil
}

func (s *Service) ListImportTasks(offset, limit int) ([]ImportTask, int64, error) {
	return s.repo.ListTasks(offset, limit)
}

func (s *Service) GetImportTask(id uint64) (*ImportTask, error) {
	return s.repo.FindTask(id)
}

func (s *Service) ImportErrors(taskID uint64) ([]ImportError, error) {
	return s.repo.ListErrors(taskID)
}
