package db

import (
	"context"
	"dex-tool/src/configs/env"
	"github.com/go-redis/redis/v8"
	"sync"
	"time"
)

type IRedisClient interface {
	Get(key string) (string, error)
	Set(key string, value string, expiration time.Duration) error
	SetNX(key string, value string, expiration time.Duration) (bool, error)
	IncrBy(key string, value int64) (int64, error)
}

type redisClient struct {
	client *redis.Client
	ctx    context.Context
}

var (
	redisClientInstance *redisClient
	once                sync.Once
)

func GetRedisInstance() IRedisClient {
	once.Do(func() {
		redisClientInstance = &redisClient{
			client: redis.NewClient(&redis.Options{
				Addr:     env.ConfigEnv.Redis.Address,
				Username: env.ConfigEnv.Redis.Username,
				Password: env.ConfigEnv.Redis.Password,
				DB:       env.ConfigEnv.Redis.Db,
			}),
			ctx: context.Background(),
		}
	})
	return redisClientInstance
}

func (r *redisClient) Get(key string) (string, error) {
	return r.client.Get(r.ctx, key).Result()
}

func (r *redisClient) Set(key string, value string, expiration time.Duration) error {
	return r.client.Set(r.ctx, key, value, expiration).Err()
}

func (r *redisClient) SetNX(key string, value string, expiration time.Duration) (bool, error) {
	return r.client.SetNX(r.ctx, key, value, expiration).Result()
}

func (r *redisClient) IncrBy(key string, value int64) (int64, error) {
	return r.client.IncrBy(r.ctx, key, value).Result()
}
