package httpclient

import (
	"time"
)

// HttpServiceConfig HTTP客户端配置
type HttpServiceConfig struct {
	// 基础URL，所有请求都会基于此URL
	BaseURL string `yaml:"base_url"`

	// 超时设置
	Timeout time.Duration `yaml:"timeout"`

	// 重试设置
	MaxRetries int           `yaml:"max_retries"`
	RetryDelay time.Duration `yaml:"retry_delay"`

	// 鉴权设置
	AuthType   string `yaml:"auth_type"` // 支持：none, basic, bearer, custom
	Username   string `yaml:"username"`  // 用于basic认证
	Password   string `yaml:"password"`  // 用于basic认证
	Token      string `yaml:"token"`     // 用于bearer认证
	AuthHeader string `yaml:"auth_header"`

	// 请求头设置
	Headers map[string]string `yaml:"headers"`

	// 代理设置
	ProxyURL string `yaml:"proxy_url"`

	// TLS设置
	InsecureSkipVerify bool   `yaml:"insecure_skip_verify"`
	CertFile           string `yaml:"cert_file"`
	KeyFile            string `yaml:"key_file"`
	CAFile             string `yaml:"ca_file"`

	// 日志设置
	EnableRequestLog  bool `yaml:"enable_request_log"`
	EnableResponseLog bool `yaml:"enable_response_log"`

	// 有效的状态码列表，为空则使用默认规则(2xx)
	ValidStatusCodes []int `yaml:"valid_status_codes"`
}

// DefaultHttpServiceConfig 返回默认配置
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
