package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zgsm/mock-kbcenter/i18n"
)

// I18n 国际化中间件
func I18n() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求中获取语言信息
		locale := getLocaleFromRequest(c)

		// 将语言信息设置到上下文中
		c.Set("locale", locale)

		// 添加翻译函数到上下文
		c.Set("translate", func(messageID string, templateData map[string]interface{}) string {
			return i18n.Translate(messageID, locale, templateData)
		})

		c.Next()
	}
}

// getLocaleFromRequest 从请求中获取语言信息
func getLocaleFromRequest(c *gin.Context) string {
	// 优先从查询参数中获取
	locale := c.Query("locale")
	if locale != "" {
		return locale
	}

	// 从Cookie中获取
	localeCookie, err := c.Cookie("locale")
	if err == nil && localeCookie != "" {
		return localeCookie
	}

	// 从Accept-Language头中获取
	acceptLanguage := c.GetHeader("Accept-Language")
	if acceptLanguage != "" {
		// 解析Accept-Language头
		// 格式如：zh-CN,zh;q=0.9,en;q=0.8
		langs := strings.Split(acceptLanguage, ",")
		if len(langs) > 0 {
			// 取第一个语言
			lang := strings.TrimSpace(langs[0])
			// 如果有权重，去掉权重部分
			if idx := strings.Index(lang, ";"); idx != -1 {
				lang = lang[:idx]
			}
			return lang
		}
	}

	// 默认使用配置中的默认语言
	return i18n.GetDefaultLocale()
}

// SetLocale 设置语言的处理器
func SetLocale(c *gin.Context) {
	locale := c.Query("locale")
	if locale == "" {
		c.JSON(400, gin.H{
			"code":    400,
			"message": i18n.Translate("i18n.locale.missing", "", nil),
		})
		return
	}

	// 设置Cookie
	c.SetCookie("locale", locale, 3600*24*30, "/", "", false, true)

	c.JSON(200, gin.H{
		"code":    200,
		"message": i18n.Translate("i18n.locale.set.success", locale, nil),
		"locale":  locale,
	})
}
