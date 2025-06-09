package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Cors middleware for handling CORS requests
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		// origin := c.Request.Header.Get("Origin")

		// Allowed domains
		allowOrigin := "*"
		// If you need to dynamically set allowed domains based on request Origin, use the following code
		// allowOrigins := []string{"http://localhost:8080", "https://example.com"}
		// for _, o := range allowOrigins {
		//     if o == origin {
		//         allowOrigin = o
		//         break
		//     }
		// }

		// Set CORS-related HTTP headers
		c.Header("Access-Control-Allow-Origin", allowOrigin)
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		// Allow all OPTIONS methods
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		// Process the request
		c.Next()
	}
}

// CorsWithConfig middleware for handling CORS requests with configuration
func CorsWithConfig(allowOrigins []string, allowMethods []string, allowHeaders []string, exposeHeaders []string, allowCredentials bool, maxAge int) gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")

		// Allowed domains
		allowOrigin := "*"
		// If allowed domains are specified, check if request Origin is in the allowed list
		if len(allowOrigins) > 0 {
			allowOrigin = ""
			for _, o := range allowOrigins {
				if o == origin {
					allowOrigin = o
					break
				}
			}
		}

		// Set CORS-related HTTP headers
		if allowOrigin != "" {
			c.Header("Access-Control-Allow-Origin", allowOrigin)
		}

		// Allowed methods
		if len(allowMethods) > 0 {
			c.Header("Access-Control-Allow-Methods", strings.Join(allowMethods, ", "))
		} else {
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		}

		// Allowed headers
		if len(allowHeaders) > 0 {
			c.Header("Access-Control-Allow-Headers", strings.Join(allowHeaders, ", "))
		} else {
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		}

		// Exposed headers
		if len(exposeHeaders) > 0 {
			c.Header("Access-Control-Expose-Headers", strings.Join(exposeHeaders, ", "))
		} else {
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		}

		// Whether to allow credentials
		if allowCredentials {
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		// Preflight request validity period
		if maxAge > 0 {
			c.Header("Access-Control-Max-Age", strconv.Itoa(maxAge))
		}

		// Allow all OPTIONS methods
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		// Process the request
		c.Next()
	}
}
