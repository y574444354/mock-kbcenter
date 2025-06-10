package asynq

import (
	"errors"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/zgsm/mock-kbcenter/config"
	"github.com/zgsm/mock-kbcenter/i18n"
	"github.com/zgsm/mock-kbcenter/pkg/logger"
)

var client *asynq.Client

// InitClient initialize Asynq client
func InitClient(cfg config.Config) error {
	if !cfg.Asynq.Enabled {
		logger.Info(i18n.Translate("asynq.client.disabled", "", nil))
		return nil
	}

	redisOpt := &asynq.RedisClientOpt{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
		PoolSize: cfg.Asynq.RedisPoolSize,
	}

	client = asynq.NewClient(redisOpt)

	logger.Info(i18n.Translate("asynq.client.init.success", "", nil))
	return nil
}

// GetClient get Asynq client
func GetClient() *asynq.Client {
	return client
}

// Close close Asynq client connection
func Close() error {
	if client != nil {
		return client.Close()
	}
	return nil
}

// EnqueueTaskFunc defines function type for enqueueing tasks
type EnqueueTaskFunc func(task *asynq.Task, queue string, retryCount ...int) (string, error)

// EnqueueTask task enqueue function variable
var EnqueueTask EnqueueTaskFunc = func(task *asynq.Task, queue string, retryCount ...int) (string, error) {
	if client == nil {
		logger.Error(i18n.Translate("asynq.client.nil", "", nil))
		return "", errors.New(i18n.Translate("asynq.client.nil", "", nil))
	}

	// Get retry count, use default from config if not provided
	count := config.GetConfig().Asynq.RetryCount
	if len(retryCount) > 0 {
		count = retryCount[0]
	}

	info, err := client.Enqueue(task, asynq.Queue(queue), asynq.MaxRetry(count))
	if err != nil {
		logger.Error(i18n.Translate("asynq.enqueue.failed", "", map[string]interface{}{
			"queue": queue,
			"error": err.Error(),
		}))
		return "", fmt.Errorf("%s", i18n.Translate("asynq.enqueue.failed", "", map[string]interface{}{
			"queue": queue,
			"error": err.Error(),
		}))
	}
	return info.ID, nil
}
