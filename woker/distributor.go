package worker

import (
	"github.com/hibiken/asynq"
)

type TaskDistributor interface {
	NewEmailDeliveryTask(payload *PayloadSendEmail, opts ...asynq.Option) error
}

type RedisTaskDistributor struct {
	client *asynq.Client
}

func NewRedisDistributor(redisAddr string) TaskDistributor {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	return &RedisTaskDistributor{
		client: client,
	}
}
