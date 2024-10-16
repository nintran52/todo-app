package item

import (
	"time"
	"todo-app/domain"
	"todo-app/pkg/clients"

	"github.com/google/uuid"
)

//go:generate mockery --name ItemRepo
type ItemRepo interface {
	Save(item *domain.ItemCreation) error
	GetAll(filter map[string]any, paging *clients.Paging) ([]domain.Item, error)
	GetItem(filter map[string]any) (domain.Item, error)
	Update(filter map[string]any, item *domain.ItemUpdate) error
	Delete(filter map[string]any) error
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

func (s *itemService) GetAllItem(userID uuid.UUID, paging *clients.Paging) ([]domain.Item, error) {
	filter := map[string]any{"user_id": userID}
	items, err := s.itemRepo.GetAll(filter, paging)
	if err != nil {
		return nil, clients.ErrCannotListEntity(domain.Item{}.TableName(), err)
	}

	return items, nil
}

func (s *itemService) GetItemByID(id, userID uuid.UUID) (domain.Item, error) {
	item, err := s.itemRepo.GetItem(map[string]any{"id": id, "user_id": userID})
	if err != nil {
		return domain.Item{}, clients.ErrCannotGetEntity(item.TableName(), err)
	}

	return item, nil
}

func (s *itemService) UpdateItem(id, userID uuid.UUID, item *domain.ItemUpdate) error {
	item.UpdatedAt = time.Now()
	err := s.itemRepo.Update(map[string]any{"id": id, "user_id": userID}, item)
	if err != nil {
		return clients.ErrCannotUpdateEntity(item.TableName(), err)
	}

	return nil
}

func (s *itemService) DeleteItem(id, userID uuid.UUID) error {
	err := s.itemRepo.Delete(map[string]any{"id": id, "user_id": userID})
	if err != nil {
		return clients.ErrCannotDeleteEntity(domain.Item{}.TableName(), err)
	}

	return nil
}
