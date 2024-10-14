package memcache

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
)

type redisCache struct {
	store *cache.Cache
}

func NewRedisCache() *redisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: "",
		DB:       0,
	})

	var ctx = context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Can not connect to Redis: %v", err)
	}
	fmt.Println("Connected to Redis")

	c := cache.New(&cache.Options{
		Redis:      rdb,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	return &redisCache{store: c}
}

func (rdc *redisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return rdc.store.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: value,
		TTL:   0,
	})
}

func (rdc *redisCache) Get(ctx context.Context, key string, value interface{}) error {
	return rdc.store.Get(ctx, key, value)
}

func (rdc *redisCache) Delete(ctx context.Context, key string) error {
	return rdc.store.Delete(ctx, key)
}
