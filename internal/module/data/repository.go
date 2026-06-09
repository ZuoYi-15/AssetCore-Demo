package data

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateTask(task *ImportTask) error {
	return r.db.Create(task).Error
}

func (r *Repository) UpdateTask(task *ImportTask) error {
	return r.db.Save(task).Error
}

func (r *Repository) CreateError(item *ImportError) error {
	return r.db.Create(item).Error
}

func (r *Repository) ListTasks(offset, limit int) ([]ImportTask, int64, error) {
	var total int64
	if err := r.db.Model(&ImportTask{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var items []ImportTask
	err := r.db.Order("id DESC").Offset(offset).Limit(limit).Find(&items).Error
	return items, total, err
}

func (r *Repository) FindTask(id uint64) (*ImportTask, error) {
	var task ImportTask
	if err := r.db.First(&task, id).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *Repository) ListErrors(taskID uint64) ([]ImportError, error) {
	var items []ImportError
	err := r.db.Where("task_id = ?", taskID).Order("row_number ASC").Find(&items).Error
	return items, err
}
