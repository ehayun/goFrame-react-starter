package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"tzlev/internal/redis"
)

type CacheManager struct {
	prefix string
}

func NewCacheManager() *CacheManager {
	return &CacheManager{
		prefix: "tzlev:cache:",
	}
}

func (cm *CacheManager) key(cacheKey string) string {
	return fmt.Sprintf("%s%s", cm.prefix, cacheKey)
}

func (cm *CacheManager) Set(ctx context.Context, cacheKey string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal cache value: %w", err)
	}

	key := cm.key(cacheKey)
	return redis.Client.Set(ctx, key, data, ttl).Err()
}

func (cm *CacheManager) Get(ctx context.Context, cacheKey string, dest interface{}) error {
	key := cm.key(cacheKey)

	data, err := redis.Client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("failed to unmarshal cache value: %w", err)
	}

	return nil
}

func (cm *CacheManager) Delete(ctx context.Context, cacheKey string) error {
	key := cm.key(cacheKey)
	return redis.Client.Del(ctx, key).Err()
}

func (cm *CacheManager) DeletePattern(ctx context.Context, pattern string) error {
	key := cm.key(pattern)

	// Find all keys matching the pattern
	keys, err := redis.Client.Keys(ctx, key).Result()
	if err != nil {
		return err
	}

	if len(keys) == 0 {
		return nil
	}

	// Delete all matching keys
	return redis.Client.Del(ctx, keys...).Err()
}

func (cm *CacheManager) Exists(ctx context.Context, cacheKey string) (bool, error) {
	key := cm.key(cacheKey)
	count, err := redis.Client.Exists(ctx, key).Result()
	return count > 0, err
}
