package memcache

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
	"todo-app/domain"

	"github.com/google/uuid"
)

type RealStore interface {
	GetUser(conditions map[string]any) (*domain.User, error)
}

type userCaching struct {
	store     Cache
	realStore RealStore
	once      *sync.Once
}

func NewUserCaching(store Cache, realStore RealStore) *userCaching {
	return &userCaching{
		store:     store,
		realStore: realStore,
		once:      new(sync.Once),
	}
}

func (uc *userCaching) GetUser(conditions map[string]interface{}) (*domain.User, error) {
	var ctx = context.Background()
	var user domain.User

	userId := conditions["id"].(uuid.UUID)
	key := fmt.Sprintf("user-%d", userId)

	err := uc.store.Get(ctx, key, &user)

	if err == nil && user.ID != uuid.Nil {
		return &user, nil
	}

	var userErr error

	uc.once.Do(func() {
		realUser, userErr := uc.realStore.GetUser(conditions)

		if userErr != nil {
			log.Println(userErr)
			return
		}

		// Update cache
		user = *realUser
		_ = uc.store.Set(ctx, key, realUser, time.Hour*2)

	})

	if userErr != nil {
		return nil, userErr
	}

	err = uc.store.Get(ctx, key, &user)

	if err == nil && user.ID != uuid.Nil {
		return &user, nil
	}

	return nil, err
}
