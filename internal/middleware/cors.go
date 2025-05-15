package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Cors 处理跨域请求的中间件
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		// origin := c.Request.Header.Get("Origin")

		// 允许的域名
		allowOrigin := "*"
		// 如果需要根据请求的Origin动态设置允许的域名，可以使用以下代码
		// allowOrigins := []string{"http://localhost:8080", "https://example.com"}
		// for _, o := range allowOrigins {
		//     if o == origin {
		//         allowOrigin = o
		//         break
		//     }
		// }

		// 设置CORS相关的HTTP头
		c.Header("Access-Control-Allow-Origin", allowOrigin)
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		// 放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		// 处理请求
		c.Next()
	}
}

// CorsWithConfig 使用配置处理跨域请求的中间件
func CorsWithConfig(allowOrigins []string, allowMethods []string, allowHeaders []string, exposeHeaders []string, allowCredentials bool, maxAge int) gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")

		// 允许的域名
		allowOrigin := "*"
		// 如果指定了允许的域名，则检查请求的Origin是否在允许列表中
		if len(allowOrigins) > 0 {
			allowOrigin = ""
			for _, o := range allowOrigins {
				if o == origin {
					allowOrigin = o
					break
				}
			}
		}

		// 设置CORS相关的HTTP头
		if allowOrigin != "" {
			c.Header("Access-Control-Allow-Origin", allowOrigin)
		}

		// 允许的方法
		if len(allowMethods) > 0 {
			c.Header("Access-Control-Allow-Methods", strings.Join(allowMethods, ", "))
		} else {
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		}

		// 允许的头
		if len(allowHeaders) > 0 {
			c.Header("Access-Control-Allow-Headers", strings.Join(allowHeaders, ", "))
		} else {
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		}

		// 暴露的头
		if len(exposeHeaders) > 0 {
			c.Header("Access-Control-Expose-Headers", strings.Join(exposeHeaders, ", "))
		} else {
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		}

		// 是否允许凭证
		if allowCredentials {
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		// 预检请求的有效期
		if maxAge > 0 {
			c.Header("Access-Control-Max-Age", strconv.Itoa(maxAge))
		}

		// 放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		// 处理请求
		c.Next()
	}
}
