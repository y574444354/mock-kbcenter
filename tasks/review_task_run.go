package tasks

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/zgsm/go-webserver/pkg/logger"
)

type RunReviewTaskPayload struct {
	ReviewTaskID string `json:"review_task_id"`
}

func NewRunReviewTaskPayload(payload RunReviewTaskPayload, queue string) (*asynq.Task, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeRunReviewTask, payloadBytes, asynq.Queue(queue)), nil
}

func HandleRunReviewTask(ctx context.Context, t *asynq.Task) error {
	var payload RunReviewTaskPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}

	// Start executing review task
	logger.Info("RunReviewTask", "payload", payload)

	return nil
}
