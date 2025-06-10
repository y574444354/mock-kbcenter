package v1

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes register API routes
func RegisterRoutes(router *gin.RouterGroup, workDir string) {
	kbcenterHandler := NewKBCenterMockHandler(workDir)
	kbcenterHandler.RegisterRoutes(router)
}
