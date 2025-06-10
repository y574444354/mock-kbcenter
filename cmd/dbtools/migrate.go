package main

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/zgsm/mock-kbcenter/config"

	// "github.com/zgsm/mock-kbcenter/internal/model"
	"github.com/zgsm/mock-kbcenter/pkg/db"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate database",
	Long:  "Execute database migrations, create or update database tables",
	Run: func(cmd *cobra.Command, args []string) {
		// Load configuration
		if err := config.LoadConfigWithDefault(); err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}

		// Initialize database connection
		if err := db.InitDB(config.GetConfig().Database); err != nil {
			log.Fatalf("Failed to initialize database: %v", err)
		}

		// Register all models that need migration
		log.Println("Migrating database...")
		if err := db.AutoMigrate(
		// Add other models here
		); err != nil {
			log.Fatalf("Failed to migrate database: %v", err)
		}
		log.Println("Database migration completed.")
	},
}
