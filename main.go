package main

// @title Go WebServer API
// @version 1.0
// @description This is the API documentation for Go WebServer
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
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

	proxy "github.com/zgsm/mock-kbcenter/cmd/proxy"
	web "github.com/zgsm/mock-kbcenter/cmd/web"

	// worker "github.com/zgsm/mock-kbcenter/cmd/worker"
	"github.com/zgsm/mock-kbcenter/config"
	"github.com/zgsm/mock-kbcenter/i18n"
)

func main() {
	// Load configuration
	if err := config.LoadConfigWithDefault(); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	// Initialize configuration
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

	if len(os.Args) > 2 && os.Args[2] == "proxy" {
		proxy.Run(cfg, workDir)
	} else {
		web.Run(cfg, workDir)
	}
}
