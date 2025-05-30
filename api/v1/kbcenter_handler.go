package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zgsm/mock-kbcenter/api"
	"github.com/zgsm/mock-kbcenter/internal/service"
)

type KBCenterMockHandler struct {
	service *service.KBCenterMockService
}

func NewKBCenterMockHandler(baseDir string) *KBCenterMockHandler {
	return &KBCenterMockHandler{
		service: service.NewKBCenterMockService(baseDir),
	}
}

func (h *KBCenterMockHandler) GetFileContent(c *gin.Context) {
	filePath := c.Query("filePath")
	startLine, _ := strconv.Atoi(c.Query("startLine"))
	endLine, _ := strconv.Atoi(c.Query("endLine"))

	content, err := h.service.GetFileContent(c.Request.Context(), filePath, startLine, endLine)
	if err != nil {
		api.Error(c, http.StatusInternalServerError, err)
		return
	}

	c.Data(http.StatusOK, "text/plain", content)
}

func (h *KBCenterMockHandler) GetDirectoryTree(c *gin.Context) {
	clientId := c.Query("clientId")
	projectPath := c.Query("projectPath")
	subDir := c.Query("subDir")
	depth, _ := strconv.Atoi(c.Query("depth"))
	includeFiles := c.Query("includeFiles") != "0"

	result, err := h.service.GetDirectoryTree(c.Request.Context(), clientId, projectPath, subDir, depth, includeFiles)
	if err != nil {
		api.Error(c, http.StatusInternalServerError, err)
		return
	}

	api.Success(c, result)
}

func (h *KBCenterMockHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/codebase-indexer/api/v1/files/content", h.GetFileContent)
	router.GET("/codebase-indexer/api/v1/codebases/directory", h.GetDirectoryTree)
}
