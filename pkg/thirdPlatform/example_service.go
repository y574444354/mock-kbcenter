package thirdPlatform

import (
	"context"
	"fmt"

	"github.com/zgsm/go-webserver/pkg/httpclient"
	"github.com/zgsm/go-webserver/pkg/logger"
)

// ExampleUserProfile 用户资料
type ExampleUserProfile struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Bio       string `json:"bio"`
	Gender    string `json:"gender"`
	Birthday  string `json:"birthday"`
	Location  string `json:"location"`
	Website   string `json:"website"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// ExampleUserSearchResult 用户搜索结果
type ExampleUserSearchResult struct {
	Total int                  `json:"total"`
	Page  int                  `json:"page"`
	Size  int                  `json:"size"`
	Users []ExampleUserProfile `json:"users"`
}

// ExampleService 示例服务
type ExampleService struct {
	*Service
}

// NewExampleService 创建示例服务
func NewExampleService(clientConfig *httpclient.HttpServiceConfig) (*ExampleService, error) {
	client, err := httpclient.NewClient(clientConfig)
	if err != nil {
		return nil, err
	}

	service := &ExampleService{
		Service: &Service{
			client: client,
		},
	}

	return service, nil
}

// GetUserProfile 获取用户资料
func (s *ExampleService) GetUserProfile(ctx context.Context, userID string) (*ExampleUserProfile, error) {
	var response struct {
		Code    int                `json:"code"`
		Message string             `json:"message"`
		Data    ExampleUserProfile `json:"data"`
	}

	// 发送请求并解析响应
	err := s.client.GetJSON(ctx, fmt.Sprintf("/users/%s/profile", userID), nil, &response)
	if err != nil {
		logger.Error("获取用户资料失败", "error", err, "user_id", userID)
		return nil, fmt.Errorf("获取用户资料失败: %w", err)
	}

	// 检查API响应状态
	if response.Code != 0 {
		logger.Error("API返回错误", "code", response.Code, "message", response.Message, "user_id", userID)
		return nil, fmt.Errorf("API错误: %s", response.Message)
	}

	return &response.Data, nil
}

// UpdateUserProfile 更新用户资料
func (s *ExampleService) UpdateUserProfile(ctx context.Context, userID string, profile *ExampleUserProfile) error {
	var response struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	// 发送请求并解析响应
	err := s.client.PutJSON(ctx, fmt.Sprintf("/users/%s/profile", userID), profile, nil, &response)
	if err != nil {
		logger.Error("更新用户资料失败", "error", err, "user_id", userID)
		return fmt.Errorf("更新用户资料失败: %w", err)
	}

	// 检查API响应状态
	if response.Code != 0 {
		logger.Error("API返回错误", "code", response.Code, "message", response.Message, "user_id", userID)
		return fmt.Errorf("API错误: %s", response.Message)
	}

	return nil
}

// SearchUsers 搜索用户
func (s *ExampleService) SearchUsers(ctx context.Context, query string, page, pageSize int) (*ExampleUserSearchResult, error) {
	var response struct {
		Code    int                     `json:"code"`
		Message string                  `json:"message"`
		Data    ExampleUserSearchResult `json:"data"`
	}

	// 构建查询参数
	params := map[string]string{
		"q":         query,
		"page":      fmt.Sprintf("%d", page),
		"page_size": fmt.Sprintf("%d", pageSize),
	}

	// 将查询参数转换为URL查询字符串
	queryString := ""
	for k, v := range params {
		if queryString == "" {
			queryString = "?"
		} else {
			queryString += "&"
		}
		queryString += k + "=" + v
	}

	// 发送请求并解析响应
	err := s.client.GetJSON(ctx, "/users/search"+queryString, nil, &response)
	if err != nil {
		logger.Error("搜索用户失败", "error", err, "query", query)
		return nil, fmt.Errorf("搜索用户失败: %w", err)
	}

	// 检查API响应状态
	if response.Code != 0 {
		logger.Error("API返回错误", "code", response.Code, "message", response.Message, "query", query)
		return nil, fmt.Errorf("API错误: %s", response.Message)
	}

	return &response.Data, nil
}
