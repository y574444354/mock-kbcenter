package service

import (
	"context"
	"errors"
	"testing"

	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/zgsm/review-manager/internal/model"
	"github.com/zgsm/review-manager/internal/repository"
	projectasynq "github.com/zgsm/review-manager/pkg/asynq"
)

// 创建一个模拟的 ReviewTaskRepository
type MockReviewTaskRepository struct {
	mock.Mock
}

func (m *MockReviewTaskRepository) Create(ctx context.Context, reviewTask *model.ReviewTask) error {
	args := m.Called(ctx, reviewTask)
	return args.Error(0)
}

func (m *MockReviewTaskRepository) Update(ctx context.Context, reviewTask *model.ReviewTask) error {
	args := m.Called(ctx, reviewTask)
	return args.Error(0)
}

func (m *MockReviewTaskRepository) GetProgress(ctx context.Context, reviewTaskID, clientID string) (float64, error) {
	args := m.Called(ctx, reviewTaskID, clientID)
	return args.Get(0).(float64), args.Error(1)
}
// 创建一个模拟的 asynq 包
type mockAsynq struct{}

// 模拟 EnqueueTask 函数
func (m *mockAsynq) EnqueueTask(task *asynq.Task, queue string) (string, error) {
	return "mock-task-id", nil
}

// 测试 Create 方法
func TestReviewTaskService_Create(t *testing.T) {
	// 创建测试用例
	tests := []struct {
		name      string
		clientID  string
		workspace string
		targets   []model.Target
		mockSetup func(*MockReviewTaskRepository)
		wantErr   bool
	}{
		{
			name:      "成功创建任务",
			clientID:  "test-client",
			workspace: "test-workspace",
			targets: []model.Target{
				{
					Type:     "file",
					FilePath: "test-file.go",
				},
			},
			mockSetup: func(mockRepo *MockReviewTaskRepository) {
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.ReviewTask")).Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "创建任务失败",
			clientID:  "test-client",
			workspace: "test-workspace",
			targets: []model.Target{
				{
					Type:     "file",
					FilePath: "test-file.go",
				},
			},
			mockSetup: func(mockRepo *MockReviewTaskRepository) {
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.ReviewTask")).Return(errors.New("创建失败"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建模拟的 repository
			mockRepo := new(MockReviewTaskRepository)
			tt.mockSetup(mockRepo)

			// 创建 service 实例，注入模拟的 repository
			s := &reviewTaskService{
				reviewTaskRepo: mockRepo,
			}

			// 调用被测试的方法
			got, err := s.Create(tt.clientID, tt.workspace, tt.targets)

			// 验证结果
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, got)
			}

			// 验证模拟对象的调用
			mockRepo.AssertExpectations(t)
		})
	}
}

// 创建一个模拟的 asynqClient 接口
type MockAsynqClient struct {
	mock.Mock
}

func (m *MockAsynqClient) Enqueue(task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	args := m.Called(task, opts)
	return &asynq.TaskInfo{ID: "mock-task-id"}, args.Error(1)
}

func (m *MockAsynqClient) Close() error {
	return nil
}

// 测试 Run 方法
func TestReviewTaskService_Run(t *testing.T) {
	// 保存原始的 EnqueueTask 函数
	originalEnqueueTask := projectasynq.EnqueueTask
	
	// 创建一个模拟的 asynq 实例
	mockAsynqInstance := &mockAsynq{}
	
	// 替换为模拟函数
	projectasynq.EnqueueTask = mockAsynqInstance.EnqueueTask
	
	// 在测试结束后恢复原始函数
	defer func() {
		projectasynq.EnqueueTask = originalEnqueueTask
	}()
	
	// 创建测试用例
	tests := []struct {
		name      string
		clientID  string
		workspace string
		targets   []model.Target
		mockSetup func(*MockReviewTaskRepository)
		wantErr   bool
	}{
		{
			name:      "成功运行任务",
			clientID:  "test-client",
			workspace: "test-workspace",
			targets: []model.Target{
				{
					Type:     "file",
					FilePath: "test-file.go",
				},
			},
			mockSetup: func(mockRepo *MockReviewTaskRepository) {
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.ReviewTask")).Return(nil)
				mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*model.ReviewTask")).Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "创建任务失败",
			clientID:  "test-client",
			workspace: "test-workspace",
			targets: []model.Target{
				{
					Type:     "file",
					FilePath: "test-file.go",
				},
			},
			mockSetup: func(mockRepo *MockReviewTaskRepository) {
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.ReviewTask")).Return(errors.New("创建失败"))
			},
			wantErr: true,
		},
		{
			name:      "更新任务失败",
			clientID:  "test-client",
			workspace: "test-workspace",
			targets: []model.Target{
				{
					Type:     "file",
					FilePath: "test-file.go",
				},
			},
			mockSetup: func(mockRepo *MockReviewTaskRepository) {
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.ReviewTask")).Return(nil)
				mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*model.ReviewTask")).Return(errors.New("更新失败"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建模拟的 repository
			mockRepo := new(MockReviewTaskRepository)
			tt.mockSetup(mockRepo)

			// 创建 service 实例，注入模拟的 repository
			s := &reviewTaskService{
				reviewTaskRepo: mockRepo,
			}

			// 调用被测试的方法
			got, err := s.Run(tt.clientID, tt.workspace, tt.targets)

			// 验证结果
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, got)
			}

			// 验证模拟对象的调用
			mockRepo.AssertExpectations(t)
		})
	}
}

// 测试 IssueIncrement 方法
func TestReviewTaskService_IssueIncrement(t *testing.T) {
	// 创建测试用例
	tests := []struct {
		name        string
		reviewTaskID string
		clientID    string
		offset      int
		mockSetup   func(*MockReviewTaskRepository)
		wantErr     bool
		wantProgress float64
	}{
		{
			name:        "成功获取进度",
			reviewTaskID: "test-task-id",
			clientID:    "test-client",
			offset:      0,
			mockSetup: func(mockRepo *MockReviewTaskRepository) {
				mockRepo.On("GetProgress", mock.Anything, "test-task-id", "test-client").Return(0.5, nil)
			},
			wantErr:     false,
			wantProgress: 0.5,
		},
		{
			name:        "获取进度失败",
			reviewTaskID: "test-task-id",
			clientID:    "test-client",
			offset:      0,
			mockSetup: func(mockRepo *MockReviewTaskRepository) {
				mockRepo.On("GetProgress", mock.Anything, "test-task-id", "test-client").Return(0.0, errors.New("获取进度失败"))
			},
			wantErr:     true,
			wantProgress: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建模拟的 repository
			mockRepo := new(MockReviewTaskRepository)
			tt.mockSetup(mockRepo)

			// 创建 service 实例，注入模拟的 repository
			s := &reviewTaskService{
				reviewTaskRepo: mockRepo,
			}

			// 调用被测试的方法
			got, err := s.IssueIncrement(tt.reviewTaskID, tt.clientID, tt.offset)

			// 验证结果
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.wantProgress, got.Progress)
				assert.Empty(t, got.Issues)
			}

			// 验证模拟对象的调用
			mockRepo.AssertExpectations(t)
		})
	}
}

// 测试 NewReviewTaskService 函数
func TestNewReviewTaskService(t *testing.T) {
	// 保存原始的 NewReviewTaskRepository 函数
	originalNewRepo := repository.NewReviewTaskRepository
	
	// 创建一个模拟的 repository
	mockRepo := new(MockReviewTaskRepository)
	
	// 替换为模拟函数
	repository.NewReviewTaskRepository = func() repository.ReviewTaskRepository {
		return mockRepo
	}
	
	// 在测试结束后恢复原始函数
	defer func() {
		repository.NewReviewTaskRepository = originalNewRepo
	}()

	// 调用被测试的函数
	service := NewReviewTaskService()

	// 验证结果
	assert.NotNil(t, service)
	assert.IsType(t, &reviewTaskService{}, service)
}