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
	"github.com/zgsm/mock-kbcenter/i18n"
	"github.com/zgsm/mock-kbcenter/pkg/logger"
)

// Recovery 从panic中恢复的中间件
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 检查连接是否已断开
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
							strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				// 请求信息
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				// 堆栈信息
				stack := string(debug.Stack())

				if brokenPipe {
					logger.Error(i18n.Translate("middleware.recovery.broken_pipe", "", nil),
						"error", err,
						"request", string(httpRequest),
					)
					// 如果连接已断开，我们无法向客户端写入状态
					c.Abort()
					return
				}

				// 记录错误日志
				logger.Error(i18n.Translate("middleware.recovery.recovered", "", nil),
					"error", err,
					"request", string(httpRequest),
					"stack", stack,
				)

				// 返回500错误
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code":    http.StatusInternalServerError,
					"message": fmt.Sprintf(i18n.Translate("middleware.recovery.internal_error", "", nil), err),
				})
			}
		}()
		c.Next()
	}
}
