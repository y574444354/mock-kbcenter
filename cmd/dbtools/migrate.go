package main

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/zgsm/go-webserver/config"
	"github.com/zgsm/go-webserver/internal/model"
	"github.com/zgsm/go-webserver/pkg/db"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate database",
	Long:  "Execute database migrations, create or update database tables",
	Run: func(cmd *cobra.Command, args []string) {
		// 加载配置
		if err := config.LoadConfigWithDefault(); err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}

		// 初始化数据库连接
		if err := db.InitDB(config.GetConfig().Database); err != nil {
			log.Fatalf("Failed to initialize database: %v", err)
		}

		// 注册所有需要迁移的模型
		log.Println("Migrating database...")
		if err := db.AutoMigrate(
			&model.ReviewTask{},
			// 在这里添加其他模型
		); err != nil {
			log.Fatalf("Failed to migrate database: %v", err)
		}
		log.Println("Database migration completed.")
	},
}
