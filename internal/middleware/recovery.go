package middleware

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zgsm/go-webserver/i18n"
	"github.com/zgsm/go-webserver/pkg/logger"
)

// Recovery middleware for recovering from panics
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check if connection is broken
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
							strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				// Request information
				httpRequest, err := httputil.DumpRequest(c.Request, false)
				if err != nil {
					httpRequest = []byte(fmt.Sprintf("failed to dump request: %v", err))
				}
				// Stack trace
				stack := string(debug.Stack())

				if brokenPipe {
					logger.Error(i18n.Translate("middleware.recovery.broken_pipe", "", nil),
						"error", err,
						"request", string(httpRequest),
					)
					// If connection is broken, we can't write status to client
					c.Abort()
					return
				}

				// Log error
				logger.Error(i18n.Translate("middleware.recovery.recovered", "", nil),
					"error", err,
					"request", string(httpRequest),
					"stack", stack,
				)

				// Return 500 error
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code":    http.StatusInternalServerError,
					"message": fmt.Sprintf(i18n.Translate("middleware.recovery.internal_error", "", nil), err),
				})
			}
		}()
		c.Next()
	}
}
