package httpclient

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/zgsm/go-webserver/i18n"
	"github.com/zgsm/go-webserver/pkg/logger"
)

// Middleware 定义HTTP客户端中间件接口
type Middleware interface {
	// ProcessRequest 处理请求，在请求发送前调用
	ProcessRequest(*http.Request) error
	// ProcessResponse 处理响应，在响应接收后调用
	ProcessResponse(*http.Response, error) (*http.Response, error)
}

// LogMiddleware 日志中间件
type LogMiddleware struct {
	EnableRequestLog  bool
	EnableResponseLog bool
}

// ProcessRequest 记录请求日志
func (m *LogMiddleware) ProcessRequest(req *http.Request) error {
	if !m.EnableRequestLog {
		return nil
	}

	// 记录请求信息
	logger.Info(i18n.Translate("httpclient.log.request", "", nil),
		"method", req.Method,
		"url", req.URL.String(),
		"headers", req.Header,
	)

	// 如果请求体不为空，记录请求体
	if req.Body != nil && req.Body != http.NoBody {
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			return err
		}
		// 重置请求体，以便后续处理
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		logger.Debug(i18n.Translate("httpclient.log.request_body", "", nil), "body", string(bodyBytes))
	}

	return nil
}

// ProcessResponse 记录响应日志
func (m *LogMiddleware) ProcessResponse(resp *http.Response, err error) (*http.Response, error) {
	if !m.EnableResponseLog || err != nil {
		return resp, err
	}

	// 记录响应信息
	logger.Info(i18n.Translate("httpclient.log.response", "", nil),
		"status", resp.Status,
		"headers", resp.Header,
	)

	// 如果响应体不为空，记录响应体
	if resp.Body != nil {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return resp, err
		}
		// 重置响应体，以便后续处理
		resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		logger.Debug(i18n.Translate("httpclient.log.response_body", "", nil), "body", string(bodyBytes))
	}

	return resp, err
}

// AuthMiddleware 认证中间件
type AuthMiddleware struct {
	AuthType   string
	Username   string
	Password   string
	Token      string
	AuthHeader string
}

// ProcessRequest 添加认证信息
func (m *AuthMiddleware) ProcessRequest(req *http.Request) error {
	switch m.AuthType {
	case "basic":
		req.SetBasicAuth(m.Username, m.Password)
	case "bearer":
		req.Header.Set("Authorization", "Bearer "+m.Token)
	case "custom":
		req.Header.Set("Authorization", m.AuthHeader)
	}
	return nil
}

// ProcessResponse 处理响应
func (m *AuthMiddleware) ProcessResponse(resp *http.Response, err error) (*http.Response, error) {
	// 认证中间件不处理响应
	return resp, err
}

// RetryMiddleware 重试中间件
type RetryMiddleware struct {
	MaxRetries int
	RetryDelay time.Duration
}

// ProcessRequest 处理请求
func (m *RetryMiddleware) ProcessRequest(req *http.Request) error {
	// 重试中间件不处理请求
	return nil
}

// ProcessResponse 处理响应，如果需要重试则返回错误
func (m *RetryMiddleware) ProcessResponse(resp *http.Response, err error) (*http.Response, error) {
	// 在Client中实现重试逻辑
	return resp, err
}

// HeaderMiddleware 请求头中间件
type HeaderMiddleware struct {
	Headers map[string]string
}

// ProcessRequest 添加请求头
func (m *HeaderMiddleware) ProcessRequest(req *http.Request) error {
	for key, value := range m.Headers {
		req.Header.Set(key, value)
	}
	return nil
}

// ProcessResponse 处理响应
func (m *HeaderMiddleware) ProcessResponse(resp *http.Response, err error) (*http.Response, error) {
	// 请求头中间件不处理响应
	return resp, err
}
