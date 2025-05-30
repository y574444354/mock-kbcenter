package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zgsm/go-webserver/i18n"
)

// Response API响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	locale, exists := c.Get("locale")
	var localeStr string
	if exists {
		localeStr = locale.(string)
	} else {
		localeStr = i18n.GetDefaultLocale()
	}

	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: i18n.Translate("common.success", localeStr, nil),
		Data:    data,
	})
}

// Fail 失败响应
func Fail(c *gin.Context, code int, messageID string) {
	locale, exists := c.Get("locale")
	var localeStr string
	if exists {
		localeStr = locale.(string)
	} else {
		localeStr = i18n.GetDefaultLocale()
	}

	message := i18n.Translate(messageID, localeStr, nil)
	c.JSON(code, Response{
		Code:    code,
		Message: message,
	})
}

func Error(c *gin.Context, code int, err error) {
	c.Errors = append(c.Errors, &gin.Error{Err: err})
	c.JSON(code, Response{
		Code:    code,
		Message: err.Error(),
	})
}

// BadRequest 400错误响应
func BadRequest(c *gin.Context, messageID string) {
	if messageID == "" {
		messageID = "common.badRequest"
	}
	Fail(c, http.StatusBadRequest, messageID)
}

// Unauthorized 401错误响应
func Unauthorized(c *gin.Context, messageID string) {
	if messageID == "" {
		messageID = "common.unauthorized"
	}
	Fail(c, http.StatusUnauthorized, messageID)
}

// Forbidden 403错误响应
func Forbidden(c *gin.Context, messageID string) {
	if messageID == "" {
		messageID = "common.forbidden"
	}
	Fail(c, http.StatusForbidden, messageID)
}

// NotFound 404错误响应
func NotFound(c *gin.Context, messageID string) {
	if messageID == "" {
		messageID = "common.notFound"
	}
	Fail(c, http.StatusNotFound, messageID)
}

// ServerError 500错误响应
func ServerError(c *gin.Context, messageID string) {
	if messageID == "" {
		messageID = "common.serverError"
	}
	Fail(c, http.StatusInternalServerError, messageID)
}

// PageResult 分页结果
type PageResult struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

// NewPageResult 创建分页结果
func NewPageResult(list interface{}, total int64, page, pageSize int) *PageResult {
	return &PageResult{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}
}
