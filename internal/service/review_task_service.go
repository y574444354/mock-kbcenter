package service

import (
	"context"
	"fmt"

	"github.com/zgsm/go-webserver/i18n"
	"github.com/zgsm/go-webserver/pkg/asynq"
	"github.com/zgsm/go-webserver/pkg/idgen"
	"github.com/zgsm/go-webserver/pkg/logger"
	"github.com/zgsm/go-webserver/pkg/types"
	"github.com/zgsm/go-webserver/tasks"

	"github.com/zgsm/go-webserver/internal/model"
	"github.com/zgsm/go-webserver/internal/repository"
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
	id, err := idgen.GenerateString()
	if err != nil {
		return "", fmt.Errorf("%s", i18n.Translate("review_task.generate_id_failed", "", map[string]interface{}{
			"error": err.Error(),
		}))
	}

	reviewTask := &model.ReviewTask{
		ID:        id,
		Status:    0,
		ClientID:  clientID,
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

	// Validate each target
	for _, target := range targets {
		if err := target.Validate(); err != nil {
			return "", err
		}
	}

	logger.Info(fmt.Sprintf("review task targets: %v", targets))

	// Create reviewTask
	reviewTaskID, err := s.Create(clientID, workspace, targets)
	if err != nil {
		return "", err
	}

	// Create async task
	task, err := tasks.NewRunReviewTaskPayload(tasks.RunReviewTaskPayload{
		ReviewTaskID: reviewTaskID,
	}, tasks.QueueCritical)
	if err != nil {
		return "", err
	}

	// Enqueue task and get task ID
	taskID, err := asynq.EnqueueTask(task, tasks.QueueCritical)
	if err != nil {
		return "", err
	}

	// Update reviewTask record to save task ID
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
	// Get issues
	var issues []types.Issue

	// Get total count
	progress, err := s.reviewTaskRepo.GetProgress(context.Background(), reviewTaskID, clientID)
	if err != nil {
		return nil, err
	}

	return &types.IssueIncrementReviewTaskResult{
		Progress: progress,
		Issues:   issues,
	}, nil
}
