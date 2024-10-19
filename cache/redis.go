package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"os"
)

var Rdb *redis.Client

// InitializeRedis initializes the Redis client
func InitializeRedis() error {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")

	if redisHost == "" {
		redisHost = "localhost"
	}
	if redisPort == "" {
		redisPort = "6379"
	}

	addr := fmt.Sprintf("%s:%s", redisHost, redisPort)
	Rdb = redis.NewClient(&redis.Options{
		Addr: addr,
	})

	// Check the Redis connection
	ctx := context.Background()
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
