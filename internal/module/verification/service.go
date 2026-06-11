package verification

import (
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

func (s *Service) Create(req CreateRequest) (*ResultResponse, error) {
	a, err := s.assetService.Get(req.AssetID)
	if err != nil {
		return nil, err
	}
	task := &Task{
		TaskNo:  "verify-" + uuid.NewString(),
		AssetID: a.ID,
		Status:  StatusPending,
		Result:  ResultUnknown,
	}
	if err := s.repo.CreateTask(task); err != nil {
		return nil, err
	}
	_ = s.producer.Publish("asset.verification.requested", kafka.NewEvent("verification.requested", task))

	score, result, conflicts := verifyAsset(task.ID, *a)
	task.Score = score
	task.Result = result
	task.Status = StatusCompleted
	if err := s.repo.UpdateTask(task); err != nil {
		return nil, err
	}
	if err := s.repo.CreateConflicts(conflicts); err != nil {
		return nil, err
	}
	if a.IdentityID != "" {
		if refreshed, err := s.assetService.RefreshIdentity(a.ID); err != nil {
			return nil, err
		} else {
			a = refreshed
		}
	}
	_ = s.producer.Publish("asset.verification.completed", kafka.NewEvent("verification.completed", task))
	return &ResultResponse{Task: *task, Conflicts: conflicts}, nil
}

func (s *Service) Get(taskID uint64) (*ResultResponse, error) {
	task, err := s.repo.FindTask(taskID)
	if err != nil {
		return nil, err
	}
	conflicts, err := s.repo.ListConflicts(task.ID)
	if err != nil {
		return nil, err
	}
	return &ResultResponse{Task: *task, Conflicts: conflicts}, nil
}

func (s *Service) List(q Query, offset, limit int) ([]TaskRecord, int64, error) {
	return s.repo.ListTasks(q, offset, limit)
}

func (s *Service) LatestByAsset(assetID uint64) (*ResultResponse, error) {
	task, err := s.repo.LatestByAsset(assetID)
	if err != nil {
		return nil, err
	}
	conflicts, err := s.repo.ListConflicts(task.ID)
	if err != nil {
		return nil, err
	}
	return &ResultResponse{Task: *task, Conflicts: conflicts}, nil
}

func verifyAsset(taskID uint64, a asset.Asset) (int, string, []Conflict) {
	score := 0
	conflicts := make([]Conflict, 0)

	addMissing := func(field, value string, weight int) {
		if strings.TrimSpace(value) == "" {
			conflicts = append(conflicts, Conflict{
				TaskID:   taskID,
				AssetID:  a.ID,
				Field:    field,
				Expected: "not empty",
				Actual:   "",
				Severity: "medium",
			})
			return
		}
		score += weight
	}

	addMissing("serial_number", a.SerialNumber, 35)
	addMissing("mac_address", a.MACAddress, 25)
	addMissing("vendor", a.Vendor, 10)
	addMissing("model", a.Model, 10)
	addMissing("owner_department", a.OwnerDepartment, 10)
	addMissing("location", a.Location, 10)

	if a.IdentityID == "" {
		conflicts = append(conflicts, Conflict{
			TaskID:   taskID,
			AssetID:  a.ID,
			Field:    "identity_id",
			Expected: "bound identity",
			Actual:   "",
			Severity: "high",
		})
	}

	result := ResultPassed
	if score < 80 {
		result = ResultWarning
	}
	if score < 50 || len(highConflicts(conflicts)) > 0 {
		result = ResultFailed
	}
	return score, result, conflicts
}

func highConflicts(items []Conflict) []Conflict {
	result := make([]Conflict, 0)
	for _, item := range items {
		if item.Severity == "high" {
			result = append(result, item)
		}
	}
	return result
}
