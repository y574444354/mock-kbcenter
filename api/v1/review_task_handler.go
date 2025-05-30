package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zgsm/mock-kbcenter/api"
	"github.com/zgsm/mock-kbcenter/internal/service"
	"github.com/zgsm/mock-kbcenter/pkg/types"
)

type ReviewTaskHandler struct {
	reviewTaskService service.ReviewTaskService
}

func NewReviewTaskHandler() *ReviewTaskHandler {
	return &ReviewTaskHandler{
		reviewTaskService: service.NewReviewTaskService(),
	}
}

type CreateReviewTaskRequest struct {
	ClientID  string         `json:"client_id" binding:"required"`
	Workspace string         `json:"workspace" binding:"required"`
	Targets   []types.Target `json:"targets" binding:"required"`
}

type CreateReviewTaskResponse struct {
	ReviewTaskID string `json:"review_task_id"`
}

func (h *ReviewTaskHandler) Create(c *gin.Context) {
	var req CreateReviewTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.BadRequest(c, "common.invalidParameter")
		return
	}

	taskID, err := h.reviewTaskService.Run(req.ClientID, req.Workspace, req.Targets)
	if err != nil {
		api.Error(c, http.StatusInternalServerError, err)
		return
	}

	api.Success(c, CreateReviewTaskResponse{
		ReviewTaskID: taskID,
	})
}

type IssueIncrementReviewTaskRequest struct {
	ClientId string `form:"client_id" binding:"required"`
	Offset   int    `form:"offset" binding:"required,min=0"`
}

type IssueIncrementReviewTaskResponse struct {
	types.IssueIncrementReviewTaskResult
}

func (h *ReviewTaskHandler) IssueIncrement(c *gin.Context) {
	reviewTaskID := c.Param("review_task_id")

	var req IssueIncrementReviewTaskRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		api.BadRequest(c, "common.invalidParameter")
		return
	}

	result, err := h.reviewTaskService.IssueIncrement(reviewTaskID, req.ClientId, req.Offset)
	if err != nil {
		api.Error(c, http.StatusInternalServerError, err)
		return
	}

	api.Success(c, result)
}
