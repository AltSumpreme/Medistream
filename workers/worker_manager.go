package workers

import (
	"context"

	"github.com/AltSumpreme/Medistream.git/queue"
)

func StartAllWorkers(ctx context.Context, queue *queue.RedisQueueConfig) {

	// go NewJobConsumer(queue, "auth_jobs", HandleAuthJobs).Start(ctx)
}
