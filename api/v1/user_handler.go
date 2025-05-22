package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zgsm/review-manager/api"
	"github.com/zgsm/review-manager/internal/model"
	"github.com/zgsm/review-manager/internal/service"
	"github.com/zgsm/review-manager/pkg/i18nlogger"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService service.UserService
}

// NewUserHandler 创建用户处理器
func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: service.NewUserService(),
	}
}

// Register 用户注册
// @Summary 用户注册
// @Description 创建新用户
// @Tags 用户
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "注册信息"
// @Success 200 {object} api.Response{data=model.User} "成功"
// @Failure 400 {object} api.Response "请求参数错误"
// @Failure 500 {object} api.Response "服务器内部错误"
// @Router /api/v1/users/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.BadRequest(c, "user.register.invalid_params")
		return
	}

	// 参数验证
	if err := h.userService.ValidateRegisterParams(req.Username, req.Email, req.Password); err != nil {
		api.BadRequest(c, err.Error())
		return
	}

	// 调用服务
	user, err := h.userService.Register(c.Request.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		locale := i18nlogger.GetLocaleFromContext(c)
		i18nlogger.Error("user.register.failed", locale, nil, "error", err)
		api.Fail(c, http.StatusInternalServerError, "user.register.failed")
		return
	}

	api.Success(c, user)
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录并返回token
// @Tags 用户
// @Accept json
// @Produce json
// @Param request body LoginRequest true "登录信息"
// @Success 200 {object} api.Response{data=LoginResponse} "成功"
// @Failure 400 {object} api.Response "请求参数错误"
// @Failure 401 {object} api.Response "未授权"
// @Failure 500 {object} api.Response "服务器内部错误"
// @Router /api/v1/users/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.BadRequest(c, "user.login.invalid_params")
		return
	}

	// 参数验证
	if err := h.userService.ValidateLoginParams(req.Username, req.Password); err != nil {
		api.BadRequest(c, err.Error())
		return
	}

	// 调用服务
	user, err, token := h.userService.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		locale := i18nlogger.GetLocaleFromContext(c)
		i18nlogger.Error("user.login.failed", locale, nil, "error", err)
		api.Unauthorized(c, "user.login.failed")
		return
	}

	api.Success(c, LoginResponse{
		User:  user,
		Token: token,
	})
}

// GetUserInfo 获取用户信息
// @Summary 获取用户信息
// @Description 获取当前登录用户的信息
// @Tags 用户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} api.Response{data=model.User} "成功"
// @Failure 401 {object} api.Response "未授权"
// @Failure 500 {object} api.Response "服务器内部错误"
// @Router /api/v1/users/info [get]
func (h *UserHandler) GetUserInfo(c *gin.Context) {
	// 从上下文中获取用户ID（实际项目中应从JWT等认证信息中获取）
	userID := uint(1) // 示例用户ID

	// 调用服务
	user, err := h.userService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		locale := i18nlogger.GetLocaleFromContext(c)
		i18nlogger.Error("user.info.get.failed", locale, nil, "error", err)
		api.ServerError(c, "user.info.get.failed")
		return
	}

	if user == nil {
		api.NotFound(c, "user.not_found")
		return
	}

	api.Success(c, user)
}

// UpdateUserInfo 更新用户信息
// @Summary 更新用户信息
// @Description 更新当前登录用户的信息
// @Tags 用户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body UpdateUserRequest true "用户信息"
// @Success 200 {object} api.Response "成功"
// @Failure 400 {object} api.Response "请求参数错误"
// @Failure 401 {object} api.Response "未授权"
// @Failure 500 {object} api.Response "服务器内部错误"
// @Router /api/v1/users/info [put]
func (h *UserHandler) UpdateUserInfo(c *gin.Context) {
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.BadRequest(c, "user.update.invalid_params")
		return
	}

	// 从上下文中获取用户ID（实际项目中应从JWT等认证信息中获取）
	userID := uint(1) // 示例用户ID

	// 调用服务
	if err := h.userService.UpdateUserInfo(c.Request.Context(), userID, req.Nickname, req.Avatar); err != nil {
		locale := i18nlogger.GetLocaleFromContext(c)
		i18nlogger.Error("user.info.update.failed", locale, nil, "error", err)
		api.ServerError(c, "user.info.update.failed")
		return
	}

	api.Success(c, nil)
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Description 修改当前登录用户的密码
// @Tags 用户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body ChangePasswordRequest true "密码信息"
// @Success 200 {object} api.Response "成功"
// @Failure 400 {object} api.Response "请求参数错误"
// @Failure 401 {object} api.Response "未授权"
// @Failure 500 {object} api.Response "服务器内部错误"
// @Router /api/v1/users/password [put]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.BadRequest(c, "user.password.invalid_params")
		return
	}

	// 参数验证
	if err := h.userService.ValidateChangePasswordParams(req.OldPassword, req.NewPassword); err != nil {
		api.BadRequest(c, err.Error())
		return
	}

	// 从上下文中获取用户ID（实际项目中应从JWT等认证信息中获取）
	userID := uint(1) // 示例用户ID

	// 调用服务
	if err := h.userService.ChangePassword(c.Request.Context(), userID, req.OldPassword, req.NewPassword); err != nil {
		locale := i18nlogger.GetLocaleFromContext(c)
		i18nlogger.Error("user.password.change.failed", locale, nil, "error", err)
		api.Fail(c, http.StatusInternalServerError, "user.password.change.failed")
		return
	}

	api.Success(c, nil)
}

// ListUsers 获取用户列表
// @Summary 获取用户列表
// @Description 分页获取用户列表
// @Tags 用户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} api.Response{data=api.PageResult{list=[]model.User}} "成功"
// @Failure 401 {object} api.Response "未授权"
// @Failure 500 {object} api.Response "服务器内部错误"
// @Router /api/v1/users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	// 获取分页参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	// 验证并解析分页参数
	page, pageSize := h.userService.ValidateAndParsePageParams(pageStr, pageSizeStr)

	// 调用服务
	users, total, err := h.userService.ListUsers(c.Request.Context(), page, pageSize)
	if err != nil {
		locale := i18nlogger.GetLocaleFromContext(c)
		i18nlogger.Error("user.list.failed", locale, nil, "error", err)
		api.ServerError(c, "user.list.failed")
		return
	}

	api.Success(c, api.NewPageResult(users, total, page, pageSize))
}

// DeleteUser 删除用户
// @Summary 删除用户
// @Description 删除指定ID的用户
// @Tags 用户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "用户ID"
// @Success 200 {object} api.Response "成功"
// @Failure 400 {object} api.Response "请求参数错误"
// @Failure 401 {object} api.Response "未授权"
// @Failure 403 {object} api.Response "禁止访问"
// @Failure 500 {object} api.Response "服务器内部错误"
// @Router /api/v1/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	// 获取用户ID
	idStr := c.Param("id")
	id, err := h.userService.ValidateAndParseUserID(idStr)
	if err != nil {
		api.BadRequest(c, err.Error())
		return
	}

	// 调用服务
	if err := h.userService.DeleteUser(c.Request.Context(), id); err != nil {
		locale := i18nlogger.GetLocaleFromContext(c)
		i18nlogger.Error("user.delete.failed", locale, nil, "error", err)
		api.ServerError(c, "user.delete.failed")
		return
	}

	api.Success(c, nil)
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	User  *model.User `json:"user"`
	Token string      `json:"token"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
