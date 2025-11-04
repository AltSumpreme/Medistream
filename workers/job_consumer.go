package workers

import (
	"context"
	"log"
	"time"

	"github.com/AltSumpreme/Medistream.git/queue"
	"github.com/redis/go-redis/v9"
)

type JobConsumer struct {
	queue     *queue.RedisQueueConfig
	queueName string
	handler   func(*queue.JobPayload) error
}

func NewJobConsumer(q *queue.RedisQueueConfig, queueName string, handler func(*queue.JobPayload) error) *JobConsumer {
	return &JobConsumer{
		queue:     q,
		queueName: queueName,
		handler:   handler,
	}
}

func (c *JobConsumer) Start(ctx context.Context) {
	log.Printf("[Worker] Started consumer for queue: %s", c.queueName)

	for {
		select {
		case <-ctx.Done():
			log.Printf("[Worker] Stopping consumer for %s: context cancelled", c.queueName)
			return

		default:
			job, err := c.queue.Dequeue(ctx, c.queueName)
			if err != nil {
				// Ignore "no job" responses
				if err == redis.Nil {
					time.Sleep(200 * time.Millisecond)
					continue
				}
				log.Printf("[%s] Dequeue error: %v", c.queueName, err)
				time.Sleep(2 * time.Second) // Backoff on error
				continue
			}

			if job == nil {
				time.Sleep(200 * time.Millisecond)
				continue
			}

			start := time.Now()
			if err := c.handler(job); err != nil {
				log.Printf("[%s] Job failed: %v", c.queueName, err)
			} else {
				log.Printf("[%s] Job processed successfully in %v", c.queueName, time.Since(start))
			}
		}
	}
}
