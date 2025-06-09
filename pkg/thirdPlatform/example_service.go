package thirdPlatform

import (
	"context"
	"fmt"

	"github.com/zgsm/go-webserver/pkg/httpclient"
	"github.com/zgsm/go-webserver/pkg/logger"
)

// ExampleUserProfile user profile
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

// ExampleUserSearchResult user search result
type ExampleUserSearchResult struct {
	Total int                  `json:"total"`
	Page  int                  `json:"page"`
	Size  int                  `json:"size"`
	Users []ExampleUserProfile `json:"users"`
}

// ExampleService example service
type ExampleService struct {
	*Service
}

// NewExampleService create example service
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

// GetUserProfile get user profile
func (s *ExampleService) GetUserProfile(ctx context.Context, userID string) (*ExampleUserProfile, error) {
	var response struct {
		Code    int                `json:"code"`
		Message string             `json:"message"`
		Data    ExampleUserProfile `json:"data"`
	}

	// Send request and parse response
	err := s.client.GetJSON(ctx, fmt.Sprintf("/users/%s/profile", userID), nil, &response)
	if err != nil {
		logger.Error("Failed to get user profile", "error", err, "user_id", userID)
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	// Check API response status
	if response.Code != 0 {
		logger.Error("API returned error", "code", response.Code, "message", response.Message, "user_id", userID)
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return &response.Data, nil
}

// UpdateUserProfile update user profile
func (s *ExampleService) UpdateUserProfile(ctx context.Context, userID string, profile *ExampleUserProfile) error {
	var response struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	// Send request and parse response
	err := s.client.PutJSON(ctx, fmt.Sprintf("/users/%s/profile", userID), profile, nil, &response)
	if err != nil {
		logger.Error("Failed to update user profile", "error", err, "user_id", userID)
		return fmt.Errorf("failed to update user profile: %w", err)
	}

	// Check API response status
	if response.Code != 0 {
		logger.Error("API returned error", "code", response.Code, "message", response.Message, "user_id", userID)
		return fmt.Errorf("API error: %s", response.Message)
	}

	return nil
}

// SearchUsers search users
func (s *ExampleService) SearchUsers(ctx context.Context, query string, page, pageSize int) (*ExampleUserSearchResult, error) {
	var response struct {
		Code    int                     `json:"code"`
		Message string                  `json:"message"`
		Data    ExampleUserSearchResult `json:"data"`
	}

	// Build query parameters
	params := map[string]string{
		"q":         query,
		"page":      fmt.Sprintf("%d", page),
		"page_size": fmt.Sprintf("%d", pageSize),
	}

	// Convert query parameters to URL query string
	queryString := ""
	for k, v := range params {
		if queryString == "" {
			queryString = "?"
		} else {
			queryString += "&"
		}
		queryString += k + "=" + v
	}

	// Send request and parse response
	err := s.client.GetJSON(ctx, "/users/search"+queryString, nil, &response)
	if err != nil {
		logger.Error("Failed to search users", "error", err, "query", query)
		return nil, fmt.Errorf("failed to search users: %w", err)
	}

	// Check API response status
	if response.Code != 0 {
		logger.Error("API returned error", "code", response.Code, "message", response.Message, "query", query)
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return &response.Data, nil
}
