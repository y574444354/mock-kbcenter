package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/zgsm/go-webserver/config"
	"github.com/zgsm/go-webserver/pkg/headerpropagation"
)

// HeaderPropagator middleware propagates specified headers to context
//
// Example usage in handler:
//
//	// Get single header value
//	value := headerpropagation.GetHeaderValue(c.Request.Context(), "X-Request-ID")
//
//	// Get all propagated headers
//	headers := headerpropagation.GetAllPropagatedHeaders(c.Request.Context())
func HeaderPropagator() gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg := config.GetConfig()
		for _, header := range cfg.HeaderPropagation.Headers {
			if value := c.GetHeader(header); value != "" {
				ctx := context.WithValue(c.Request.Context(), headerpropagation.ContextKeyPrefix+headerpropagation.ContextKey(header), value)
				c.Request = c.Request.WithContext(ctx)
			}
		}
		c.Next()
	}
}
