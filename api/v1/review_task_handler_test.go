package v1_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"unsafe"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	v1 "github.com/zgsm/go-webserver/api/v1"
	"github.com/zgsm/go-webserver/config"
	"github.com/zgsm/go-webserver/i18n"
	"github.com/zgsm/go-webserver/internal/service"
	"github.com/zgsm/go-webserver/pkg/types"
)

func init() {
	// 初始化i18n用于测试
	cfg := config.Config{}
	cfg.I18n.DefaultLocale = "zh-CN"
	cfg.I18n.BundlePath = "i18n/locales"
	if err := i18n.InitI18n(cfg); err != nil {
		panic(err)
	}
}

type mockReviewTaskService struct {
	createFunc         func(clientID, workspace string, targets []types.Target) (string, error)
	runFunc            func(clientID, workspace string, targets []types.Target) (string, error)
	issueIncrementFunc func(reviewTaskID, clientID string, offset int) (*types.IssueIncrementReviewTaskResult, error)
}

func (m *mockReviewTaskService) Create(clientID, workspace string, targets []types.Target) (string, error) {
	if m.createFunc != nil {
		return m.createFunc(clientID, workspace, targets)
	}
	return "", errors.New("mock create function not set")
}

func (m *mockReviewTaskService) Run(clientID, workspace string, targets []types.Target) (string, error) {
	if m.runFunc != nil {
		return m.runFunc(clientID, workspace, targets)
	}
	return "", errors.New("mock run function not set")
}

func (m *mockReviewTaskService) IssueIncrement(reviewTaskID, clientID string, offset int) (*types.IssueIncrementReviewTaskResult, error) {
	if m.issueIncrementFunc != nil {
		return m.issueIncrementFunc(reviewTaskID, clientID, offset)
	}
	return nil, errors.New("mock issueIncrement function not set")
}

func setReviewTaskService(handler *v1.ReviewTaskHandler, svc service.ReviewTaskService) {
	// 使用 unsafe 设置未导出字段
	handlerPtr := (*struct {
		reviewTaskService service.ReviewTaskService
	})(unsafe.Pointer(handler))
	handlerPtr.reviewTaskService = svc
}

func TestReviewTaskHandler_Create(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  string
		mockRun      func(clientID, workspace string, targets []types.Target) (string, error)
		expectedCode int
	}{
		{
			name: "success",
			requestBody: `{
				"client_id": "client123",
				"workspace": "workspace1",
				"targets": [{"id": "target1", "type": "file"}]
			}`,
			mockRun: func(clientID, workspace string, targets []types.Target) (string, error) {
				return "task-123", nil
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "invalid request",
			requestBody: `{
				"client_id": "",
				"workspace": "workspace1",
				"targets": []
			}`,
			mockRun:      nil,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "service error",
			requestBody: `{
				"client_id": "client123",
				"workspace": "workspace1",
				"targets": [{"id": "target1", "type": "file"}]
			}`,
			mockRun: func(clientID, workspace string, targets []types.Target) (string, error) {
				return "", errors.New("service error")
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock service
			mockSvc := &mockReviewTaskService{
				runFunc: tt.mockRun,
			}
			handler := &v1.ReviewTaskHandler{}
			setReviewTaskService(handler, mockSvc)

			// Setup Gin
			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.POST("/review_tasks", handler.Create)

			// Create request
			req, _ := http.NewRequest(http.MethodPost, "/review_tasks", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")

			// Record response
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Verify
			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}

func TestReviewTaskHandler_IssueIncrement(t *testing.T) {
	tests := []struct {
		name               string
		queryParams        string
		pathParam          string
		mockIssueIncrement func(reviewTaskID, clientID string, offset int) (*types.IssueIncrementReviewTaskResult, error)
		expectedCode       int
	}{
		{
			name:        "success",
			queryParams: "client_id=client123&offset=10",
			pathParam:   "task-123",
			mockIssueIncrement: func(reviewTaskID, clientID string, offset int) (*types.IssueIncrementReviewTaskResult, error) {
				return &types.IssueIncrementReviewTaskResult{
					IsDone:     false,
					Progress:   0.5,
					Total:      100,
					NextOffset: 10,
					Issues: []types.Issue{{
						IssueID:   "issue1",
						Message:   "test issue",
						FilePath:  "test.go",
						StartLine: 1,
						EndLine:   1,
						Severity:  "low",
						Status:    0,
						CreatedAt: "2025-01-01T00:00:00Z",
						UpdatedAt: "2025-01-01T00:00:00Z",
					}},
				}, nil
			},
			expectedCode: http.StatusOK,
		},
		{
			name:               "invalid params",
			queryParams:        "client_id=&offset=-1",
			pathParam:          "task-123",
			mockIssueIncrement: nil,
			expectedCode:       http.StatusBadRequest,
		},
		{
			name:        "service error",
			queryParams: "client_id=client123&offset=10",
			pathParam:   "task-123",
			mockIssueIncrement: func(reviewTaskID, clientID string, offset int) (*types.IssueIncrementReviewTaskResult, error) {
				return nil, errors.New("service error")
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock service
			mockSvc := &mockReviewTaskService{
				issueIncrementFunc: tt.mockIssueIncrement,
			}
			handler := &v1.ReviewTaskHandler{}
			setReviewTaskService(handler, mockSvc)

			// Setup Gin
			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.GET("/review_tasks/:review_task_id/issues/increment", handler.IssueIncrement)

			// Create request
			req, _ := http.NewRequest(http.MethodGet, "/review_tasks/"+tt.pathParam+"/issues/increment?"+tt.queryParams, nil)

			// Record response
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Verify
			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}
