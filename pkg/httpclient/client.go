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

	"github.com/zgsm/mock-kbcenter/i18n"
	"github.com/zgsm/mock-kbcenter/pkg/logger"
)

// Client HTTP client
type Client struct {
	config      *HttpServiceConfig
	httpClient  *http.Client
	middlewares []Middleware
}

// NewClient create new HTTP client
func NewClient(config *HttpServiceConfig) (*Client, error) {
	if config == nil {
		config = DefaultHttpServiceConfig()
	}

	// Create Transport
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
	}

	// Configure proxy
	if config.ProxyURL != "" {
		proxyURL, err := url.Parse(config.ProxyURL)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", i18n.Translate("httpclient.proxy.parse_failed", "", nil), err)
		}
		transport.Proxy = http.ProxyURL(proxyURL)
	}

	// Configure TLS
	if config.InsecureSkipVerify || config.CertFile != "" || config.CAFile != "" {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: config.InsecureSkipVerify,
		}

		// Load client certificate
		if config.CertFile != "" && config.KeyFile != "" {
			cert, err := tls.LoadX509KeyPair(config.CertFile, config.KeyFile)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", i18n.Translate("httpclient.cert.load_failed", "", nil), err)
			}
			tlsConfig.Certificates = []tls.Certificate{cert}
		}

		// Load CA certificate
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

	// Create HTTP client
	httpClient := &http.Client{
		Timeout:   config.Timeout,
		Transport: transport,
	}

	// Create client instance
	client := &Client{
		config:      config,
		httpClient:  httpClient,
		middlewares: make([]Middleware, 0),
	}

	// Add default middlewares
	client.AddMiddleware(&LogMiddleware{
		EnableRequestLog:  config.EnableRequestLog,
		EnableResponseLog: config.EnableResponseLog,
	})

	client.AddMiddleware(&HeaderMiddleware{
		Headers: config.Headers,
	})

	// Add status code validation middleware
	client.AddMiddleware(&StatusCodeMiddleware{
		ValidStatusCodes: config.ValidStatusCodes,
	})

	// Add authentication middleware
	if config.AuthType != "none" {
		client.AddMiddleware(&AuthMiddleware{
			AuthType:   config.AuthType,
			Username:   config.Username,
			Password:   config.Password,
			Token:      config.Token,
			AuthHeader: config.AuthHeader,
		})
	}

	// Add retry middleware
	if config.MaxRetries > 0 {
		client.AddMiddleware(&RetryMiddleware{
			MaxRetries: config.MaxRetries,
			RetryDelay: config.RetryDelay,
		})
	}

	return client, nil
}

// AddMiddleware add middleware
func (c *Client) AddMiddleware(middleware Middleware) {
	if middleware == nil {
		return
	}
	c.middlewares = append(c.middlewares, middleware)
}

// Request send HTTP request
func (c *Client) Request(ctx context.Context, method, path string, body interface{}, headers map[string]string) (*http.Response, error) {
	// Build full URL
	fullURL := path
	if !strings.HasPrefix(path, "http") && c.config.BaseURL != "" {
		fullURL = c.config.BaseURL + path
	}

	// Prepare request body
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
			// Default to serialize object to JSON
			jsonData, err := json.Marshal(body)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", i18n.Translate("httpclient.request.serialize_failed", "", nil), err)
			}
			reqBody = bytes.NewReader(jsonData)
			// If Content-Type not specified, set to application/json
			if headers == nil {
				headers = make(map[string]string)
			}
			if _, ok := headers["Content-Type"]; !ok {
				headers["Content-Type"] = "application/json"
			}
		}
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, method, fullURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", i18n.Translate("httpclient.request.create_failed", "", nil), err)
	}

	// Add custom headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Apply request middlewares
	for _, middleware := range c.middlewares {
		if err := middleware.ProcessRequest(req); err != nil {
			return nil, fmt.Errorf("%s: %w", i18n.Translate("httpclient.middleware.process_failed", "", nil), err)
		}
	}

	// Send request
	var resp *http.Response
	var respErr error

	// Implement retry logic
	retries := 0
	for {
		resp, respErr = c.httpClient.Do(req)

		// Apply response middlewares
		for i := len(c.middlewares) - 1; i >= 0; i-- {
			resp, respErr = c.middlewares[i].ProcessResponse(resp, respErr)
		}

		// Check if retry needed
		shouldRetry := false
		if c.config.MaxRetries > retries {
			if respErr != nil {
				// Retry on network error
				shouldRetry = true
			} else if resp.StatusCode >= 500 {
				// Retry on server error
				shouldRetry = true
			}
		}

		if !shouldRetry {
			break
		}

		// Close response body and prepare for retry
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}

		retries++
		logger.Debug(i18n.Translate("httpclient.retry.attempt", "", map[string]interface{}{
			"attempt": retries,
			"max":     c.config.MaxRetries,
		}))
		time.Sleep(c.config.RetryDelay)

		// Recreate request body
		if body != nil {
			switch v := body.(type) {
			case string:
				reqBody = strings.NewReader(v)
			case []byte:
				reqBody = bytes.NewReader(v)
			case io.Reader:
				// Cannot reset io.Reader, this is a limitation
				return nil, errors.New(i18n.Translate("httpclient.retry.reader_unsupported", "", nil))
			default:
				// Reserialize JSON
				jsonData, _ := json.Marshal(body)
				reqBody = bytes.NewReader(jsonData)
			}
			req, err = http.NewRequestWithContext(ctx, method, fullURL, reqBody)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", i18n.Translate("httpclient.retry.request_failed", "", nil), err)
			}
			// Re-add headers
			for key, value := range headers {
				req.Header.Set(key, value)
			}
			// Re-apply request middlewares
			for _, middleware := range c.middlewares {
				if err := middleware.ProcessRequest(req); err != nil {
					return nil, fmt.Errorf("%s: %w", i18n.Translate("httpclient.retry.middleware_failed", "", nil), err)
				}
			}
		}
	}

	return resp, respErr
}

// Get send GET request
func (c *Client) Get(ctx context.Context, path string, headers map[string]string) (*http.Response, error) {
	return c.Request(ctx, http.MethodGet, path, nil, headers)
}

// Post send POST request
func (c *Client) Post(ctx context.Context, path string, body interface{}, headers map[string]string) (*http.Response, error) {
	return c.Request(ctx, http.MethodPost, path, body, headers)
}

// Put send PUT request
func (c *Client) Put(ctx context.Context, path string, body interface{}, headers map[string]string) (*http.Response, error) {
	return c.Request(ctx, http.MethodPut, path, body, headers)
}

// Delete send DELETE request
func (c *Client) Delete(ctx context.Context, path string, headers map[string]string) (*http.Response, error) {
	return c.Request(ctx, http.MethodDelete, path, nil, headers)
}

// Patch send PATCH request
func (c *Client) Patch(ctx context.Context, path string, body interface{}, headers map[string]string) (*http.Response, error) {
	return c.Request(ctx, http.MethodPatch, path, body, headers)
}

// GetJSON send GET request and parse JSON response
func (c *Client) GetJSON(ctx context.Context, path string, headers map[string]string, v interface{}) error {
	resp, err := c.Get(ctx, path, headers)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return c.parseJSONResponse(resp, v)
}

// PostJSON send POST request and parse JSON response
func (c *Client) PostJSON(ctx context.Context, path string, body interface{}, headers map[string]string, v interface{}) error {
	resp, err := c.Post(ctx, path, body, headers)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return c.parseJSONResponse(resp, v)
}

// PutJSON send PUT request and parse JSON response
func (c *Client) PutJSON(ctx context.Context, path string, body interface{}, headers map[string]string, v interface{}) error {
	resp, err := c.Put(ctx, path, body, headers)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return c.parseJSONResponse(resp, v)
}

// DeleteJSON send DELETE request and parse JSON response
func (c *Client) DeleteJSON(ctx context.Context, path string, headers map[string]string, v interface{}) error {
	resp, err := c.Delete(ctx, path, headers)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return c.parseJSONResponse(resp, v)
}

// PatchJSON send PATCH request and parse JSON response
func (c *Client) PatchJSON(ctx context.Context, path string, body interface{}, headers map[string]string, v interface{}) error {
	resp, err := c.Patch(ctx, path, body, headers)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return c.parseJSONResponse(resp, v)
}

// parseJSONResponse parse JSON response
func (c *Client) parseJSONResponse(resp *http.Response, v interface{}) error {
	// Check status code
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

	// Parse JSON response
	if v != nil {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return fmt.Errorf("%s: %w", i18n.Translate("httpclient.json.parse_failed", "", nil), err)
		}
	}

	return nil
}
