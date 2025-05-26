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
	// 添加子命令
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(initCmd)
}

func main() {
	// 加载配置
	if err := config.LoadConfigWithDefault(); err != nil {
		log.Fatalln("config.load.failed: %w", err)
	}
	// 初始化配置
	cfg := config.GetConfig()

	if err := i18n.InitI18n(*cfg); err != nil {
		fmt.Println("i18n.init.failed: %w", err)
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
