package workers

import (
	"context"

	"github.com/AltSumpreme/Medistream.git/queue"
)

func StartAllWorkers(ctx context.Context, queue *queue.RedisQueueConfig) {

	go NewJobConsumer(queue, "appointment_jobs", HandleAppointmentJobs).Start(ctx)
}
