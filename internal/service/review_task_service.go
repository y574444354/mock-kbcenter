package service

import (
	"context"
	"fmt"

	"github.com/zgsm/review-manager/i18n"
	"github.com/zgsm/review-manager/pkg/asynq"
	"github.com/zgsm/review-manager/pkg/idgen"
	"github.com/zgsm/review-manager/pkg/logger"
	"github.com/zgsm/review-manager/pkg/types"
	"github.com/zgsm/review-manager/tasks"

	"github.com/zgsm/review-manager/internal/model"
	"github.com/zgsm/review-manager/internal/repository"
)

type ReviewTaskService interface {
	Create(clientID, workspace string, targets []types.Target) (string, error)
	Run(clientID, workspace string, targets []types.Target) (string, error)
	IssueIncrement(reviewTaskID, clientID string, offset int) (*types.IssueIncrementReviewTaskResult, error)
}

type reviewTaskService struct {
	reviewTaskRepo repository.ReviewTaskRepository
}

func NewReviewTaskService() ReviewTaskService {
	return &reviewTaskService{
		reviewTaskRepo: repository.NewReviewTaskRepository(),
	}
}

func (s *reviewTaskService) Create(clientID, workspace string, targets []types.Target) (string, error) {
	reviewTask := &model.ReviewTask{
		ID:        idgen.GenerateString(), // 使用雪花算法生成ID
		Status:    0,
		ClientId:  clientID,
		Workspace: workspace,
		Targets:   targets,
	}

	if err := s.reviewTaskRepo.Create(context.Background(), reviewTask); err != nil {
		return "", err
	}

	return reviewTask.ID, nil
}

func (s *reviewTaskService) Run(clientID, workspace string, targets []types.Target) (string, error) {
	if len(targets) == 0 {
		return "", fmt.Errorf("%s", i18n.Translate("review_task.empty_targets", "", nil))
	}

	// 验证每个target
	for _, target := range targets {
		if err := target.Validate(); err != nil {
			return "", err
		}
	}

	logger.Info(fmt.Sprintf("review task targets: %v", targets))

	// 创建reviewTask
	reviewTaskID, err := s.Create(clientID, workspace, targets)
	if err != nil {
		return "", err
	}

	// 创建异步任务
	task, err := tasks.NewRunReviewTaskPayload(tasks.RunReviewTaskPayload{
		ReviewTaskID: reviewTaskID,
	}, tasks.QueueCritical)
	if err != nil {
		return "", err
	}

	// 将任务加入队列并获取任务ID
	taskID, err := asynq.EnqueueTask(task, tasks.QueueCritical)
	if err != nil {
		return "", err
	}

	// 更新reviewTask记录保存任务ID
	reviewTask := &model.ReviewTask{
		ID:        reviewTaskID,
		RunTaskID: taskID,
	}
	if err := s.reviewTaskRepo.Update(context.Background(), reviewTask); err != nil {
		return "", err
	}

	return reviewTaskID, nil
}

func (s *reviewTaskService) IssueIncrement(reviewTaskID, clientID string, offset int) (*types.IssueIncrementReviewTaskResult, error) {
	// 获取issues
	var issues []types.Issue

	// 获取总数
	progress, err := s.reviewTaskRepo.GetProgress(context.Background(), reviewTaskID, clientID)
	if err != nil {
		return nil, err
	}

	return &types.IssueIncrementReviewTaskResult{
		Progress: progress,
		Issues:   issues,
	}, nil
}
