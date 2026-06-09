package workflow

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) AutoMigrate() error {
	return r.db.AutoMigrate(&Definition{}, &Node{}, &Instance{}, &Task{})
}

func (r *Repository) Transaction(fn func(*Repository) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return fn(&Repository{db: tx})
	})
}

func (r *Repository) CreateDefinition(item *Definition) error {
	return r.db.Create(item).Error
}

func (r *Repository) CountDefinitions() (int64, error) {
	var count int64
	err := r.db.Model(&Definition{}).Count(&count).Error
	return count, err
}

func (r *Repository) UpdateDefinition(item *Definition) error {
	return r.db.Save(item).Error
}

func (r *Repository) DeleteNodes(definitionID uint64) error {
	return r.db.Where("definition_id = ?", definitionID).Delete(&Node{}).Error
}

func (r *Repository) CreateNodes(nodes []Node) error {
	if len(nodes) == 0 {
		return nil
	}
	return r.db.Create(&nodes).Error
}

func (r *Repository) FindDefinitionByType(flowType string) (*Definition, error) {
	var item Definition
	if err := r.db.Preload("Nodes", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order ASC, id ASC")
	}).Where("flow_type = ?", flowType).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *Repository) FindDefinitionByID(id uint64) (*Definition, error) {
	var item Definition
	if err := r.db.Preload("Nodes", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order ASC, id ASC")
	}).First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *Repository) ListDefinitions() ([]Definition, error) {
	var items []Definition
	err := r.db.Preload("Nodes", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order ASC, id ASC")
	}).Order("id ASC").Find(&items).Error
	return items, err
}

func (r *Repository) DeleteDefinition(id uint64) error {
	return r.db.Delete(&Definition{}, id).Error
}

func (r *Repository) CreateInstance(item *Instance) error {
	return r.db.Create(item).Error
}

func (r *Repository) UpdateInstance(item *Instance) error {
	return r.db.Save(item).Error
}

func (r *Repository) CreateTask(item *Task) error {
	return r.db.Create(item).Error
}

func (r *Repository) UpdateTask(item *Task) error {
	return r.db.Save(item).Error
}

func (r *Repository) FindTask(id uint64) (*Task, error) {
	var task Task
	err := r.db.Preload("Instance").First(&task, id).Error
	return &task, err
}

func (r *Repository) FindInstance(id uint64) (*Instance, error) {
	var item Instance
	err := r.db.Preload("Tasks", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order ASC, id ASC")
	}).First(&item, id).Error
	return &item, err
}

func (r *Repository) ListInstances(assetID uint64, offset, limit int) ([]Instance, int64, error) {
	db := r.db.Model(&Instance{})
	if assetID > 0 {
		db = db.Where("asset_id = ?", assetID)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var items []Instance
	err := db.Preload("Tasks", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order ASC, id ASC")
	}).Order("id DESC").Offset(offset).Limit(limit).Find(&items).Error
	return items, total, err
}

func (r *Repository) ListApprovedRetireInstances() ([]Instance, error) {
	var items []Instance
	err := r.db.Where("flow_type = ? AND status = ? AND asset_id > 0", FlowRetire, InstanceApproved).Find(&items).Error
	return items, err
}

func (r *Repository) ListApprovedTransferInstances() ([]Instance, error) {
	var items []Instance
	err := r.db.Where("flow_type = ? AND status = ? AND asset_id > 0", FlowTransfer, InstanceApproved).Find(&items).Error
	return items, err
}

func (r *Repository) ListTasks(status string, offset, limit int) ([]Task, int64, error) {
	db := r.db.Model(&Task{})
	if status != "" {
		db = db.Where("status = ?", status)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var items []Task
	err := db.Preload("Instance").Order("id DESC").Offset(offset).Limit(limit).Find(&items).Error
	return items, total, err
}
