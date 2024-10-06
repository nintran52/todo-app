package postgres

import (
	"fmt"
	"todo-app/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type itemRepo struct {
	db *gorm.DB
}

func NewItemRepo(db *gorm.DB) *itemRepo {
	return &itemRepo{
		db: db,
	}
}

func (r *itemRepo) Save(item *domain.ItemCreation) error {
	if err := r.db.Create(&item).Error; err != nil {
		return fmt.Errorf("failed to create item: %w", err)
	}

	return nil
}

func (r *itemRepo) GetAll() ([]domain.Item, error) {
	items := []domain.Item{}

	if err := r.db.Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to get all item: %w", err)
	}

	return items, nil
}

func (r *itemRepo) GetByID(id uuid.UUID) (domain.Item, error) {
	item := domain.Item{}

	if err := r.db.Where("id = ?", id).Find(&item).Error; err != nil {
		return domain.Item{}, fmt.Errorf("failed to get item: %w", err)
	}

	return item, nil
}

func (r *itemRepo) Update(id uuid.UUID, item *domain.ItemUpdate) error {
	if err := r.db.Where("id = ?", id).Updates(&item).Error; err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	return nil
}

func (r *itemRepo) Delete(id uuid.UUID) error {
	if err := r.db.Table(domain.Item{}.TableName()).Where("id = ?", id).Delete(nil).Error; err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}

	return nil
}
