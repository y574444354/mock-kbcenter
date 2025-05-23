package thirdPlatform

import (
	"fmt"
	"time"

	"github.com/zgsm/review-manager/config"
	"github.com/zgsm/review-manager/i18n"
	"github.com/zgsm/review-manager/pkg/httpclient"
)

// Service 定义
type Service struct {
	client *httpclient.Client
}

type HttpServices struct {
	IssueManager IssueManagerService
}

var serverManager *HttpServices

// InitHTTPClient 初始化HTTP客户端
func InitHTTPClient() error {
	issueManagerService, err := NewIssueManagerService()
	if err != nil {
		return err
	}

	serverManager = &HttpServices{
		IssueManager: *issueManagerService,
		// 添加其他服务
	}

	return nil
}

func GetServerManager() (*HttpServices, error) {
	if serverManager == nil {
		return nil, fmt.Errorf("%s", i18n.Translate("httpclient.service.not_initialized", "", nil))
	}
	return serverManager, nil
}

// GetServiceConfig 从应用配置中获取服务配置并转换为HTTP客户端配置
func GetServiceConfig(serviceName string) (*httpclient.HttpServiceConfig, error) {
	// 获取应用配置
	cfg := config.GetConfig()

	// 获取服务配置
	serviceCfg, ok := cfg.HTTPClient.Services[serviceName]
	if !ok {
		return nil, fmt.Errorf("%s", i18n.Translate("httpclient.service.config_not_found", "", map[string]interface{}{"service": serviceName}))
	}

	// 创建HTTP客户端配置
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

	// 合并默认请求头
	if clientConfig.Headers == nil {
		clientConfig.Headers = make(map[string]string)
	}
	for k, v := range cfg.HTTPClient.Headers {
		// 只有当服务配置中没有同名请求头时，才使用默认请求头
		if _, exists := clientConfig.Headers[k]; !exists {
			clientConfig.Headers[k] = v
		}
	}

	// 如果服务配置中没有指定超时时间，使用默认超时时间
	if serviceCfg.Timeout <= 0 {
		clientConfig.Timeout = time.Duration(cfg.HTTPClient.Timeout) * time.Second
	}

	// 如果服务配置中没有指定重试次数，使用默认重试次数
	if serviceCfg.MaxRetries <= 0 {
		clientConfig.MaxRetries = cfg.HTTPClient.MaxRetries
	}

	return clientConfig, nil
}
