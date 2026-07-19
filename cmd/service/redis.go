package service

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisClient(dsn string) (*RedisClient, error) {
	opt, err := redis.ParseURL(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}
	opt.PoolSize = 10
	opt.MinIdleConns = 10
	opt.MaxRetries = 3
	opt.DialTimeout = time.Second * 5
	opt.ReadTimeout = time.Second * 3
	opt.WriteTimeout = time.Second * 3
	opt.PoolTimeout = time.Second * 4

	client := redis.NewClient(opt)
	ctx := context.Background()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}
	return &RedisClient{
		client: client,
		ctx:    ctx,
	}, nil
}
