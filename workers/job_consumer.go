package workers

import (
	"context"
	"log"

	"github.com/AltSumpreme/Medistream.git/queue"
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
	log.Printf("Worker started for queue %s", c.queueName)
	for {
		job, err := c.queue.Dequeue(ctx, c.queueName)
		if err != nil {
			log.Printf("[%s] Dequeue error %v", c.queueName, err)
			continue
		}
		if job == nil {
			continue

		}
		if err := c.handler(job); err != nil {
			log.Printf("[%s]job failed: %v", c.queueName, err)
		}

	}
}
