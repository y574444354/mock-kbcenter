package v1

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册API路由
func RegisterRoutes(router *gin.RouterGroup, workDir string) {
	kbcenterHandler := NewKBCenterMockHandler(workDir)
	kbcenterHandler.RegisterRoutes(router)
}
