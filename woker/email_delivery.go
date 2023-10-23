package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

const (
	TypeEmailDelivery = "email:deliver"
)

type PayloadSendEmail struct {
	User    string `json:"user"`
	Content string `json:"content"`
}

func (distributor *RedisTaskDistributor) NewEmailDeliveryTask(payload *PayloadSendEmail, opts ...asynq.Option) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}
	task := asynq.NewTask(TypeEmailDelivery, jsonPayload, opts...)
	info, err := distributor.client.Enqueue(task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("queue", info.Queue).Int("max_retry", info.MaxRetry).Msg("enqueued task")

	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	subject := fmt.Sprintf("New ticket from %s", payload.User)
	// TODO: replace this URL with an environment variable that points to a front-end page
	content := fmt.Sprintf(`Hello Cypod engineer,<br/>
	Ticket content: %s <br/>
	Please check it out ASAP <ahref="http://cypod-test.web.app">click here</a>.<br/>`, payload.Content)

	to := []string{viper.GetString("SUPPORT_EMAIL")}

	err := processor.mailer.SendEmail(subject, content, to, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to send verify email: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Msg("processed task")
		// Str("email", user.Email).Msg("processed task")
	return nil
}
