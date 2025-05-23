package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zgsm/review-manager/pkg/i18nlogger"
)

// Logger 日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// 处理请求
		c.Next()

		// 结束时间
		end := time.Now()
		latency := end.Sub(start)

		// 请求方法
		method := c.Request.Method
		// 状态码
		statusCode := c.Writer.Status()
		// 客户端IP
		clientIP := c.ClientIP()
		// 错误信息
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		// 如果有查询参数，添加到路径
		if raw != "" {
			path = path + "?" + raw
		}

		// 获取语言设置
		locale := i18nlogger.GetLocaleFromContext(c)

		// 记录日志
		i18nlogger.Info("log.http.request", locale, map[string]interface{}{
			"method": method,
			"path":   path,
		},
			"status", statusCode,
			"method", method,
			"path", path,
			"ip", clientIP,
			"latency", latency,
			"error", errorMessage,
			"user-agent", c.Request.UserAgent(),
		)
	}
}
