package httpclient

import (
	"time"
)

// HttpServiceConfig HTTP client configuration
type HttpServiceConfig struct {
	// Base URL, all requests will be based on this URL
	BaseURL string `yaml:"base_url"`

	// Timeout settings
	Timeout time.Duration `yaml:"timeout"`

	// Retry settings
	MaxRetries int           `yaml:"max_retries"`
	RetryDelay time.Duration `yaml:"retry_delay"`

	// Authentication settings
	AuthType   string `yaml:"auth_type"` // Supported: none, basic, bearer, custom
	Username   string `yaml:"username"`  // For basic auth
	Password   string `yaml:"password"`  // For basic auth
	Token      string `yaml:"token"`     // For bearer auth
	AuthHeader string `yaml:"auth_header"`

	// Request headers
	Headers map[string]string `yaml:"headers"`

	// Proxy settings
	ProxyURL string `yaml:"proxy_url"`

	// TLS settings
	InsecureSkipVerify bool   `yaml:"insecure_skip_verify"`
	CertFile           string `yaml:"cert_file"`
	KeyFile            string `yaml:"key_file"`
	CAFile             string `yaml:"ca_file"`

	// Logging settings
	EnableRequestLog  bool `yaml:"enable_request_log"`
	EnableResponseLog bool `yaml:"enable_response_log"`

	// Valid status codes, empty means use default rule (2xx)
	ValidStatusCodes []int `yaml:"valid_status_codes"`
}

// DefaultHttpServiceConfig returns default configuration
func DefaultHttpServiceConfig() *HttpServiceConfig {
	return &HttpServiceConfig{
		Timeout:           30 * time.Second,
		MaxRetries:        3,
		RetryDelay:        1 * time.Second,
		AuthType:          "none",
		Headers:           make(map[string]string),
		EnableRequestLog:  true,
		EnableResponseLog: true,
		ValidStatusCodes:  []int{},
	}
}
