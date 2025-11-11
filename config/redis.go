package config

import (
	"context"
	"crypto/tls"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

var Rdb *redis.Client
var Ctx = context.Background()

func InitRedis() {
	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Fatalf("failed to parse redis url:%v", err)
	}
	opt.TLSConfig = &tls.Config{
		MinVersion:         tls.VersionTLS12, // enforce secure TLS versions
		InsecureSkipVerify: false,            // ensure proper cert validation
	}
	Rdb = redis.NewClient(opt)

	_, err = Rdb.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("failed to connect to Redis: %v", err)
	}

	log.Println("Connected to Redis successfully")
}
