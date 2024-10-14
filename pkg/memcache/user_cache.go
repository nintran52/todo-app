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
	mutex     sync.Mutex
}

func NewUserCaching(store Cache, realStore RealStore) *userCaching {
	return &userCaching{
		store:     store,
		realStore: realStore,
	}
}

func (uc *userCaching) GetUser(conditions map[string]interface{}) (*domain.User, error) {
	var ctx = context.Background()
	var user domain.User

	// Safely extract userId with comma-ok pattern to avoid panics
	userId, ok := conditions["id"].(uuid.UUID)
	if !ok {
		return nil, fmt.Errorf("invalid user ID type")
	}

	key := fmt.Sprintf("user-%d", userId)

	// Try to get the user from the cache
	if err := uc.store.Get(ctx, key, &user); err == nil && user.ID != uuid.Nil {
		return &user, nil
	}

	// Use a sync.Mutex or singleflight group to prevent duplicate requests
	uc.mutex.Lock()
	defer uc.mutex.Unlock()

	// Double-check cache after acquiring the lock
	if err := uc.store.Get(ctx, key, &user); err == nil && user.ID != uuid.Nil {
		return &user, nil
	}

	// If not in cache, get the user from the real store
	realUser, err := uc.realStore.GetUser(conditions)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Update the cache with the retrieved user
	if cacheErr := uc.store.Set(ctx, key, realUser, time.Hour*2); cacheErr != nil {
		log.Printf("failed to set cache: %v", cacheErr)
	}

	return realUser, nil
}
