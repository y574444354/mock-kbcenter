package main

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/zgsm/go-webserver/config"
	"github.com/zgsm/go-webserver/i18n"
	"github.com/zgsm/go-webserver/internal/model"
	"github.com/zgsm/go-webserver/pkg/db"
	"github.com/zgsm/go-webserver/pkg/logger"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: i18n.Translate("db.migrate.start", "", nil),
	Long:  i18n.Translate("dbtools.migrate.description", "", nil),
	Run: func(cmd *cobra.Command, args []string) {
		// 加载配置
		if err := config.LoadConfigWithDefault(); err != nil {
			log.Fatalf(i18n.Translate("config.load.failed", "", nil)+": %v", err)
		}

		// 初始化日志
		if err := logger.InitLogger(config.GetConfig().Log); err != nil {
			log.Fatalf(i18n.Translate("logger.init.failed", "", nil)+": %v", err)
		}
		defer logger.Sync()

		// 初始化数据库连接
		if err := db.InitDB(*config.GetConfig()); err != nil {
			logger.Error(i18n.Translate("db.connection.init", "", nil), "error", err)
			os.Exit(1)
		}

		// 执行数据库迁移
		logger.Info(i18n.Translate("db.migrate.start", "", nil))

		// 注册所有需要迁移的模型
		logger.Info(i18n.Translate("db.model.register", "", nil))
		if err := db.AutoMigrate(
			&model.User{},
			&model.UserProfile{},
			// 在这里添加其他模型
		); err != nil {
			logger.Error(i18n.Translate("db.migrate.failed", "", nil), "error", err)
			os.Exit(1)
		}
		logger.Info(i18n.Translate("db.migrate.success", "", nil))
	},
}
