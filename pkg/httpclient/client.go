package httpclient

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/zgsm/go-webserver/i18n"
	"github.com/zgsm/go-webserver/pkg/logger"
)

// Client HTTP客户端
type Client struct {
	config      *HttpServiceConfig
	httpClient  *http.Client
	middlewares []Middleware
}

// NewClient 创建新的HTTP客户端
func NewClient(config *HttpServiceConfig) (*Client, error) {
	if config == nil {
		config = DefaultHttpServiceConfig()
	}

	// 创建Transport
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
	}

	// 配置代理
	if config.ProxyURL != "" {
		proxyURL, err := url.Parse(config.ProxyURL)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", i18n.Translate("httpclient.proxy.parse_failed", "", nil), err)
		}
		transport.Proxy = http.ProxyURL(proxyURL)
	}

	// 配置TLS
	if config.InsecureSkipVerify || config.CertFile != "" || config.CAFile != "" {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: config.InsecureSkipVerify,
		}

		// 加载客户端证书
		if config.CertFile != "" && config.KeyFile != "" {
			cert, err := tls.LoadX509KeyPair(config.CertFile, config.KeyFile)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", i18n.Translate("httpclient.cert.load_failed", "", nil), err)
			}
			tlsConfig.Certificates = []tls.Certificate{cert}
		}

		// 加载CA证书
		if config.CAFile != "" {
			caCert, err := os.ReadFile(config.CAFile)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", i18n.Translate("httpclient.ca.read_failed", "", nil), err)
			}
			caCertPool := x509.NewCertPool()
			caCertPool.AppendCertsFromPEM(caCert)
			tlsConfig.RootCAs = caCertPool
		}

		transport.TLSClientConfig = tlsConfig
	}

	// 创建HTTP客户端
	httpClient := &http.Client{
		Timeout:   config.Timeout,
		Transport: transport,
	}

	// 创建客户端实例
	client := &Client{
		config:      config,
		httpClient:  httpClient,
		middlewares: make([]Middleware, 0),
	}

	// 添加默认中间件
	client.AddMiddleware(&LogMiddleware{
		EnableRequestLog:  config.EnableRequestLog,
		EnableResponseLog: config.EnableResponseLog,
	})

	client.AddMiddleware(&HeaderMiddleware{
		Headers: config.Headers,
	})

	// 添加状态码校验中间件
	client.AddMiddleware(&StatusCodeMiddleware{
		ValidStatusCodes: config.ValidStatusCodes,
	})

	// 添加认证中间件
	if config.AuthType != "none" {
		client.AddMiddleware(&AuthMiddleware{
			AuthType:   config.AuthType,
			Username:   config.Username,
			Password:   config.Password,
			Token:      config.Token,
			AuthHeader: config.AuthHeader,
		})
	}

	// 添加重试中间件
	if config.MaxRetries > 0 {
		client.AddMiddleware(&RetryMiddleware{
			MaxRetries: config.MaxRetries,
			RetryDelay: config.RetryDelay,
		})
	}

	return client, nil
}

// AddMiddleware 添加中间件
func (c *Client) AddMiddleware(middleware Middleware) {
	c.middlewares = append(c.middlewares, middleware)
}

// Request 发送HTTP请求
func (c *Client) Request(ctx context.Context, method, path string, body interface{}, headers map[string]string) (*http.Response, error) {
	// 构建完整URL
	fullURL := path
	if !strings.HasPrefix(path, "http") && c.config.BaseURL != "" {
		fullURL = c.config.BaseURL + path
	}

	// 准备请求体
	var reqBody io.Reader
	if body != nil {
		switch v := body.(type) {
		case string:
			reqBody = strings.NewReader(v)
		case []byte:
			reqBody = bytes.NewReader(v)
		case io.Reader:
			reqBody = v
		default:
			// 默认将对象序列化为JSON
			jsonData, err := json.Marshal(body)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", i18n.Translate("httpclient.request.serialize_failed", "", nil), err)
			}
			reqBody = bytes.NewReader(jsonData)
			// 如果没有指定Content-Type，则设置为application/json
			if headers == nil {
				headers = make(map[string]string)
			}
			if _, ok := headers["Content-Type"]; !ok {
				headers["Content-Type"] = "application/json"
			}
		}
	}

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, method, fullURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", i18n.Translate("httpclient.request.create_failed", "", nil), err)
	}

	// 添加自定义请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// 应用请求中间件
	for _, middleware := range c.middlewares {
		if err := middleware.ProcessRequest(req); err != nil {
			return nil, fmt.Errorf("%s: %w", i18n.Translate("httpclient.middleware.process_failed", "", nil), err)
		}
	}

	// 发送请求
	var resp *http.Response
	var respErr error

	// 实现重试逻辑
	retries := 0
	for {
		resp, respErr = c.httpClient.Do(req)

		// 应用响应中间件
		for i := len(c.middlewares) - 1; i >= 0; i-- {
			resp, respErr = c.middlewares[i].ProcessResponse(resp, respErr)
		}

		// 检查是否需要重试
		shouldRetry := false
		if c.config.MaxRetries > retries {
			if respErr != nil {
				// 网络错误重试
				shouldRetry = true
			} else if resp.StatusCode >= 500 {
				// 服务器错误重试
				shouldRetry = true
			}
		}

		if !shouldRetry {
			break
		}

		// 关闭响应体，准备重试
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}

		retries++
		logger.Info(i18n.Translate("httpclient.retry.attempt", "", map[string]interface{}{
			"attempt": retries,
			"max":     c.config.MaxRetries,
		}))
		time.Sleep(c.config.RetryDelay)

		// 重新创建请求体
		if body != nil {
			switch v := body.(type) {
			case string:
				reqBody = strings.NewReader(v)
			case []byte:
				reqBody = bytes.NewReader(v)
			case io.Reader:
				// 无法重置io.Reader，这是一个限制
				return nil, errors.New(i18n.Translate("httpclient.retry.reader_unsupported", "", nil))
			default:
				// 重新序列化JSON
				jsonData, _ := json.Marshal(body)
				reqBody = bytes.NewReader(jsonData)
			}
			req, err = http.NewRequestWithContext(ctx, method, fullURL, reqBody)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", i18n.Translate("httpclient.retry.request_failed", "", nil), err)
			}
			// 重新添加请求头
			for key, value := range headers {
				req.Header.Set(key, value)
			}
			// 重新应用请求中间件
			for _, middleware := range c.middlewares {
				if err := middleware.ProcessRequest(req); err != nil {
					return nil, fmt.Errorf("%s: %w", i18n.Translate("httpclient.retry.middleware_failed", "", nil), err)
				}
			}
		}
	}

	return resp, respErr
}

// Get 发送GET请求
func (c *Client) Get(ctx context.Context, path string, headers map[string]string) (*http.Response, error) {
	return c.Request(ctx, http.MethodGet, path, nil, headers)
}

// Post 发送POST请求
func (c *Client) Post(ctx context.Context, path string, body interface{}, headers map[string]string) (*http.Response, error) {
	return c.Request(ctx, http.MethodPost, path, body, headers)
}

// Put 发送PUT请求
func (c *Client) Put(ctx context.Context, path string, body interface{}, headers map[string]string) (*http.Response, error) {
	return c.Request(ctx, http.MethodPut, path, body, headers)
}

// Delete 发送DELETE请求
func (c *Client) Delete(ctx context.Context, path string, headers map[string]string) (*http.Response, error) {
	return c.Request(ctx, http.MethodDelete, path, nil, headers)
}

// Patch 发送PATCH请求
func (c *Client) Patch(ctx context.Context, path string, body interface{}, headers map[string]string) (*http.Response, error) {
	return c.Request(ctx, http.MethodPatch, path, body, headers)
}

// GetJSON 发送GET请求并解析JSON响应
func (c *Client) GetJSON(ctx context.Context, path string, headers map[string]string, v interface{}) error {
	resp, err := c.Get(ctx, path, headers)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return c.parseJSONResponse(resp, v)
}

// PostJSON 发送POST请求并解析JSON响应
func (c *Client) PostJSON(ctx context.Context, path string, body interface{}, headers map[string]string, v interface{}) error {
	resp, err := c.Post(ctx, path, body, headers)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return c.parseJSONResponse(resp, v)
}

// PutJSON 发送PUT请求并解析JSON响应
func (c *Client) PutJSON(ctx context.Context, path string, body interface{}, headers map[string]string, v interface{}) error {
	resp, err := c.Put(ctx, path, body, headers)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return c.parseJSONResponse(resp, v)
}

// DeleteJSON 发送DELETE请求并解析JSON响应
func (c *Client) DeleteJSON(ctx context.Context, path string, headers map[string]string, v interface{}) error {
	resp, err := c.Delete(ctx, path, headers)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return c.parseJSONResponse(resp, v)
}

// PatchJSON 发送PATCH请求并解析JSON响应
func (c *Client) PatchJSON(ctx context.Context, path string, body interface{}, headers map[string]string, v interface{}) error {
	resp, err := c.Patch(ctx, path, body, headers)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return c.parseJSONResponse(resp, v)
}

// parseJSONResponse 解析JSON响应
func (c *Client) parseJSONResponse(resp *http.Response, v interface{}) error {
	// 检查状态码
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf(i18n.Translate("httpclient.response.read_failed", "", map[string]interface{}{"code": resp.StatusCode})+": %w", err)
		}
		return fmt.Errorf("%s", i18n.Translate("httpclient.response.failed", "", map[string]interface{}{
			"code": resp.StatusCode,
			"body": string(bodyBytes),
		}))
	}

	// 解析JSON响应
	if v != nil {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return fmt.Errorf("%s: %w", i18n.Translate("httpclient.json.parse_failed", "", nil), err)
		}
	}

	return nil
}
