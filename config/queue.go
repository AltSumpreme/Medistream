package config

import (
	"crypto/tls"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/hibiken/asynq"
)

var QueueRedisOpt asynq.RedisClientOpt

func InitAsynqQueue() {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://redis:6379/0" // default
	}
	opt.TLSConfig = &tls.Config{
		MinVersion:         tls.VersionTLS12, // enforce secure TLS versions
		InsecureSkipVerify: false,            // ensure proper cert validation
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
