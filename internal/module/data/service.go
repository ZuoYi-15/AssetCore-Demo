package data

import (
	"encoding/json"
	"errors"
	"mime/multipart"
	"strings"

	"asset-core/internal/infrastructure/kafka"
	"asset-core/internal/module/asset"

	"github.com/google/uuid"
)

type Service struct {
	repo         *Repository
	assetService *asset.Service
	producer     kafka.Producer
}

func NewService(repo *Repository, assetService *asset.Service, producer kafka.Producer) *Service {
	return &Service{repo: repo, assetService: assetService, producer: producer}
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

func (s *Service) ImportAssetsExcel(file multipart.File, filename, operatorID string) (*ImportAssetsResult, error) {
	if s.assetService == nil {
		return nil, errors.New("asset service is not configured")
	}
	rows, err := ReadAssetRows(file)
	if err != nil {
		return nil, err
	}
	task := &ImportTask{
		TaskNo:     "import-" + uuid.NewString(),
		FileName:   filename,
		FileURL:    "upload://" + filename,
		Status:     ImportRunning,
		OperatorID: operatorID,
		TotalCount: len(rows),
	}
	if err := s.repo.CreateTask(task); err != nil {
		return nil, err
	}

	for _, row := range rows {
		req := row.Request
		if strings.TrimSpace(req.Source) == "" {
			req.Source = "excel"
		}
		if strings.TrimSpace(req.AssetName) == "" {
			s.recordImportError(task, row.RowNumber, "asset_name", "asset_name is required", row.Raw)
			continue
		}
		if _, err := s.assetService.Create(req); err != nil {
			s.recordImportError(task, row.RowNumber, "asset", err.Error(), row.Raw)
			continue
		}
		task.SuccessCount++
	}

	task.FailedCount = task.TotalCount - task.SuccessCount
	if task.FailedCount > 0 {
		task.Status = ImportFailed
	} else {
		task.Status = ImportCompleted
	}
	if err := s.repo.UpdateTask(task); err != nil {
		return nil, err
	}
	_ = s.producer.Publish("data.import.assets.completed", kafka.NewEvent("data.import.assets.completed", task))
	return &ImportAssetsResult{Task: task}, nil
}

func (s *Service) recordImportError(task *ImportTask, rowNumber int, field, message string, raw map[string]string) {
	task.FailedCount++
	payload, _ := json.Marshal(raw)
	_ = s.repo.CreateError(&ImportError{
		TaskID:       task.ID,
		RowNumber:    rowNumber,
		ErrorField:   field,
		ErrorMessage: message,
		RawData:      string(payload),
	})
}
