package httpclient

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/zgsm/mock-kbcenter/i18n"
	"github.com/zgsm/mock-kbcenter/pkg/logger"
)

// Middleware defines HTTP client middleware interface
type Middleware interface {
	// ProcessRequest process request, called before sending request
	ProcessRequest(*http.Request) error
	// ProcessResponse process response, called after receiving response
	ProcessResponse(*http.Response, error) (*http.Response, error)
}

// StatusCodeMiddleware status code validation middleware
type StatusCodeMiddleware struct {
	// ValidStatusCodes allowed status code ranges, if nil use default rules
	ValidStatusCodes []int
}

// ProcessRequest process request
func (m *StatusCodeMiddleware) ProcessRequest(req *http.Request) error {
	// Status code middleware doesn't process requests
	return nil
}

// ProcessResponse validate status code
func (m *StatusCodeMiddleware) ProcessResponse(resp *http.Response, err error) (*http.Response, error) {
	if err != nil {
		return resp, err
	}

	// If no valid status codes specified, use default rules (2xx is valid)
	if len(m.ValidStatusCodes) == 0 {
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			// Read response body content
			var bodyContent string
			if resp.Body != nil {
				bodyBytes, err := io.ReadAll(resp.Body)
				if err == nil && bodyBytes != nil {
					bodyContent = string(bodyBytes)
					// Reset response body for further processing
					resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
				}
			}

			url := ""
			if resp.Request != nil {
				url = resp.Request.URL.String()
			}
			return nil, &StatusError{
				StatusCode: resp.StatusCode,
				Message: i18n.Translate("httpclient.error.invalid_status_code", "", map[string]interface{}{
					"status": resp.StatusCode,
					"url":    url,
					"body":   bodyContent,
				}),
			}
		}
		return resp, nil
	}

	// Check if status code is in allowed range
	for _, code := range m.ValidStatusCodes {
		if resp.StatusCode == code {
			return resp, nil
		}
	}

	// Read response body content
	var bodyContent string
	if resp.Body != nil {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err == nil && bodyBytes != nil {
			bodyContent = string(bodyBytes)
			// Reset response body for further processing
			resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}
	}

	url := ""
	if resp.Request != nil {
		url = resp.Request.URL.String()
	}
	return nil, &StatusError{
		StatusCode: resp.StatusCode,
		Message: i18n.Translate("httpclient.error.invalid_status_code", "", map[string]interface{}{
			"status": resp.StatusCode,
			"url":    url,
			"body":   bodyContent,
		}),
	}
}

// StatusError status code error
type StatusError struct {
	StatusCode int
	Message    string
}

func (e *StatusError) Error() string {
	return e.Message
}

// LogMiddleware logging middleware
type LogMiddleware struct {
	EnableRequestLog  bool
	EnableResponseLog bool
}

// ProcessRequest log request
func (m *LogMiddleware) ProcessRequest(req *http.Request) error {
	if !m.EnableRequestLog {
		return nil
	}

	// Log request info
	logger.Debug(i18n.Translate("httpclient.log.request", "", nil),
		"method", req.Method,
		"url", req.URL.String(),
		"headers", req.Header,
	)

	// If request body not empty, log request body
	if req.Body != nil && req.Body != http.NoBody {
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			return err
		}
		// Reset request body for further processing
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		logger.Debug(i18n.Translate("httpclient.log.request_body", "", nil), "body", string(bodyBytes))
	}

	return nil
}

// ProcessResponse log response
func (m *LogMiddleware) ProcessResponse(resp *http.Response, err error) (*http.Response, error) {
	if !m.EnableResponseLog || err != nil {
		return resp, err
	}

	// Choose log level based on status code
	logFunc := logger.Debug
	if resp.StatusCode >= 400 {
		logFunc = logger.Warn
	}

	// Log response info
	logFunc(i18n.Translate("httpclient.log.response", "", nil),
		"status", resp.Status,
		"url", resp.Request.URL.String(),
		// "headers", resp.Header,
	)

	// If response body not empty, log response body
	// if resp.Body != nil {
	// 	bodyBytes, err := io.ReadAll(resp.Body)
	// 	if err != nil {
	// 		return resp, err
	// 	}
	// 	// Reset response body for further processing
	// 	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	// 	logger.Debug(i18n.Translate("httpclient.log.response_body", "", nil), "body", string(bodyBytes))
	// }

	return resp, err
}

// AuthMiddleware authentication middleware
type AuthMiddleware struct {
	AuthType   string
	Username   string
	Password   string
	Token      string
	AuthHeader string
}

// ProcessRequest add authentication info
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

// ProcessResponse process response
func (m *AuthMiddleware) ProcessResponse(resp *http.Response, err error) (*http.Response, error) {
	// Auth middleware doesn't process response
	return resp, err
}

// RetryMiddleware retry middleware
type RetryMiddleware struct {
	MaxRetries int
	RetryDelay time.Duration
}

// ProcessRequest process request
func (m *RetryMiddleware) ProcessRequest(req *http.Request) error {
	// Retry middleware doesn't process request
	return nil
}

// ProcessResponse process response, return error if retry needed
func (m *RetryMiddleware) ProcessResponse(resp *http.Response, err error) (*http.Response, error) {
	// Retry logic implemented in Client
	return resp, err
}

// HeaderMiddleware header middleware
type HeaderMiddleware struct {
	Headers map[string]string
}

// ProcessRequest add headers
func (m *HeaderMiddleware) ProcessRequest(req *http.Request) error {
	if req == nil {
		return errors.New("http request cannot be nil")
	}
	for key, value := range m.Headers {
		req.Header.Set(key, value)
	}
	return nil
}

// ProcessResponse process response
func (m *HeaderMiddleware) ProcessResponse(resp *http.Response, err error) (*http.Response, error) {
	// Header middleware doesn't process response
	return resp, err
}
