package main

import (
	"fmt"
	"log"
	"os"

	web "github.com/zgsm/mock-kbcenter/cmd/web"
	// worker "github.com/zgsm/mock-kbcenter/cmd/worker"
	"github.com/zgsm/mock-kbcenter/config"
	"github.com/zgsm/mock-kbcenter/i18n"
)

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

	workDir := ""
	if len(os.Args) > 1 {
		workDir = os.Args[1]
	}
	if workDir == "" {
		var err error
		workDir, err = os.Getwd()
		if err != nil {
			log.Fatalln(i18n.Translate("kbcenter.getwd_failed", "", map[string]interface{}{"error": err.Error()}))
		}
	}
	fmt.Println(i18n.Translate("kbcenter.workdir", "", map[string]interface{}{"workdir": workDir}))
	web.Run(cfg, workDir)
}
