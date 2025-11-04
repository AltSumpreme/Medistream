package queue

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/AltSumpreme/Medistream.git/config"
)

var GlobalQueue *RedisQueueConfig

func InitQueue() (*RedisQueueConfig, error) {
	addr := os.Getenv("REDIS_QUEUE_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	q, err := NewRedisQueue(config.Rdb)

	if err != nil {
		return nil, fmt.Errorf("failed to initialize Redis Queue: %w", err)
	}
	GlobalQueue = q

	if err := q.client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}

	log.Println("Redis Queue initialized at", addr)
	return q, nil
}

func (r *RedisQueueConfig) Close() {
	err := r.client.Close()
	if err != nil {
		log.Printf("Error closing Redis client: %v", err)
	} else {
		log.Println("Redis client closed successfully")
	}
}
