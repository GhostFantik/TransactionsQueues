package redis_client

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

type ConfigRedis struct {
	Host     string
	Port     string
	Password string
	DB       int
}

func NewRedisDB(cfg ConfigRedis) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: "",
		DB:       0,
	})
	_, err := rdb.Ping(context.TODO()).Result()
	if err != nil {
		return nil, err
	}
	return rdb, nil
}
