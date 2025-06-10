package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zgsm/mock-kbcenter/api"
	"github.com/zgsm/mock-kbcenter/internal/service"
	"github.com/zgsm/mock-kbcenter/pkg/language"
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

type FileStructureRequest struct {
	ClientId     string `form:"clientId" binding:"required"`
	CodebasePath string `form:"codebasePath" binding:"required"`
	FilePath     string `form:"filePath"`
}

type FunctionDefinition struct {
	Type     string `json:"type"`
	Name     string `json:"name"`
	Position struct {
		StartLine   int `json:"startLine"`
		StartColumn int `json:"startColumn"`
		EndLine     int `json:"endLine"`
		EndColumn   int `json:"endColumn"`
	} `json:"position"`
	Content string `json:"content"`
}

func (h *KBCenterMockHandler) GetFileStructure(c *gin.Context) {
	var req FileStructureRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		api.BadRequest(c, "common.invalidParams")
		return
	}

	funcs, err := h.service.GetFileStructure(c.Request.Context(), req.FilePath)
	if err != nil {
		api.Error(c, http.StatusInternalServerError, err)
		return
	}

	var list []FunctionDefinition
	// Get language from file extension
	lang, err := language.Detect(req.FilePath)
	if err != nil {
		api.Error(c, http.StatusInternalServerError, err)
		return
	}

	for _, f := range funcs {
		funcName, err := language.GetFunctionName(lang, f.Code)
		if err != nil {
			api.Error(c, http.StatusInternalServerError, err)
			return
		}

		fd := FunctionDefinition{
			Type: "function_definition",
			Name: funcName,
			Position: struct {
				StartLine   int `json:"startLine"`
				StartColumn int `json:"startColumn"`
				EndLine     int `json:"endLine"`
				EndColumn   int `json:"endColumn"`
			}{
				StartLine:   f.StartLine,
				StartColumn: 0, // Default to 0 for now
				EndLine:     f.EndLine,
				EndColumn:   0, // Default to 1 for now
			},
			Content: f.Code,
		}
		list = append(list, fd)
	}

	api.Success(c, gin.H{
		"list": list,
	})
}

func (h *KBCenterMockHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/files/content", h.GetFileContent)
	router.GET("/codebases/directory", h.GetDirectoryTree)
	router.GET("/files/structure", h.GetFileStructure)
}
