package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/AltSumpreme/Medistream.git/config"
	"github.com/AltSumpreme/Medistream.git/metrics"
	"github.com/AltSumpreme/Medistream.git/queue"
	"github.com/AltSumpreme/Medistream.git/utils"
	"github.com/AltSumpreme/Medistream.git/workers"
)

func main() {

	utils.InitLogger()

	// Initialize Redis
	config.InitRedis()
	utils.Log.Info("Worker started")

	// Initialize metrics
	metrics.MetricsInit()

	// Initialize the database connection
	config.ConnectDB()
	defer config.CloseDB()
	// Initialize Job Queue

	q, err := queue.InitQueue()
	if err != nil {
		utils.Log.Fatalf("Failed to initialize queue: %v", err)
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	workers.StartAllWorkers(ctx, q)

	defer q.Close()

	select {}

}
