package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Has(ctx context.Context, key string) (bool, error)
	Ping() error
}

type ClientConfig struct {
	Host     string
	Port     string
	Password string
}

type client struct {
	rdb *redis.Client
}

func NewClient(config ClientConfig) (Client, error) {
	addr := makeAddr(config)
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: config.Password,
		DB:       0,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return client{rdb}, nil
}

func (c client) Get(ctx context.Context, key string) (string, error) {
	val, err := c.rdb.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func (c client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	err := c.rdb.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c client) Has(ctx context.Context, key string) (bool, error) {
	exists, err := c.rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

func (c client) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := c.rdb.Ping(ctx).Result()
	return err
}

func makeAddr(config ClientConfig) string {
	addr := config.Host
	if config.Port != "" {
		addr += ":" + config.Port
	}
	return addr
}
