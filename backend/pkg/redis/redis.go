package redis

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func InitRedis() *redis.Client {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}

	password := os.Getenv("REDIS_PASSWORD")

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	// Test connection
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v. Continuing without cache.", err)
		return nil
	}

	log.Println("Redis connection established")
	return client
}

// Cache expiration times
const (
	PortfolioCacheTTL = 30 * time.Second
	HoldingsCacheTTL  = 30 * time.Second
	PriceCacheTTL     = 5 * time.Second
)
