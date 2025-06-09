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

	web "github.com/zgsm/go-webserver/cmd/web"
	worker "github.com/zgsm/go-webserver/cmd/worker"
	"github.com/zgsm/go-webserver/config"
	"github.com/zgsm/go-webserver/i18n"
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
