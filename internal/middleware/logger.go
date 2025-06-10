package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zgsm/mock-kbcenter/i18n"
	"github.com/zgsm/mock-kbcenter/pkg/logger"
)

// Logger middleware for logging
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start time
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// End time
		end := time.Now()
		latency := end.Sub(start)

		// Request method
		method := c.Request.Method
		// Status code
		statusCode := c.Writer.Status()
		// Client IP
		clientIP := c.ClientIP()
		// Error message
		errorMessage := c.Errors.String()

		// If has query parameters, add to path
		if raw != "" {
			path = path + "?" + raw
		}

		// Log the request
		logger.Info(i18n.Translate("log.http.request", "", map[string]interface{}{
			"method": method,
			"path":   path,
		}),
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
