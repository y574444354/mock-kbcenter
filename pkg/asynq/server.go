package asynq

import (
	"context"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/zgsm/go-webserver/config"
	"github.com/zgsm/go-webserver/i18n"
	"github.com/zgsm/go-webserver/pkg/logger"
)

type errorHandler struct {
	logger asynq.Logger
}

func (h *errorHandler) HandleError(ctx context.Context, task *asynq.Task, err error) {
	taskID := ""
	if task.ResultWriter() != nil {
		taskID = task.ResultWriter().TaskID()
	}
	h.logger.Error(i18n.Translate("asynq.server.error.handler", "", map[string]interface{}{
		"taskID":   taskID,
		"taskType": task.Type(),
		"error":    err.Error(),
	}))
}

var (
	server *asynq.Server
)

// InitServer initialize Asynq server
func InitServer(cfg config.Config) error {
	redisOpt := &asynq.RedisClientOpt{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
		PoolSize: cfg.Asynq.RedisPoolSize,
	}

	// Create Asynq logger
	// Initialize Asynq dedicated log
	if err := logger.InitLogger(cfg.Asynq.Log); err != nil {
		return err
	}

	// Create Asynq server
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

// NewServeMux create task handler router
func NewServeMux() *asynq.ServeMux {
	return asynq.NewServeMux()
}

// Start start task handler
func Start(mux *asynq.ServeMux) error {
	if server == nil {
		msg := i18n.Translate("asynq.server.not.initialized", "", nil)
		return fmt.Errorf("%s", msg)
	}
	return server.Run(mux)
}

// Shutdown graceful shutdown server
func Shutdown() {
	if server != nil {
		server.Shutdown()
	}
}

// GetAsynqLogger get Asynq logger
func GetAsynqLogger() asynq.Logger {
	return logger.GetAsynqLogger()
}
