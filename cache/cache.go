package cache

import (
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
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
	Addr := viper.GetString("REDDIS_ADDR")
	client := redis.NewClient(&redis.Options{Addr: Addr})
	return client
}
