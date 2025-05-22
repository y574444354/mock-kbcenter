package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/zgsm/review-manager/api"
	"github.com/zgsm/review-manager/internal/service"
)

type ReviewTaskHandler struct {
	reviewTaskService service.ReviewTaskService
}

func NewReviewTaskHandler() *ReviewTaskHandler {
	return &ReviewTaskHandler{
		reviewTaskService: service.NewReviewTaskService(),
	}
}

func (h *ReviewTaskHandler) Create(c *gin.Context) {
	api.Success(c, nil)
}