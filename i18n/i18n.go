package i18n

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/zgsm/go-webserver/config"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

var (
	bundle       *i18n.Bundle
	defaultLocal string
)

// InitI18n 初始化国际化支持
func InitI18n(cfg config.Config) error {
	// 创建bundle
	bundle = i18n.NewBundle(language.MustParse(cfg.I18n.DefaultLocale))
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

	defaultLocal = cfg.I18n.DefaultLocale

	// 加载语言文件
	err := loadMessageFiles(cfg.I18n.BundlePath)
	if err != nil {
		return fmt.Errorf("i18n.load.failed: %w", err)
	}

	log.Printf("i18n.init.success: default_locale=%s", defaultLocal)
	return nil
}

// loadMessageFiles 加载指定目录下的所有语言文件
func loadMessageFiles(dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// 只处理yaml文件
		if !strings.HasSuffix(file.Name(), ".yaml") && !strings.HasSuffix(file.Name(), ".yml") {
			continue
		}

		// 加载语言文件
		path := filepath.Join(dir, file.Name())
		_, err := bundle.LoadMessageFile(path)
		if err != nil {
			return fmt.Errorf("i18n.load.file.failed: %s: %w", path, err)
		}

		log.Printf("i18n.load.file.success: file=%s", path)
	}

	return nil
}

// Translate 翻译消息
func Translate(messageID string, locale string, templateData map[string]interface{}) string {
	if locale == "" {
		locale = defaultLocal
	}

	localizer := i18n.NewLocalizer(bundle, locale)
	message, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: templateData,
	})

	if err != nil {
		log.Printf("i18n.translate.failed: messageID=%s, locale=%s, error=%v", messageID, locale, err)
		return messageID
	}

	return message
}

// GetDefaultLocale 获取默认语言
func GetDefaultLocale() string {
	return defaultLocal
}

// GetBundle 获取国际化bundle
func GetBundle() *i18n.Bundle {
	return bundle
}
