package main

// @title Go WebServer API
// @version 1.0
// @description 这是Go WebServer的API文档
// @termsOfService http://swagger.io/terms/

// @contact.name API支持
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

import (
	"fmt"
	"log"
	"os"

	web "github.com/zgsm/go-webserver/cmd/web"
	worker "github.com/zgsm/go-webserver/cmd/worker"
	"github.com/zgsm/go-webserver/config"
	"github.com/zgsm/go-webserver/i18n"
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

	if len(os.Args) < 2 {
		fmt.Println(i18n.Translate("service.usage", "", nil))
		fmt.Println(i18n.Translate("service.web", "", nil))
		fmt.Println(i18n.Translate("service.worker", "", nil))
		web.Run(cfg)
		return
	}

	switch os.Args[1] {
	case "web":
		web.Run(cfg)
	case "worker":
		worker.Run(cfg)
	default:
		log.Fatalln(i18n.Translate("service.unknown", "", nil), "service", os.Args[1])
	}
}
