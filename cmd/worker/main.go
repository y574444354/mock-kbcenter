package worker

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/zgsm/go-webserver/config"
	"github.com/zgsm/go-webserver/i18n"
	"github.com/zgsm/go-webserver/pkg/asynq"
	"github.com/zgsm/go-webserver/pkg/logger"
	"github.com/zgsm/go-webserver/tasks"
)

func Run(cfg *config.Config) {
	// 初始化日志
	if err := logger.InitLogger(cfg.Asynq.Log); err != nil {
		log.Fatalln(i18n.Translate("asynq.server.init.failed", "", nil), "error", err)
		os.Exit(1)
	}
	defer logger.Sync()

	// 初始化Asynq服务器
	if err := asynq.InitServer(*cfg); err != nil {
		logger.Error(i18n.Translate("asynq.server.init.failed", "", nil), "error", err)
		os.Exit(1)
	}

	// 注册任务处理器
	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.TypeRunReviewTask, tasks.HandleRunReviewTask)

	// 启动worker
	logger.Info(i18n.Translate("worker.process.start", "", nil), "pid", os.Getpid())

	// 优雅退出
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
}
