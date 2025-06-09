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

// InitI18n initializes internationalization support
func InitI18n(cfg config.Config) error {
	// Log the incoming configuration values
	log.Printf("i18n.init.config: default_locale=%q, bundle_path=%q",
		cfg.I18n.DefaultLocale, cfg.I18n.BundlePath)

	// Validate language tag
	if cfg.I18n.DefaultLocale == "" {
		return fmt.Errorf("i18n.init.failed: default locale is empty")
	}

	// Try to parse language tag
	lang, err := language.Parse(cfg.I18n.DefaultLocale)
	if err != nil {
		return fmt.Errorf("i18n.init.failed: invalid locale tag %q: %w", cfg.I18n.DefaultLocale, err)
	}

	// Create bundle
	bundle = i18n.NewBundle(lang)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

	defaultLocal = cfg.I18n.DefaultLocale

	// Load language files
	if err := loadMessageFiles(cfg.I18n.BundlePath); err != nil {
		return fmt.Errorf("i18n.load.failed: %w", err)
	}

	log.Println(Translate("i18n.init.success", defaultLocal, map[string]interface{}{"locale": defaultLocal}))
	return nil
}

// loadMessageFiles loads all language files in the specified directory
func loadMessageFiles(dir string) error {
	// First try to load from embedded file system
	files, err := localesFS.ReadDir("locales")
	if err == nil {
		for _, file := range files {
			if file.IsDir() {
				continue
			}

			// Only process yaml files
			if !strings.HasSuffix(file.Name(), ".yaml") && !strings.HasSuffix(file.Name(), ".yml") {
				continue
			}

			// Load language file
			path := filepath.Join("locales", file.Name())
			data, err := localesFS.ReadFile(path)
			if err != nil {
				return fmt.Errorf("i18n.load.embed.file.failed: %s: %w", path, err)
			}

			// Load using in-memory content
			_, err = bundle.ParseMessageFileBytes(data, path)
			if err != nil {
				return fmt.Errorf("i18n.parse.embed.file.failed: %s: %w", path, err)
			}

			log.Printf("i18n.load.embed.file.success: file=%s", path)
		}
		return nil
	}

	// If embedded file loading fails, fall back to filesystem
	files, err = os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// Only process yaml files
		if !strings.HasSuffix(file.Name(), ".yaml") && !strings.HasSuffix(file.Name(), ".yml") {
			continue
		}

		// Load language file
		path := filepath.Join(dir, file.Name())
		_, err := bundle.LoadMessageFile(path)
		if err != nil {
			return fmt.Errorf("i18n.load.file.failed: %s: %w", path, err)
		}

		log.Printf("i18n.load.file.success: file=%s", path)
	}

	return nil
}

// Translate translates messages
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

// GetDefaultLocale get default locale
func GetDefaultLocale() string {
	return defaultLocal
}

// GetBundle gets the internationalization bundle
func GetBundle() *i18n.Bundle {
	return bundle
}
