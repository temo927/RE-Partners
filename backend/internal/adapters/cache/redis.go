package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"pack-calculator/internal/ports"
	pkgerrors "pack-calculator/pkg/errors"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(addr string, password string, db int) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	return &RedisCache{client: client}, nil
}

func (c *RedisCache) Close() error {
	return c.client.Close()
}

func (c *RedisCache) Get(key string) ([]int, error) {
	ctx := context.Background()
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, pkgerrors.ErrNotFound
	}
	if err != nil {
		return nil, pkgerrors.WrapWithDomain(err, pkgerrors.ErrCache, "failed to get from cache")
	}

	var sizes []int
	if err := json.Unmarshal([]byte(val), &sizes); err != nil {
		return nil, pkgerrors.WrapWithDomain(err, pkgerrors.ErrCache, "failed to unmarshal cache value")
	}

	return sizes, nil
}

func (c *RedisCache) Set(key string, value []int, ttl int) error {
	ctx := context.Background()
	data, err := json.Marshal(value)
	if err != nil {
		return pkgerrors.WrapWithDomain(err, pkgerrors.ErrCache, "failed to marshal cache value")
	}

	if err := c.client.Set(ctx, key, data, time.Duration(ttl)*time.Second).Err(); err != nil {
		return pkgerrors.WrapWithDomain(err, pkgerrors.ErrCache, "failed to set cache")
	}

	return nil
}

func (c *RedisCache) Delete(key string) error {
	ctx := context.Background()
	if err := c.client.Del(ctx, key).Err(); err != nil {
		return pkgerrors.WrapWithDomain(err, pkgerrors.ErrCache, "failed to delete from cache")
	}

	return nil
}

var _ ports.Cache = (*RedisCache)(nil)
