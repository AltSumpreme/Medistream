package queue

import (
	"log"
	"os"

	"github.com/AltSumpreme/Medistream.git/config"
)

var GlobalQueue *RedisQueueConfig

func InitQueue() {
	addr := os.Getenv("REDIS_QUEUE_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	GlobalQueue = NewRedisQueue(config.Rdb)

	log.Println("Redis Queue initialized at", addr)
}
