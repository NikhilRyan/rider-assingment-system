package cache

import (
	"github.com/go-redis/redis/v8"
	"log"
	"rider-assignment-system/config"
)

var RedisClient *redis.Client

func InitRedis() {
	cfg := config.Cfg.Redis
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	_, err := RedisClient.Ping(RedisClient.Context()).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	log.Println("Connected to Redis.")
}
