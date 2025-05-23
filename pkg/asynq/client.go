package asynq

import (
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/zgsm/review-manager/config"
	"github.com/zgsm/review-manager/i18n"
	"github.com/zgsm/review-manager/pkg/logger"
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
type EnqueueTaskFunc func(task *asynq.Task, queue string) (string, error)

// EnqueueTask 入队任务函数变量
var EnqueueTask EnqueueTaskFunc = func(task *asynq.Task, queue string) (string, error) {
	info, err := client.Enqueue(task, asynq.Queue(queue))
	if err != nil {
		return "", err
	}
	return info.ID, nil
}
