package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zgsm/review-manager/i18n"
)

var (
	rootCmd = &cobra.Command{
		Use:   "dbtools",
		Short: i18n.Translate("dbtools.description", "", nil),
		Long:  i18n.Translate("dbtools.description", "", nil),
	}
)

func init() {

	// 添加子命令
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(initCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
