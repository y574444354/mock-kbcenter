package main

import (
	"fmt"
	"os"

	web "github.com/zgsm/review-manager/cmd/web"
	worker "github.com/zgsm/review-manager/cmd/worker"
	"github.com/zgsm/review-manager/i18n"
	"github.com/zgsm/review-manager/pkg/logger"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println(i18n.Translate("service.usage", "", nil))
		fmt.Println(i18n.Translate("service.web", "", nil))
		fmt.Println(i18n.Translate("service.worker", "", nil))
		os.Exit(1)
	}

	switch os.Args[1] {
	case "web":
		web.Run()
	case "worker":
		worker.Run()
	default:
		logger.Error(i18n.Translate("service.unknown", "", nil), "service", os.Args[1])
		os.Exit(1)
	}
}
