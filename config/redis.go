package config

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client
var Ctx = context.Background()

func InitRedis() {
	host := os.Getenv("REDIS_HOST")

	if host == "" {
		log.Fatal("REDIS_HOST is not set")
	}

	portStr := os.Getenv("REDIS_PORT")
	if portStr == "" {
		log.Fatal("REDIS_PORT is not set")
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Invalid REDIS_PORT: %v", err)
	}
	// Create Redis client
	Rdb = redis.NewClient(&redis.Options{
		Addr:     host + ":" + strconv.Itoa(port),
		Password: "", // default no password set
		DB:       0,
	})
	// Test Redis connection
	_, err = Rdb.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Connected to Redis successfully")
}
