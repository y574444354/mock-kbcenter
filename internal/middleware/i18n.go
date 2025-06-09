package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zgsm/go-webserver/i18n"
)

// I18n internationalization middleware
func I18n() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get language information from request
		locale := getLocaleFromRequest(c)

		// Set language information to context
		c.Set("locale", locale)

		// Add translation function to context
		c.Set("translate", func(messageID string, templateData map[string]interface{}) string {
			return i18n.Translate(messageID, locale, templateData)
		})

		c.Next()
	}
}

// getLocaleFromRequest gets language information from request
func getLocaleFromRequest(c *gin.Context) string {
	// First try to get from query parameters
	locale := c.Query("locale")
	if locale != "" {
		return locale
	}

	// Then try to get from Cookie
	localeCookie, err := c.Cookie("locale")
	if err == nil && localeCookie != "" {
		return localeCookie
	}

	// Then try to get from Accept-Language header
	acceptLanguage := c.GetHeader("Accept-Language")
	if acceptLanguage != "" {
		// Parse Accept-Language header
		// Format like: zh-CN,zh;q=0.9,en;q=0.8
		langs := strings.Split(acceptLanguage, ",")
		if len(langs) > 0 {
			// Get the first language
			lang := strings.TrimSpace(langs[0])
			// If has weight, remove the weight part
			if idx := strings.Index(lang, ";"); idx != -1 {
				lang = lang[:idx]
			}
			return lang
		}
	}

	// Finally use default locale from config
	return i18n.GetDefaultLocale()
}

// SetLocale handler for setting language
func SetLocale(c *gin.Context) {
	locale := c.Query("locale")
	if locale == "" {
		c.JSON(400, gin.H{
			"code":    400,
			"message": i18n.Translate("i18n.locale.missing", "", nil),
		})
		return
	}

	// Set Cookie
	c.SetCookie("locale", locale, 3600*24*30, "/", "", false, true)

	c.JSON(200, gin.H{
		"code":    200,
		"message": i18n.Translate("i18n.locale.set.success", locale, nil),
		"locale":  locale,
	})
}
