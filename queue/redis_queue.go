package queue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisQueueConfig struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisQueue(client *redis.Client) (*RedisQueueConfig, error) {
	if client == nil {
		return nil, errors.New("redis client is nil")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("unable to connect to redis: %v", err)
	}

	return &RedisQueueConfig{
		client: client,
		ctx:    context.Background(),
	}, nil
}

func (q *RedisQueueConfig) Enqueue(ctx context.Context, queueName string, payload JobPayload) error {
	if q == nil {
		return errors.New("redis queue is not initialized")
	}
	data, _ := json.Marshal(payload)
	return q.client.LPush(ctx, queueName, data).Err()
}

func (q *RedisQueueConfig) Dequeue(ctx context.Context, queueName string) (*JobPayload, error) {
	data, err := q.client.RPop(ctx, queueName).Result()
	if err != nil {
		return nil, err
	}

	var payload JobPayload
	if err := json.Unmarshal([]byte(data), &payload); err != nil {
		return nil, err
	}
	return &payload, nil
}
