package asynq

import (
	"context"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/zgsm/mock-kbcenter/config"
	"github.com/zgsm/mock-kbcenter/i18n"
	"github.com/zgsm/mock-kbcenter/pkg/logger"
)

type errorHandler struct {
	logger asynq.Logger
}

func (h *errorHandler) HandleError(ctx context.Context, task *asynq.Task, err error) {
	h.logger.Error("asynq task error",
		"task", task.Type(),
		"error", err,
		"message", i18n.Translate("asynq.server.error.handler", "", nil),
	)
}

var (
	server *asynq.Server
)

// InitServer 初始化Asynq服务器
func InitServer(cfg config.Config) error {
	redisOpt := &asynq.RedisClientOpt{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
		PoolSize: cfg.Asynq.RedisPoolSize,
	}

	// 创建Asynq日志记录器
	// 初始化Asynq专用日志
	if err := logger.InitLogger(cfg.Asynq.Log); err != nil {
		return err
	}

	// 创建Asynq服务器
	server = asynq.NewServer(
		redisOpt,
		asynq.Config{
			Concurrency: cfg.Asynq.Concurrency,
			Queues:      cfg.Asynq.Queues,
			RetryDelayFunc: func(n int, e error, t *asynq.Task) time.Duration {
				return time.Duration(cfg.Asynq.RetryDelay) * time.Second
			},
			Logger: logger.GetAsynqLogger(),
			ErrorHandler: &errorHandler{
				logger: logger.GetAsynqLogger(),
			},
		},
	)

	logger.Info(i18n.Translate("asynq.server.init.success", "", nil))
	return nil
}

// NewServeMux 创建任务处理器路由器
func NewServeMux() *asynq.ServeMux {
	return asynq.NewServeMux()
}

// Start 启动任务处理器
func Start(mux *asynq.ServeMux) error {
	if server == nil {
		msg := i18n.Translate("asynq.server.not.initialized", "", nil)
		return fmt.Errorf("%s", msg)
	}
	return server.Run(mux)
}

// Shutdown 优雅关闭服务器
func Shutdown() {
	if server != nil {
		server.Shutdown()
	}
}

// GetAsynqLogger 获取Asynq日志记录器
func GetAsynqLogger() asynq.Logger {
	return logger.GetAsynqLogger()
}
