package v1

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zgsm/go-webserver/api"
	"github.com/zgsm/go-webserver/pkg/i18nlogger"
	"github.com/zgsm/go-webserver/pkg/thirdPlatform"
)

// ExternalServiceHandler 外部服务处理器
// 这是一个示例处理器，展示如何在API中使用HTTP客户端服务
type ExternalServiceHandler struct{}

// NewExternalServiceHandler 创建外部服务处理器
func NewExternalServiceHandler() *ExternalServiceHandler {
	return &ExternalServiceHandler{}
}

// GetUserProfile 获取用户资料
// @Summary 获取外部用户资料
// @Description 从外部服务获取用户资料
// @Tags 外部服务
// @Accept json
// @Produce json
// @Param user_id path string true "用户ID"
// @Success 200 {object} api.Response{data=httpclient.UserProfile} "成功"
// @Failure 400 {object} api.Response "请求参数错误"
// @Failure 500 {object} api.Response "服务器内部错误"
// @Router /api/v1/external/users/{user_id}/profile [get]
func (h *ExternalServiceHandler) GetUserProfile(c *gin.Context) {
	// 获取路径参数
	userID := c.Param("user_id")
	if userID == "" {
		api.BadRequest(c, "external.user.id.empty")
		return
	}

	// 获取示例服务
	serverManager, err := thirdPlatform.GetServerManager()
	if err != nil {
		locale := i18nlogger.GetLocaleFromContext(c)
		i18nlogger.Error("external.service.get.failed", locale, nil, "error", err)
		api.ServerError(c, "external.service.get.failed")
		return
	}

	exampleService := serverManager.Example

	// 调用服务方法
	profile, err := exampleService.GetUserProfile(c.Request.Context(), userID)
	if err != nil {
		locale := i18nlogger.GetLocaleFromContext(c)
		i18nlogger.Error("external.user.profile.get.failed", locale, nil, "error", err, "user_id", userID)
		api.ServerError(c, "external.user.profile.get.failed")
		return
	}

	api.Success(c, profile)
}

// SearchUsers 搜索用户
// @Summary 搜索外部用户
// @Description 从外部服务搜索用户
// @Tags 外部服务
// @Accept json
// @Produce json
// @Param q query string true "搜索关键词"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} api.Response{data=httpclient.UserSearchResult} "成功"
// @Failure 400 {object} api.Response "请求参数错误"
// @Failure 500 {object} api.Response "服务器内部错误"
// @Router /api/v1/external/users/search [get]
func (h *ExternalServiceHandler) SearchUsers(c *gin.Context) {
	// 获取查询参数
	query := c.Query("q")
	if query == "" {
		api.BadRequest(c, "external.search.query.empty")
		return
	}

	// 获取分页参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// 获取示例服务
	serverManager, err := thirdPlatform.GetServerManager()
	if err != nil {
		locale := i18nlogger.GetLocaleFromContext(c)
		i18nlogger.Error("external.service.get.failed", locale, nil, "error", err)
		api.ServerError(c, "external.service.get.failed")
		return
	}

	exampleService := serverManager.Example

	// 调用服务方法
	result, err := exampleService.SearchUsers(c.Request.Context(), query, page, pageSize)
	if err != nil {
		locale := i18nlogger.GetLocaleFromContext(c)
		i18nlogger.Error("external.user.search.failed", locale, nil, "error", err, "query", query)
		api.ServerError(c, "external.user.search.failed")
		return
	}

	api.Success(c, result)
}

// UpdateUserProfile 更新用户资料
// @Summary 更新外部用户资料
// @Description 更新外部服务中的用户资料
// @Tags 外部服务
// @Accept json
// @Produce json
// @Param user_id path string true "用户ID"
// @Param profile body httpclient.UserProfile true "用户资料"
// @Success 200 {object} api.Response "成功"
// @Failure 400 {object} api.Response "请求参数错误"
// @Failure 500 {object} api.Response "服务器内部错误"
// @Router /api/v1/external/users/{user_id}/profile [put]
func (h *ExternalServiceHandler) UpdateUserProfile(c *gin.Context) {
	// 获取路径参数
	userID := c.Param("user_id")
	if userID == "" {
		api.BadRequest(c, "external.user.id.empty")
		return
	}

	// 解析请求体
	var profile thirdPlatform.ExampleUserProfile
	if err := c.ShouldBindJSON(&profile); err != nil {
		api.BadRequest(c, "external.user.profile.invalid")
		return
	}

	// 获取示例服务
	serverManager, err := thirdPlatform.GetServerManager()
	if err != nil {
		locale := i18nlogger.GetLocaleFromContext(c)
		i18nlogger.Error("external.service.get.failed", locale, nil, "error", err)
		api.ServerError(c, "external.service.get.failed")
		return
	}

	exampleService := serverManager.Example

	// 调用服务方法
	err = exampleService.UpdateUserProfile(c.Request.Context(), userID, &profile)
	if err != nil {
		locale := i18nlogger.GetLocaleFromContext(c)
		i18nlogger.Error("external.user.profile.update.failed", locale, nil, "error", err, "user_id", userID)
		api.ServerError(c, "external.user.profile.update.failed")
		return
	}

	api.Success(c, nil)
}
