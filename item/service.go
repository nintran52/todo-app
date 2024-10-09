package item

import (
	"time"
	"todo-app/domain"
	"todo-app/pkg/clients"

	"github.com/google/uuid"
)

type ItemRepo interface {
	Save(item *domain.ItemCreation) error
	GetAll() ([]domain.Item, error)
	GetByID(id uuid.UUID) (domain.Item, error)
	Update(id uuid.UUID, item *domain.ItemUpdate) error
	Delete(id uuid.UUID) error
}

type itemService struct {
	itemRepo ItemRepo
}

func NewItemService(repo ItemRepo) *itemService {
	return &itemService{
		itemRepo: repo,
	}
}

func (s *itemService) CreateItem(item *domain.ItemCreation) error {
	if err := item.Validate(); err != nil {
		return clients.ErrInvalidRequest(err)
	}

	item.ID = uuid.New()
	if err := s.itemRepo.Save(item); err != nil {
		return clients.ErrCannotCreateEntity(item.TableName(), err)
	}

	return nil
}

func (s *itemService) GetAllItem() ([]domain.Item, error) {
	items, err := s.itemRepo.GetAll()
	if err != nil {
		return nil, clients.ErrCannotListEntity(domain.Item{}.TableName(), err)
	}

	return items, nil
}

func (s *itemService) GetItemByID(id uuid.UUID) (domain.Item, error) {
	item, err := s.itemRepo.GetByID(id)
	if err != nil {
		return domain.Item{}, clients.ErrCannotGetEntity(item.TableName(), err)
	}

	return item, nil
}

func (s *itemService) UpdateItem(id uuid.UUID, item *domain.ItemUpdate) error {
	item.UpdatedAt = time.Now()
	err := s.itemRepo.Update(id, item)
	if err != nil {
		return clients.ErrCannotUpdateEntity(item.TableName(), err)
	}

	return nil
}

func (s *itemService) DeleteItem(id uuid.UUID) error {
	err := s.itemRepo.Delete(id)
	if err != nil {
		return clients.ErrCannotDeleteEntity(domain.Item{}.TableName(), err)
	}

	return nil
}
