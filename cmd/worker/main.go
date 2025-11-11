package main

import (
	"time"

	"github.com/AltSumpreme/Medistream.git/config"
	"github.com/AltSumpreme/Medistream.git/metrics"
	"github.com/AltSumpreme/Medistream.git/queue"
	"github.com/AltSumpreme/Medistream.git/utils"
	"github.com/AltSumpreme/Medistream.git/workers"
	"github.com/hibiken/asynq"
)

func main() {

	utils.InitLogger()
	// Initialize metrics
	metrics.MetricsInit()

	// Initialize the database connection
	config.ConnectDB()
	defer config.CloseDB()
	// Initialize Job Queue
	config.InitAsynqQueue()
	srv := asynq.NewServer(
		config.QueueRedisOpt,
		asynq.Config{
			Concurrency: 10,

			Queues: map[string]int{
				"appointments": 5,
				"emails":       3,
				//	"reports":      2,
			},

			RetryDelayFunc: func(n int, err error, task *asynq.Task) time.Duration {
				// n is the number of retries already attempted
				// err is the error that caused the retry
				// task is the failed task
				delay := time.Duration(1<<n) * time.Second // Exponential backoff
				if delay > 10*time.Minute {                // Cap the maximum delay
					delay = 10 * time.Minute
				}
				return delay
			},
		},
	)

	mux := asynq.NewServeMux()

	mux.HandleFunc(string(queue.JobTypeCreateAppointment), workers.ProcessCreateAppointmentTask)
	workers.RegisterEmailHandlers(mux)
	//muz.HandleFunc(string(queue.JobTypeGenerateReport),workers.ProcessReportTask);

	if err := srv.Run(mux); err != nil {
		utils.Log.Fatalf("could not run asynq server: %v", err)
	}

}
