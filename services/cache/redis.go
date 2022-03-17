package cache

import (
	"github.com/go-redis/redis"
	"os"
)

var Client *redis.Client

func Setup() {
	Client = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
}
