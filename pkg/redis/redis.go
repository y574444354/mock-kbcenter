package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/zgsm/review-manager/config"
	"github.com/zgsm/review-manager/i18n"
	"github.com/zgsm/review-manager/pkg/logger"
)

var (
	client *redis.Client
	ctx    = context.Background()
)

// InitRedis 初始化Redis连接
func InitRedis(cfg config.Config) error {
	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
		// 连接池配置
		PoolSize:     10,
		MinIdleConns: 5,
		// 超时配置
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolTimeout:  4 * time.Second,
	})

	// 测试连接
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf(i18n.Translate("redis.connect.failed", "", nil)+": %w", err)
	}

	logger.Info(i18n.Translate("redis.connect.success", "", nil), "addr", fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port))
	return nil
}

// GetClient 获取Redis客户端
func GetClient() *redis.Client {
	return client
}

// Close 关闭Redis连接
func Close() error {
	if client != nil {
		return client.Close()
	}
	return nil
}

// Set 设置键值对
func Set(key string, value interface{}, expiration time.Duration) error {
	return client.Set(ctx, key, value, expiration).Err()
}

// Get 获取值
func Get(key string) (string, error) {
	return client.Get(ctx, key).Result()
}

// Del 删除键
func Del(keys ...string) error {
	return client.Del(ctx, keys...).Err()
}

// Exists 检查键是否存在
func Exists(keys ...string) (bool, error) {
	result, err := client.Exists(ctx, keys...).Result()
	return result > 0, err
}

// Expire 设置过期时间
func Expire(key string, expiration time.Duration) error {
	return client.Expire(ctx, key, expiration).Err()
}

// TTL 获取剩余过期时间
func TTL(key string) (time.Duration, error) {
	return client.TTL(ctx, key).Result()
}

// Incr 自增
func Incr(key string) (int64, error) {
	return client.Incr(ctx, key).Result()
}

// HSet 设置哈希表字段值
func HSet(key string, field string, value interface{}) error {
	return client.HSet(ctx, key, field, value).Err()
}

// HGet 获取哈希表字段值
func HGet(key, field string) (string, error) {
	return client.HGet(ctx, key, field).Result()
}

// HGetAll 获取哈希表所有字段和值
func HGetAll(key string) (map[string]string, error) {
	return client.HGetAll(ctx, key).Result()
}

// HDel 删除哈希表字段
func HDel(key string, fields ...string) error {
	return client.HDel(ctx, key, fields...).Err()
}
