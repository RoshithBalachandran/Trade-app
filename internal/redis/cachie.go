package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
}

func NewCache(addr string) *Cache {
	rdb := redis.NewClient(&redis.Options{Addr: addr})
	return &Cache{client: rdb}
}

func (c *Cache) Set(ctx context.Context, key string, value []byte) error {
	return c.client.Set(ctx, key, value, 0).Err()
}

func (c *Cache) Get(ctx context.Context, key string) ([]byte, error) {
	val, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	return val, err
}
