package queue

import (
	"github.com/AltSumpreme/Medistream.git/config"
	"github.com/hibiken/asynq"
)

var Client *asynq.Client

func Init() *asynq.Client {
	Client = asynq.NewClient(config.QueueRedisOpt)
	return Client
}

func Close() error {
	if Client != nil {
		return Client.Close()
	}
	return nil
}
