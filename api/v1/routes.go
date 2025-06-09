package v1

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes register API routes
func RegisterRoutes(router *gin.RouterGroup) {
	reviewTaskHandler := NewReviewTaskHandler()
	reviewTaskGroup := router.Group("/review_tasks")
	{
		reviewTaskGroup.POST("/", reviewTaskHandler.Create)
		reviewTaskGroup.GET("/:review_task_id/issues/increment", reviewTaskHandler.IssueIncrement)
	}
}
