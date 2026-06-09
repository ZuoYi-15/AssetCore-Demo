package verification

import (
	"strings"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateTask(task *Task) error {
	return r.db.Create(task).Error
}

func (r *Repository) UpdateTask(task *Task) error {
	return r.db.Save(task).Error
}

func (r *Repository) FindTask(id uint64) (*Task, error) {
	var task Task
	if err := r.db.First(&task, id).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *Repository) ListTasks(q Query, offset, limit int) ([]TaskRecord, int64, error) {
	db := r.db.Table("asset_verification_task").
		Select(`asset_verification_task.id,
			asset_verification_task.task_no,
			asset_verification_task.asset_id,
			asset.asset_name,
			asset.asset_type,
			asset.identity_id,
			asset.serial_number,
			asset.owner_department,
			asset_verification_task.status,
			asset_verification_task.score,
			asset_verification_task.result,
			asset_verification_task.created_at,
			asset_verification_task.updated_at`).
		Joins("LEFT JOIN asset ON asset.id = asset_verification_task.asset_id")
	if strings.TrimSpace(q.Keyword) != "" {
		like := "%" + strings.TrimSpace(q.Keyword) + "%"
		db = db.Where("asset.asset_name LIKE ? OR asset.serial_number LIKE ? OR asset.identity_id LIKE ? OR asset_verification_task.task_no LIKE ?", like, like, like, like)
	}
	if strings.TrimSpace(q.Result) != "" {
		db = db.Where("asset_verification_task.result = ?", strings.TrimSpace(q.Result))
	}
	if strings.TrimSpace(q.Status) != "" {
		db = db.Where("asset_verification_task.status = ?", strings.TrimSpace(q.Status))
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var items []TaskRecord
	if err := db.Order("asset_verification_task.id DESC").Offset(offset).Limit(limit).Scan(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *Repository) LatestByAsset(assetID uint64) (*Task, error) {
	var task Task
	if err := r.db.Where("asset_id = ?", assetID).Order("id DESC").First(&task).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *Repository) CreateConflicts(items []Conflict) error {
	if len(items) == 0 {
		return nil
	}
	return r.db.Create(&items).Error
}

func (r *Repository) ListConflicts(taskID uint64) ([]Conflict, error) {
	var items []Conflict
	err := r.db.Where("task_id = ?", taskID).Order("id ASC").Find(&items).Error
	return items, err
}
