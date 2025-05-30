package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/zgsm/go-webserver/config"
	"github.com/zgsm/go-webserver/i18n"
	"github.com/zgsm/go-webserver/pkg/logger"
	"github.com/zgsm/go-webserver/pkg/redis"
)

func main() {
	// 初始化配置
	if err := config.LoadConfigWithDefault(); err != nil {
		log.Fatalln("config.load.failed: %w", err)
		os.Exit(1)
	}

	// 初始化i18n
	if err := i18n.InitI18n(*config.GetConfig()); err != nil {
		fmt.Println("i18n.init.failed: %w", err)
		os.Exit(1)
	}

	// 初始化redis
	if err := redis.InitRedis(*config.GetConfig()); err != nil {
		fmt.Printf(i18n.Translate("redis.connect.failed", "", nil)+": %v\n", err)
		os.Exit(1)
	}
	defer redis.Close()

	// 解析命令行参数
	flag.Parse()

	// 执行清除缓存操作
	if err := redis.FlushDB(); err != nil {
		logger.Error(i18n.Translate("redis.flushdb.failed", "", nil), "error", err)
		os.Exit(1)
	}

	fmt.Println(i18n.Translate("redis.flushdb.success", "", nil))
}
