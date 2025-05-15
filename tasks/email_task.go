package tasks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/zgsm/go-webserver/pkg/logger"
)

// EmailDeliveryPayload 邮件任务负载
type EmailDeliveryPayload struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

// NewEmailDeliveryTask 创建邮件发送任务
func NewEmailDeliveryTask(payload EmailDeliveryPayload, queue string) (*asynq.Task, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("序列化邮件任务负载失败: %w", err)
	}
	return asynq.NewTask(TypeEmailDelivery, payloadBytes, asynq.Queue(queue)), nil
}

// HandleEmailDeliveryTask 处理邮件发送任务
func HandleEmailDeliveryTask(ctx context.Context, task *asynq.Task) error {
	var payload EmailDeliveryPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("反序列化邮件任务负载失败: %w", err)
	}

	logger.Info("处理邮件发送任务",
		"to", payload.To,
		"subject", payload.Subject,
	)

	// 这里实现实际的邮件发送逻辑
	return nil
}
