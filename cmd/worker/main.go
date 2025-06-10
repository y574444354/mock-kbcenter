package worker

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/zgsm/mock-kbcenter/config"
	"github.com/zgsm/mock-kbcenter/i18n"
	"github.com/zgsm/mock-kbcenter/pkg/asynq"
	"github.com/zgsm/mock-kbcenter/pkg/db"
	"github.com/zgsm/mock-kbcenter/pkg/logger"
	"github.com/zgsm/mock-kbcenter/pkg/redis"
	"github.com/zgsm/mock-kbcenter/pkg/thirdPlatform"
	"github.com/zgsm/mock-kbcenter/tasks"
)

func Run(cfg *config.Config) {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				logger.Error("worker panic", "error", err)
			}
			logger.Sync()
			os.Exit(1)
		}
	}()

	// Check if Asynq is enabled
	if !cfg.Asynq.Enabled {
		log.Println(i18n.Translate("asynq.server.disabled", "", nil))
		return
	}

	// Initialize logger
	if err := logger.InitLogger(cfg.Asynq.Log); err != nil {
		logger.Error(i18n.Translate("asynq.server.init.failed", "", nil), "error", err)
		panic(err)
	}
	defer logger.Sync()

	// Initialize database
	if err := db.InitDB(cfg.Database); err != nil {
		logger.Error(i18n.Translate("db.init.failed", "", nil), "error", err)
		panic(err)
	}

	// Initialize Redis
	if err := redis.InitRedis(*cfg); err != nil {
		logger.Error(i18n.Translate("redis.init.failed", "", nil), "error", err)
		panic(err)
	}

	// Initialize HTTP client
	if err := thirdPlatform.InitHTTPClient(); err != nil {
		logger.Error(i18n.Translate("httpclient.init.failed", "", nil), "error", err)
		panic(err)
	}

	// Initialize Asynq server
	if err := asynq.InitServer(*cfg); err != nil {
		logger.Error(i18n.Translate("asynq.server.init.failed", "", nil), "error", err)
		panic(err)
	}

	// Register task handlers
	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.TypeRunReviewTask, tasks.HandleRunReviewTask)

	// Start worker
	logger.Info(i18n.Translate("worker.process.start", "", nil), "pid", os.Getpid())

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := asynq.Start(mux); err != nil {
			logger.Error(i18n.Translate("asynq.server.start.failed", "", nil), "error", err)
			quit <- syscall.SIGTERM
		}
	}()

	<-quit
	logger.Info(i18n.Translate("worker.process.stop", "", nil))

	// Resource cleanup
	asynq.Shutdown()

	if err := redis.Close(); err != nil {
		logger.Error(i18n.Translate("redis.client.close.failed", "", nil), "error", err)
	}

	if err := db.CloseDB(); err != nil {
		logger.Error(i18n.Translate("db.connection.close.failed", "", nil), "error", err)
	}

}
