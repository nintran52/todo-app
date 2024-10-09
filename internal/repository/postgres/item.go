package postgres

import (
	"errors"
	"todo-app/domain"
	"todo-app/pkg/clients"

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
		return clients.ErrDB(err)
	}

	return nil
}

func (r *itemRepo) GetAll() ([]domain.Item, error) {
	items := []domain.Item{}

	if err := r.db.Find(&items).Error; err != nil {
		return nil, clients.ErrDB(err)
	}

	return items, nil
}

func (r *itemRepo) GetByID(id uuid.UUID) (domain.Item, error) {
	var item domain.Item

	if err := r.db.Where("id = ?", id).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Item{}, clients.ErrRecordNotFound
		}

		return domain.Item{}, clients.ErrDB(err)
	}

	return item, nil
}

func (r *itemRepo) Update(id uuid.UUID, item *domain.ItemUpdate) error {
	if err := r.db.Where("id = ?", id).Updates(&item).Error; err != nil {
		return clients.ErrDB(err)
	}

	return nil
}

func (r *itemRepo) Delete(id uuid.UUID) error {
	if err := r.db.Table(domain.Item{}.TableName()).Where("id = ?", id).Delete(nil).Error; err != nil {
		return clients.ErrDB(err)
	}

	return nil
}
