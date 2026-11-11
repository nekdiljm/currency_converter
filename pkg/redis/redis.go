package redis

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Set(ctx context.Context, key string, v any, ttl time.Duration) error
	Get(ctx context.Context, key string, v any) error
	Exists(ctx context.Context, key string) (bool, error)
	Delete(ctx context.Context, key string) error
	Close() error
}

type cache struct {
	client *redis.Client
}

func New(ctx context.Context, host, port string) (Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr: net.JoinHostPort(host, port),
	})

	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("ping redis: %w", err)
	}
	return &cache{client: client}, nil
}

func (c *cache) Set(ctx context.Context, key string, v any, ttl time.Duration) error {
	return c.client.Set(ctx, key, v, ttl).Err()
}

func (c *cache) Get(ctx context.Context, key string, v any) error {
	return c.client.Get(ctx, key).Scan(&v)
}

func (c *cache) Exists(ctx context.Context, key string) (bool, error) {
	n, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

func (c *cache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

func (c *cache) Close() error {
	return c.client.Close()
}
