package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"os"
)

var Rdb *redis.Client
var ctx = context.Background()

// InitRedis initializes the Redis connection
func InitRedis() error {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	Rdb = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	// Test the Redis connection
	_, err := Rdb.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %v", err)
	}

	log.Println("Connected to Redis successfully.")
	return nil
}

// GetRedisClient returns the Redis client
func GetRedisClient() *redis.Client {
	return Rdb
}
