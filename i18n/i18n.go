package i18n

import (
	"embed"
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

//go:embed locales/*.yaml
var localesFS embed.FS

var (
	bundle       *i18n.Bundle
	defaultLocal string
)

// InitI18n 初始化国际化支持
func InitI18n(cfg config.Config) error {
	// 记录传入的配置值
	log.Printf("i18n.init.config: default_locale=%q, bundle_path=%q",
		cfg.I18n.DefaultLocale, cfg.I18n.BundlePath)

	// 验证语言标签
	if cfg.I18n.DefaultLocale == "" {
		return fmt.Errorf("i18n.init.failed: default locale is empty")
	}

	// 尝试解析语言标签
	lang, err := language.Parse(cfg.I18n.DefaultLocale)
	if err != nil {
		return fmt.Errorf("i18n.init.failed: invalid locale tag %q: %w", cfg.I18n.DefaultLocale, err)
	}

	// 创建bundle
	bundle = i18n.NewBundle(lang)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

	defaultLocal = cfg.I18n.DefaultLocale

	// 加载语言文件
	if err := loadMessageFiles(cfg.I18n.BundlePath); err != nil {
		return fmt.Errorf("i18n.load.failed: %w", err)
	}

	log.Println(Translate("i18n.init.success", defaultLocal, map[string]interface{}{"locale": defaultLocal}))
	return nil
}

// loadMessageFiles 加载指定目录下的所有语言文件
func loadMessageFiles(dir string) error {
	// 首先尝试从嵌入文件系统加载
	files, err := localesFS.ReadDir("locales")
	if err == nil {
		for _, file := range files {
			if file.IsDir() {
				continue
			}

			// 只处理yaml文件
			if !strings.HasSuffix(file.Name(), ".yaml") && !strings.HasSuffix(file.Name(), ".yml") {
				continue
			}

			// 加载语言文件
			path := filepath.Join("locales", file.Name())
			data, err := localesFS.ReadFile(path)
			if err != nil {
				return fmt.Errorf("i18n.load.embed.file.failed: %s: %w", path, err)
			}

			// 使用内存中的内容加载
			_, err = bundle.ParseMessageFileBytes(data, path)
			if err != nil {
				return fmt.Errorf("i18n.parse.embed.file.failed: %s: %w", path, err)
			}

			log.Printf("i18n.load.embed.file.success: file=%s", path)
		}
		return nil
	}

	// 如果嵌入文件加载失败，回退到文件系统
	files, err = os.ReadDir(dir)
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
	if bundle == nil {
		log.Printf("i18n.translate.failed: bundle not initialized, messageID=%s", messageID)
		return messageID
	}

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
