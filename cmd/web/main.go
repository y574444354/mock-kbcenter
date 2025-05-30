package web

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	v1 "github.com/zgsm/mock-kbcenter/api/v1"
	"github.com/zgsm/mock-kbcenter/config"
	"github.com/zgsm/mock-kbcenter/i18n"
	"github.com/zgsm/mock-kbcenter/internal/middleware"
	"github.com/zgsm/mock-kbcenter/pkg/logger"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/zgsm/mock-kbcenter/swagger"
)

func Run(cfg *config.Config, workDir string) {
	locale := i18n.GetDefaultLocale()

	// 初始化日志
	if err := logger.InitLogger(cfg.Log); err != nil {
		log.Fatalln(i18n.Translate("logger.init.failed", locale, nil), "error", err)
	}
	defer logger.Sync()

	// 设置Gin模式
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else if cfg.Server.Mode == "test" {
		gin.SetMode(gin.TestMode)
	}

	// // 初始化数据库
	// if err := db.InitDB(cfg.Database); err != nil {
	// 	logger.Error(i18n.Translate("db.connection.init", locale, nil), "error", err)
	// 	os.Exit(1)
	// }

	// // 初始化Redis
	// if err := redis.InitRedis(*cfg); err != nil {
	// 	logger.Error(i18n.Translate("redis.connect.failed", locale, nil), "error", err)
	// 	os.Exit(1)
	// }

	// // 初始化HTTP客户端
	// if err := thirdPlatform.InitHTTPClient(); err != nil {
	// 	logger.Error(i18n.Translate("httpclient.init.failed", locale, nil), "error", err)
	// 	os.Exit(1)
	// }

	// // 初始化Asynq客户端
	// if err := asynq.InitClient(*cfg); err != nil {
	// 	logger.Error(i18n.Translate("asynq.client.init.failed", locale, nil), "error", err)
	// 	os.Exit(1)
	// }
	// defer asynq.Close()

	// 创建Gin引擎
	r := gin.New()

	// 注册中间件
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.Cors())
	r.Use(middleware.I18n())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// Swagger文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 注册API路由
	apiV1 := r.Group("/api/v1")
	v1.RegisterRoutes(apiV1, workDir)

	// 启动服务器
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: r,
	}

	// 优雅关闭
	go func() {
		locale := i18n.GetDefaultLocale()
		logger.Info(i18n.Translate("server.start.success", locale, nil), "port", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(i18n.Translate("server.start.failed", locale, nil), "error", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info(i18n.Translate("server.shutdown.starting", locale, nil))

	// 设置关闭超时
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal(i18n.Translate("server.shutdown.forced", locale, nil), "error", err)
	}

	logger.Info(i18n.Translate("server.shutdown.success", locale, nil))
}
