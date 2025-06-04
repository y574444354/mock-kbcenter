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
	ClientID  string         `json:"client_id" binding:"required" example:"client123"`  // 客户端ID
	Workspace string         `json:"workspace" binding:"required" example:"workspace1"` // 工作空间名称
	Targets   []types.Target `json:"targets" binding:"required"`                        // 评审目标列表
}

// CreateReviewTaskResponse 创建评审任务响应
type CreateReviewTaskResponse struct {
	ReviewTaskID string `json:"review_task_id" example:"task-123"` // 评审任务ID
}

// Create 创建评审任务
// @Summary 创建新的评审任务
// @Description 创建一个新的代码评审任务
// @Tags review_tasks
// @Accept json
// @Produce json
// @Param request body CreateReviewTaskRequest true "创建评审任务请求参数"
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
	ClientId string `form:"client_id" binding:"required" example:"client123"` // 客户端ID
	Offset   int    `form:"offset" binding:"required,min=0" example:"0"`      // 偏移量
}

// IssueIncrementReviewTaskResponse 增量问题响应
type IssueIncrementReviewTaskResponse struct {
	types.IssueIncrementReviewTaskResult
}

// IssueIncrement 获取增量问题
// @Summary 获取评审任务的增量问题
// @Description 获取指定评审任务从某个偏移量开始的增量问题
// @Tags review_tasks
// @Accept json
// @Produce json
// @Param review_task_id path string true "评审任务ID"
// @Param client_id query string true "客户端ID"
// @Param offset query int true "偏移量" minimum(0)
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

	result, err := h.reviewTaskService.IssueIncrement(reviewTaskID, req.ClientId, req.Offset)
	if err != nil {
		api.Error(c, http.StatusInternalServerError, err)
		return
	}

	api.Success(c, result)
}
