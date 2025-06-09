package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zgsm/go-webserver/i18n"
)

// Response API response structure
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success success response
func Success(c *gin.Context, data interface{}) {
	locale, exists := c.Get("locale")
	var localeStr string
	if exists {
		localeStr = locale.(string)
	} else {
		localeStr = i18n.GetDefaultLocale()
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: i18n.Translate("common.success", localeStr, nil),
		Data:    data,
	})
}

// Fail failure response
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

	locale, exists := c.Get("locale")
	var localeStr string
	if exists {
		localeStr = locale.(string)
	} else {
		localeStr = i18n.GetDefaultLocale()
	}

	messageID := "common.internalServerError"
	var message string
	if err != nil {
		message = err.Error()
	} else {
		message = i18n.Translate(messageID, localeStr, nil)
	}

	c.JSON(code, Response{
		Code:    code,
		Message: message,
	})
}

// BadRequest 400 bad request response
func BadRequest(c *gin.Context, messageID string) {
	if messageID == "" {
		messageID = "common.badRequest"
	}
	Fail(c, http.StatusBadRequest, messageID)
}

// Unauthorized 401 unauthorized response
func Unauthorized(c *gin.Context, messageID string) {
	if messageID == "" {
		messageID = "common.unauthorized"
	}
	Fail(c, http.StatusUnauthorized, messageID)
}

// Forbidden 403 forbidden response
func Forbidden(c *gin.Context, messageID string) {
	if messageID == "" {
		messageID = "common.forbidden"
	}
	Fail(c, http.StatusForbidden, messageID)
}

// NotFound 404 not found response
func NotFound(c *gin.Context, messageID string) {
	if messageID == "" {
		messageID = "common.notFound"
	}
	Fail(c, http.StatusNotFound, messageID)
}

// ServerError 500 server error response
func ServerError(c *gin.Context, messageID string) {
	if messageID == "" {
		messageID = "common.serverError"
	}
	Fail(c, http.StatusInternalServerError, messageID)
}

// PageResult pagination result
type PageResult struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

// NewPageResult create pagination result
func NewPageResult(list interface{}, total int64, page, pageSize int) *PageResult {
	return &PageResult{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}
}
