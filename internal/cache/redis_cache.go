// internal/cache/redis_cache.go
package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"cloud-disk/internal/config"
	"cloud-disk/pkg/logger"

	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisCache(cfg *config.RedisConfig) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	ctx := context.Background()

	// 测试连接
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	logger.Info("Redis cache connected")
	return &RedisCache{
		client: client,
		ctx:    ctx,
	}, nil
}

// Set 设置缓存
func (c *RedisCache) Set(key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(c.ctx, key, data, expiration).Err()
}

// Get 获取缓存
func (c *RedisCache) Get(key string, dest interface{}) error {
	data, err := c.client.Get(c.ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

// Delete 删除缓存
func (c *RedisCache) Delete(key string) error {
	return c.client.Del(c.ctx, key).Err()
}

// Close 关闭连接
func (c *RedisCache) Close() error {
	return c.client.Close()
}
