package mysql

import (
	"todo-app/domain"
)

type itemRepo struct {
}

func NewItemRepo() *itemRepo {
	return &itemRepo{}
}

func (r *itemRepo) Save(item *domain.ItemCreation) error {
	return nil
}
