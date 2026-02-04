package cache

import (
	"context"
	"errors"
	"testing"
	"time"

	pkgerrors "pack-calculator/pkg/errors"
)

func setupTestRedis(t *testing.T) *RedisCache {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	cache, err := NewRedisCache("localhost:6379", "", 0)
	if err != nil {
		t.Skipf("Skipping test: failed to connect to Redis: %v", err)
	}

	// Clean test keys
	ctx := context.Background()
	cache.client.FlushDB(ctx)

	return cache
}

func TestRedisCache_Get(t *testing.T) {
	cache := setupTestRedis(t)
	defer cache.Close()

	t.Run("key not found returns ErrNotFound", func(t *testing.T) {
		_, err := cache.Get("nonexistent")
		if err == nil {
			t.Error("Get() error = nil, want ErrNotFound")
		}
		if !errors.Is(err, pkgerrors.ErrNotFound) {
			t.Errorf("Get() error should be ErrNotFound, got: %v", err)
		}
	})

	t.Run("get existing key", func(t *testing.T) {
		key := "test:get"
		value := []int{250, 500, 1000}

		err := cache.Set(key, value, 60)
		if err != nil {
			t.Fatalf("Set() error = %v", err)
		}

		got, err := cache.Get(key)
		if err != nil {
			t.Errorf("Get() error = %v", err)
		}
		if len(got) != len(value) {
			t.Errorf("Get() = %v, want %v", got, value)
		}
	})
}

func TestRedisCache_Set(t *testing.T) {
	cache := setupTestRedis(t)
	defer cache.Close()

	t.Run("set and get", func(t *testing.T) {
		key := "test:set"
		value := []int{100, 200, 300}

		err := cache.Set(key, value, 60)
		if err != nil {
			t.Errorf("Set() error = %v", err)
		}

		got, err := cache.Get(key)
		if err != nil {
			t.Errorf("Get() error = %v", err)
		}
		if len(got) != len(value) {
			t.Errorf("Get() = %v, want %v", got, value)
		}
	})

	t.Run("set with TTL", func(t *testing.T) {
		key := "test:ttl"
		value := []int{1, 2, 3}

		err := cache.Set(key, value, 1)
		if err != nil {
			t.Fatalf("Set() error = %v", err)
		}

		// Wait for expiration
		time.Sleep(2 * time.Second)

		_, err = cache.Get(key)
		if err == nil {
			t.Error("Get() after expiration error = nil, want error")
		}
	})
}

func TestRedisCache_Delete(t *testing.T) {
	cache := setupTestRedis(t)
	defer cache.Close()

	t.Run("delete existing key", func(t *testing.T) {
		key := "test:delete"
		value := []int{250, 500}

		err := cache.Set(key, value, 60)
		if err != nil {
			t.Fatalf("Set() error = %v", err)
		}

		err = cache.Delete(key)
		if err != nil {
			t.Errorf("Delete() error = %v", err)
		}

		_, err = cache.Get(key)
		if err == nil {
			t.Error("Get() after delete error = nil, want error")
		}
	})

	t.Run("delete nonexistent key", func(t *testing.T) {
		err := cache.Delete("nonexistent")
		if err != nil {
			t.Errorf("Delete() nonexistent key error = %v, want nil", err)
		}
	})
}

func TestRedisCache_ErrorWrapping(t *testing.T) {
	t.Run("invalid address returns error", func(t *testing.T) {
		_, err := NewRedisCache("invalid:6379", "", 0)
		if err == nil {
			t.Error("NewRedisCache() error = nil, want error")
		}
	})

	t.Run("cache errors are wrapped", func(t *testing.T) {
		cache := setupTestRedis(t)
		defer cache.Close()

		// Force an error by using invalid data
		ctx := context.Background()
		cache.client.Set(ctx, "invalid", "not json", time.Hour)

		_, err := cache.Get("invalid")
		if err != nil {
			// Error should be wrapped with ErrCache
			if !errors.Is(err, pkgerrors.ErrCache) {
				t.Errorf("Get() error should be wrapped with ErrCache, got: %v", err)
			}
		}
	})
}
