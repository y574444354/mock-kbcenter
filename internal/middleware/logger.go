package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zgsm/go-webserver/i18n"
	"github.com/zgsm/go-webserver/pkg/logger"
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
		latency := fmt.Sprintf("%.2fms", end.Sub(start).Seconds()*1000)

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
			"ip", clientIP,
			"latency", latency,
			"error", errorMessage,
			"user-agent", c.Request.UserAgent(),
		)
	}
}
