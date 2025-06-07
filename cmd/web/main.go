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
	v1 "github.com/zgsm/go-webserver/api/v1"
	"github.com/zgsm/go-webserver/config"
	"github.com/zgsm/go-webserver/i18n"
	"github.com/zgsm/go-webserver/internal/middleware"
	"github.com/zgsm/go-webserver/pkg/asynq"
	"github.com/zgsm/go-webserver/pkg/db"
	"github.com/zgsm/go-webserver/pkg/logger"
	"github.com/zgsm/go-webserver/pkg/redis"
	"github.com/zgsm/go-webserver/pkg/thirdPlatform"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/zgsm/go-webserver/swagger"
)

func Run(cfg *config.Config) {
	// 初始化日志
	if err := logger.InitLogger(cfg.Log); err != nil {
		log.Fatalln(i18n.Translate("logger.init.failed", "", nil), "error", err)
	}
	defer logger.Sync()

	// 设置Gin模式
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else if cfg.Server.Mode == "test" {
		gin.SetMode(gin.TestMode)
	}

	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				logger.Error("web server panic", "error", err)
			}
			logger.Sync()
			os.Exit(1)
		}
	}()

	// 初始化数据库
	if err := db.InitDB(cfg.Database); err != nil {
		logger.Error(i18n.Translate("db.connection.init", "", nil), "error", err)
		panic(err)
	}
	defer func() {
		if err := db.CloseDB(); err != nil {
			logger.Error(i18n.Translate("db.connection.close.failed", "", nil), "error", err)
		}
	}()

	// 初始化Redis
	if err := redis.InitRedis(*cfg); err != nil {
		logger.Error(i18n.Translate("redis.connect.failed", "", nil), "error", err)
		panic(err)
	}
	defer func() {
		if err := redis.Close(); err != nil {
			logger.Error(i18n.Translate("redis.client.close.failed", "", nil), "error", err)
		}
	}()

	// 初始化HTTP客户端
	if err := thirdPlatform.InitHTTPClient(); err != nil {
		logger.Error(i18n.Translate("httpclient.init.failed", "", nil), "error", err)
		panic(err)
	}

	// 初始化Asynq客户端（如果启用）
	if cfg.Asynq.Enabled {
		if err := asynq.InitClient(*cfg); err != nil {
			logger.Error(i18n.Translate("asynq.client.init.failed", "", nil), "error", err)
			panic(err)
		}
		defer asynq.Close()
	}

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
	v1.RegisterRoutes(apiV1)

	// 启动服务器
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: r,
	}

	// 优雅关闭
	go func() {
		logger.Info(i18n.Translate("server.start.success", "", nil), "port", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(i18n.Translate("server.start.failed", "", nil), "error", err)
		}
	}()

	// 等待中断信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	logger.Info(i18n.Translate("server.shutdown.starting", "", nil))

	// 设置关闭超时
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal(i18n.Translate("server.shutdown.forced", "", nil), "error", err)
	}

	logger.Info(i18n.Translate("server.shutdown.success", "", nil))
}
