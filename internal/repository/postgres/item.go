package postgres

import (
	"errors"
	"todo-app/domain"
	"todo-app/pkg/clients"

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

func (r *itemRepo) GetAll(filter map[string]any, paging *clients.Paging) ([]domain.Item, error) {
	items := []domain.Item{}
	var query *gorm.DB

	if f := filter; f != nil {
		if v := f["user_id"]; v != "" {
			query = r.db.Where("user_id = ?", v)
		}
	}

	if err := query.Table(domain.Item{}.TableName()).Select("id").Count(&paging.Total).Error; err != nil {
		return nil, clients.ErrDB(err)
	}

	query = r.db.Limit(paging.Limit).Offset((paging.Page - 1) * paging.Limit)

	if err := query.Find(&items).Error; err != nil {
		return nil, clients.ErrDB(err)
	}

	return items, nil
}

func (r *itemRepo) GetItem(filter map[string]any) (domain.Item, error) {
	var item domain.Item

	if err := r.db.Where(filter).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Item{}, clients.ErrRecordNotFound
		}

		return domain.Item{}, clients.ErrDB(err)
	}

	return item, nil
}

func (r *itemRepo) Update(filter map[string]any, item *domain.ItemUpdate) error {
	if err := r.db.Where(filter).Updates(&item).Error; err != nil {
		return clients.ErrDB(err)
	}

	return nil
}

func (r *itemRepo) Delete(filter map[string]any) error {
	if err := r.db.Table(domain.Item{}.TableName()).Where(filter).Delete(nil).Error; err != nil {
		return clients.ErrDB(err)
	}

	return nil
}
