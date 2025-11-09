package config

import (
	"log"
	"os"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
)

var QueueRedisOpt asynq.RedisClientOpt

func InitAsynqQueue() {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://redis:6379/0" // default
	}

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatalf("Invalid REDIS_URL: %v", err)
	}

	QueueRedisOpt = asynq.RedisClientOpt{
		Addr:     opt.Addr,
		DB:       opt.DB,
		Password: opt.Password,
	}

	log.Printf("Asynq queue initialized: %s DB=%d", opt.Addr, opt.DB)
}
