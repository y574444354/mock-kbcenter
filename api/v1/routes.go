package v1

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册API路由
func RegisterRoutes(router *gin.RouterGroup) {
	// 用户相关路由
	userHandler := NewUserHandler()
	userGroup := router.Group("/users")
	{
		userGroup.POST("/register", userHandler.Register)
		userGroup.POST("/login", userHandler.Login)
		userGroup.GET("/info", userHandler.GetUserInfo)
		userGroup.PUT("/info", userHandler.UpdateUserInfo)
		userGroup.PUT("/password", userHandler.ChangePassword)
		userGroup.GET("", userHandler.ListUsers)
		userGroup.DELETE("/:id", userHandler.DeleteUser)
	}

	// 可以添加其他API路由
	// 例如：商品、订单、支付等
	reviewTaskHandler := NewReviewTaskHandler()
	reviewTaskGroup := router.Group("/review-tasks")
	{
		reviewTaskGroup.POST("/", reviewTaskHandler.Create)
	}
}
