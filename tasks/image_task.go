package tasks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/zgsm/go-webserver/pkg/logger"
)

// ImageResizePayload 图片处理任务负载
type ImageResizePayload struct {
	SourceURL string `json:"source_url"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
}

// NewImageResizeTask 创建图片处理任务
func NewImageResizeTask(payload ImageResizePayload, queue string) (*asynq.Task, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("序列化图片任务负载失败: %w", err)
	}
	return asynq.NewTask(TypeImageResize, payloadBytes, asynq.Queue(queue)), nil
}

// HandleImageResizeTask 处理图片缩放任务
func HandleImageResizeTask(ctx context.Context, task *asynq.Task) error {
	var payload ImageResizePayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("反序列化图片任务负载失败: %w", err)
	}

	logger.Info("处理图片缩放任务",
		"source_url", payload.SourceURL,
		"width", payload.Width,
		"height", payload.Height,
	)

	// 这里实现实际的图片处理逻辑
	return nil
}
