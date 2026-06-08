package verification

import "gorm.io/gorm"

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
