package asynq

import (
	"errors"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/zgsm/go-webserver/config"
	"github.com/zgsm/go-webserver/i18n"
	"github.com/zgsm/go-webserver/pkg/logger"
)

var client *asynq.Client

// InitClient 初始化Asynq客户端
func InitClient(cfg config.Config) error {
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

// GetClient 获取Asynq客户端
func GetClient() *asynq.Client {
	return client
}

// Close 关闭Asynq客户端连接
func Close() error {
	if client != nil {
		return client.Close()
	}
	return nil
}

// EnqueueTaskFunc 定义入队任务的函数类型
type EnqueueTaskFunc func(task *asynq.Task, queue string, retryCount ...int) (string, error)

// EnqueueTask 入队任务函数变量
var EnqueueTask EnqueueTaskFunc = func(task *asynq.Task, queue string, retryCount ...int) (string, error) {
	if client == nil {
		logger.Error(i18n.Translate("asynq.client.nil", "", nil))
		return "", errors.New(i18n.Translate("asynq.client.nil", "", nil))
	}

	// 获取重试次数，如果未传入则使用配置中的默认值
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
