package thirdPlatform

import (
	"fmt"
	"time"

	"github.com/zgsm/mock-kbcenter/config"
	"github.com/zgsm/mock-kbcenter/i18n"
	"github.com/zgsm/mock-kbcenter/pkg/httpclient"
)

// Service definition
type Service struct {
	client *httpclient.Client
}

type HttpServices struct {
	IssueManager IssueManagerService
}

var serverManager *HttpServices

// InitHTTPClient initialize HTTP client
func InitHTTPClient() error {
	issueManagerService, err := NewIssueManagerService()
	if err != nil {
		return err
	}

	serverManager = &HttpServices{
		IssueManager: *issueManagerService,
		// Add other services
	}

	return nil
}

func GetServerManager() (*HttpServices, error) {
	if serverManager == nil {
		return nil, fmt.Errorf("%s", i18n.Translate("httpclient.service.not_initialized", "", nil))
	}
	return serverManager, nil
}

// GetServiceConfig get service config from app config and convert to HTTP client config
func GetServiceConfig(serviceName string) (*httpclient.HttpServiceConfig, error) {
	// Get application config
	cfg := config.GetConfig()

	// Get service config
	serviceCfg, ok := cfg.HTTPClient.Services[serviceName]
	if !ok {
		return nil, fmt.Errorf("%s", i18n.Translate("httpclient.service.config_not_found", "", map[string]interface{}{"service": serviceName}))
	}

	// Create HTTP client config
	clientConfig := &httpclient.HttpServiceConfig{
		BaseURL:            serviceCfg.BaseURL,
		Timeout:            time.Duration(serviceCfg.Timeout) * time.Second,
		MaxRetries:         serviceCfg.MaxRetries,
		RetryDelay:         time.Duration(cfg.HTTPClient.RetryDelay) * time.Second,
		AuthType:           serviceCfg.AuthType,
		Username:           serviceCfg.Username,
		Password:           serviceCfg.Password,
		Token:              serviceCfg.Token,
		AuthHeader:         serviceCfg.AuthHeader,
		Headers:            serviceCfg.Headers,
		ProxyURL:           cfg.HTTPClient.ProxyURL,
		InsecureSkipVerify: cfg.HTTPClient.InsecureSkipVerify,
		EnableRequestLog:   cfg.HTTPClient.EnableRequestLog,
		EnableResponseLog:  cfg.HTTPClient.EnableResponseLog,
	}

	// Merge default headers
	if clientConfig.Headers == nil {
		clientConfig.Headers = make(map[string]string)
	}
	for k, v := range cfg.HTTPClient.Headers {
		// Only use default header when service config doesn't have same header
		if _, exists := clientConfig.Headers[k]; !exists {
			clientConfig.Headers[k] = v
		}
	}

	// Use default timeout if not specified in service config
	if serviceCfg.Timeout <= 0 {
		clientConfig.Timeout = time.Duration(cfg.HTTPClient.Timeout) * time.Second
	}

	// Use default retry count if not specified in service config
	if serviceCfg.MaxRetries <= 0 {
		clientConfig.MaxRetries = cfg.HTTPClient.MaxRetries
	}

	return clientConfig, nil
}
