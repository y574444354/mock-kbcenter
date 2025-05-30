package main

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/zgsm/mock-kbcenter/config"
	"github.com/zgsm/mock-kbcenter/internal/model"
	"github.com/zgsm/mock-kbcenter/pkg/db"
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
			os.Exit(1)
		}

		// 注册所有需要迁移的模型
		log.Println("Migrating database...")
		if err := db.AutoMigrate(
			&model.ReviewTask{},
			// 在这里添加其他模型
		); err != nil {
			log.Fatalf("Failed to migrate database: %v", err)
			os.Exit(1)
		}
		log.Println("Database migration completed.")
	},
}
