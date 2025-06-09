package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/zgsm/go-webserver/config"
	"github.com/zgsm/go-webserver/i18n"
	"github.com/zgsm/go-webserver/pkg/logger"
)

var (
	client *redis.Client
	ctx    = context.Background()
)

// InitRedis initialize Redis connection
func InitRedis(cfg config.Config) error {
	if !cfg.Redis.Enabled {
		return nil
	}
	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
		// Connection pool configuration
		PoolSize:     10,
		MinIdleConns: 5,
		// Timeout configuration
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolTimeout:  4 * time.Second,
	})

	// Test connection
	_, err := client.Ping(ctx).Result()
	if err != nil {
		// Close client if connection failed
		if closeErr := client.Close(); closeErr != nil {
			logger.Error(i18n.Translate("redis.client.close.failed", "", map[string]interface{}{"error": closeErr}))
		}
		return fmt.Errorf(i18n.Translate("redis.connect.failed", "", nil)+": %w", err)
	}

	logger.Info(i18n.Translate("redis.connect.success", "", nil), "addr", fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port))
	return nil
}

// GetClient get Redis client
func GetClient() (*redis.Client, error) {
	if client == nil {
		return nil, fmt.Errorf("%s", i18n.Translate("redis.not_initialized_or_disabled", "", nil))
	}
	return client, nil
}

// TryLock attempt to acquire distributed lock
func TryLock(key string, expiration time.Duration) (bool, error) {
	if client == nil {
		return false, fmt.Errorf("%s", i18n.Translate("redis.not_initialized_or_disabled", "", nil))
	}
	return client.SetNX(ctx, key, 1, expiration).Result()
}

// Unlock release distributed lock
func Unlock(key string) error {
	if client == nil {
		return fmt.Errorf("%s", i18n.Translate("redis.not_initialized_or_disabled", "", nil))
	}
	return client.Del(ctx, key).Err()
}

// Close close Redis connection
func Close() error {
	if client != nil {
		return client.Close()
	}
	return nil
}

// Set set key-value pair
func Set(key string, value interface{}, expiration time.Duration) error {
	if client == nil {
		return fmt.Errorf("%s", i18n.Translate("redis.not_initialized_or_disabled", "", nil))
	}
	return client.Set(ctx, key, value, expiration).Err()
}

// Get get value
func Get(key string) (string, error) {
	if client == nil {
		return "", fmt.Errorf("%s", i18n.Translate("redis.not_initialized_or_disabled", "", nil))
	}
	return client.Get(ctx, key).Result()
}

// Del delete keys
func Del(keys ...string) error {
	if client == nil {
		return fmt.Errorf("%s", i18n.Translate("redis.not_initialized_or_disabled", "", nil))
	}
	return client.Del(ctx, keys...).Err()
}

// Exists check if key exists
func Exists(keys ...string) (bool, error) {
	if client == nil {
		return false, fmt.Errorf("%s", i18n.Translate("redis.not_initialized_or_disabled", "", nil))
	}
	result, err := client.Exists(ctx, keys...).Result()
	return result > 0, err
}

// Expire set expiration time
func Expire(key string, expiration time.Duration) error {
	if client == nil {
		return fmt.Errorf("%s", i18n.Translate("redis.not_initialized_or_disabled", "", nil))
	}
	return client.Expire(ctx, key, expiration).Err()
}

// TTL get time to live
func TTL(key string) (time.Duration, error) {
	if client == nil {
		return 0, fmt.Errorf("%s", i18n.Translate("redis.not_initialized_or_disabled", "", nil))
	}
	return client.TTL(ctx, key).Result()
}

// Incr increment value
func Incr(key string) (int64, error) {
	if client == nil {
		return 0, fmt.Errorf("%s", i18n.Translate("redis.not_initialized_or_disabled", "", nil))
	}
	return client.Incr(ctx, key).Result()
}

// HSet set hash field value
func HSet(key string, field string, value interface{}) error {
	if client == nil {
		return fmt.Errorf("%s", i18n.Translate("redis.not_initialized_or_disabled", "", nil))
	}
	return client.HSet(ctx, key, field, value).Err()
}

// HGet get hash field value
func HGet(key, field string) (string, error) {
	if client == nil {
		return "", fmt.Errorf("%s", i18n.Translate("redis.not_initialized_or_disabled", "", nil))
	}
	return client.HGet(ctx, key, field).Result()
}

// HGetAll get all hash fields and values
func HGetAll(key string) (map[string]string, error) {
	if client == nil {
		return nil, fmt.Errorf("%s", i18n.Translate("redis.not_initialized_or_disabled", "", nil))
	}
	return client.HGetAll(ctx, key).Result()
}

// HDel delete hash fields
func HDel(key string, fields ...string) error {
	if client == nil {
		return fmt.Errorf("%s", i18n.Translate("redis.not_initialized_or_disabled", "", nil))
	}
	return client.HDel(ctx, key, fields...).Err()
}

// FlushDB flush current database
func FlushDB() error {
	if client == nil {
		return fmt.Errorf("%s", i18n.Translate("redis.not_initialized_or_disabled", "", nil))
	}
	if err := client.FlushDB(ctx).Err(); err != nil {
		return fmt.Errorf(i18n.Translate("redis.flushdb.failed", "", nil)+": %w", err)
	}
	logger.Info(i18n.Translate("redis.flushdb.success", "", nil))
	return nil
}
