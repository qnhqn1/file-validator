package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/qnhqn1/file-validator/config"
)


type RedisCache struct {
	client *redis.Client
}


type Cache interface {
	Set(ctx context.Context, key string, val []byte, ttl time.Duration) error
	Get(ctx context.Context, key string) ([]byte, error)
	Close()
}


func NewRedisCache(cfg config.RedisConfig) Cache {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	client := redis.NewClient(&redis.Options{Addr: addr, DB: cfg.DB})
	return &RedisCache{client: client}
}


func (r *RedisCache) Set(ctx context.Context, key string, val []byte, ttl time.Duration) error {
	return r.client.Set(ctx, key, val, ttl).Err()
}


func (r *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	res, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}
	return res, nil
}


func (r *RedisCache) Close() {
	_ = r.client.Close()
}


