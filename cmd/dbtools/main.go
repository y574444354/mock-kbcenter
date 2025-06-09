package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/zgsm/go-webserver/config"
	"github.com/zgsm/go-webserver/i18n"
)

var (
	rootCmd = &cobra.Command{
		Use:   "dbtools",
		Short: "Database tools for executing migrations and initialization operations",
		Long:  "Database tools for executing migrations and initialization operations",
	}
)

func init() {
	// Add subcommands
	if migrateCmd == nil {
		log.Fatalf("migrateCmd is nil")
	}
	rootCmd.AddCommand(migrateCmd)

	if initCmd == nil {
		log.Fatalf("initCmd is nil")
	}
	rootCmd.AddCommand(initCmd)
}

func main() {
	// Load configuration
	if err := config.LoadConfigWithDefault(); err != nil {
		log.Fatalf("config.load.failed: %v", err)
	}
	// Initialize configuration
	cfg := config.GetConfig()

	if err := i18n.InitI18n(*cfg); err != nil {
		fmt.Printf("i18n.init.failed: %v\n", err)
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
