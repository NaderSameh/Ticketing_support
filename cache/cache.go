package cache

import (
	"github.com/redis/go-redis/v9"
)

type CacheClient struct {
	RedisClient *redis.Client
}

func NewCacheClient() *CacheClient {
	return &CacheClient{
		RedisClient: NewRedisClient(),
	}
}

func NewRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{Addr: ":6379"})
	return client
}
