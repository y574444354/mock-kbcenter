package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zgsm/go-webserver/api"
	"github.com/zgsm/go-webserver/internal/service"
	"github.com/zgsm/go-webserver/pkg/types"
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
	ClientID  string         `json:"client_id" binding:"required" example:"client123"`  // Client ID
	Workspace string         `json:"workspace" binding:"required" example:"workspace1"` // Workspace name
	Targets   []types.Target `json:"targets" binding:"required"`                        // List of review targets
}

// CreateReviewTaskResponse represents the response for creating a review task
type CreateReviewTaskResponse struct {
	ReviewTaskID string `json:"review_task_id" example:"task-123"` // Review task ID
}

// Create creates a new review task
// @Summary Create a new review task
// @Description Create a new code review task
// @Tags review_tasks
// @Accept json
// @Produce json
// @Param request body CreateReviewTaskRequest true "Request parameters for creating a review task"
// @Success 200 {object} CreateReviewTaskResponse
// @Failure 400 {object} api.Response
// @Failure 500 {object} api.Response
// @Router /review_tasks [post]
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
	ClientID string `form:"client_id" binding:"required" example:"client123"` // Client ID
	Offset   int    `form:"offset" binding:"required,min=0" example:"0"`      // Offset
}

// IssueIncrementReviewTaskResponse represents the response for incremental issues
type IssueIncrementReviewTaskResponse struct {
	types.IssueIncrementReviewTaskResult
}

// IssueIncrement get incremental issues
// @Summary Get incremental issues for a review task
// @Description Get incremental issues starting from a specified offset for a review task
// @Tags review_tasks
// @Accept json
// @Produce json
// @Param review_task_id path string true "Review task ID"
// @Param client_id query string true "Client ID"
// @Param offset query int true "Offset" minimum(0)
// @Success 200 {object} IssueIncrementReviewTaskResponse
// @Failure 400 {object} api.Response
// @Failure 500 {object} api.Response
// @Router /review_tasks/{review_task_id}/issues/increment [get]
func (h *ReviewTaskHandler) IssueIncrement(c *gin.Context) {
	reviewTaskID := c.Param("review_task_id")

	var req IssueIncrementReviewTaskRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		api.BadRequest(c, "common.invalidParameter")
		return
	}

	result, err := h.reviewTaskService.IssueIncrement(reviewTaskID, req.ClientID, req.Offset)
	if err != nil {
		api.Error(c, http.StatusInternalServerError, err)
		return
	}

	api.Success(c, result)
}
