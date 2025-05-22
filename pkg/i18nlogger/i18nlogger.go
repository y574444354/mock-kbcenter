package i18nlogger

import (
	"github.com/gin-gonic/gin"
	"github.com/zgsm/review-manager/i18n"
	"github.com/zgsm/review-manager/pkg/logger"
)

// Debug 记录国际化的调试级别日志
func Debug(messageID string, locale string, templateData map[string]interface{}, keysAndValues ...interface{}) {
	msg := i18n.Translate(messageID, locale, templateData)
	logger.Debug(msg, keysAndValues...)
}

// Info 记录国际化的信息级别日志
func Info(messageID string, locale string, templateData map[string]interface{}, keysAndValues ...interface{}) {
	msg := i18n.Translate(messageID, locale, templateData)
	logger.Info(msg, keysAndValues...)
}

// Warn 记录国际化的警告级别日志
func Warn(messageID string, locale string, templateData map[string]interface{}, keysAndValues ...interface{}) {
	msg := i18n.Translate(messageID, locale, templateData)
	logger.Warn(msg, keysAndValues...)
}

// Error 记录国际化的错误级别日志
func Error(messageID string, locale string, templateData map[string]interface{}, keysAndValues ...interface{}) {
	msg := i18n.Translate(messageID, locale, templateData)
	logger.Error(msg, keysAndValues...)
}

// DPanic 记录国际化的开发环境恐慌级别日志
func DPanic(messageID string, locale string, templateData map[string]interface{}, keysAndValues ...interface{}) {
	msg := i18n.Translate(messageID, locale, templateData)
	logger.DPanic(msg, keysAndValues...)
}

// Panic 记录国际化的恐慌级别日志
func Panic(messageID string, locale string, templateData map[string]interface{}, keysAndValues ...interface{}) {
	msg := i18n.Translate(messageID, locale, templateData)
	logger.Panic(msg, keysAndValues...)
}

// Fatal 记录国际化的致命级别日志
func Fatal(messageID string, locale string, templateData map[string]interface{}, keysAndValues ...interface{}) {
	msg := i18n.Translate(messageID, locale, templateData)
	logger.Fatal(msg, keysAndValues...)
}

// GetLocaleFromContext 从Gin上下文中获取语言设置
func GetLocaleFromContext(c interface{}) string {
	if ctx, ok := c.(*gin.Context); ok {
		if locale, exists := ctx.Get("locale"); exists {
			if localeStr, ok := locale.(string); ok {
				return localeStr
			}
		}
	}
	return i18n.GetDefaultLocale()
}
