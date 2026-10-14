package infrastructure

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type memoryItem struct {
	value     string
	expiresAt time.Time
}

type RedisStore struct {
	client     *redis.Client
	memory     sync.Map
	defaultTTL time.Duration
}

func NewRedisStore(client *redis.Client) *RedisStore {
	return &RedisStore{
		client:     client,
		defaultTTL: 10 * time.Minute,
	}
}

func (r *RedisStore) Exists(ctx context.Context, key string) (bool, error) {
	// Check Redis first if client exists
	if r.client != nil {
		res, err := r.client.Exists(ctx, key).Result()
		if err == nil {
			return res > 0, nil
		}
		log.Printf("[RedisStore] Warning: Redis Exists failed (%v), checking in-memory fallback for key: %s\n", err, key)
	}

	// Fallback to in-memory store
	if val, ok := r.memory.Load(key); ok {
		item, ok := val.(memoryItem)
		if ok && time.Now().Before(item.expiresAt) {
			return true, nil
		}
		r.memory.Delete(key)
	}

	return false, nil
}

func (r *RedisStore) Set(ctx context.Context, key string) error {
	ttl := r.defaultTTL

	// Always save to in-memory fallback
	r.memory.Store(key, memoryItem{
		value:     "pending",
		expiresAt: time.Now().Add(ttl),
	})

	if r.client != nil {
		err := r.client.Set(ctx, key, "pending", ttl).Err()
		if err != nil {
			log.Printf("[RedisStore] Warning: Redis Set failed (%v), stored in in-memory fallback for key: %s\n", err, key)
			// Return nil so notification handler / mediatr pipeline does not fail
			return nil
		}
	}

	return nil
}

func (r *RedisStore) Delete(ctx context.Context, key string) error {
	r.memory.Delete(key)

	if r.client != nil {
		err := r.client.Del(ctx, key).Err()
		if err != nil {
			log.Printf("[RedisStore] Warning: Redis Delete failed (%v)\n", err)
			return nil
		}
	}

	return nil
}
