package queue

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

type RedisQueueConfig struct {
	client *redis.Client
}

func NewRedisQueue(client *redis.Client) *RedisQueueConfig {
	return &RedisQueueConfig{
		client: client,
	}
}

func (q *RedisQueueConfig) Enqueue(ctx context.Context, queueName string, payload JobPayload) error {
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
