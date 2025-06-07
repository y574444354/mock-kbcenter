package repository

import (
	"context"

	"github.com/zgsm/go-webserver/internal/model"
	"github.com/zgsm/go-webserver/pkg/db"
	"gorm.io/gorm"
)

type ReviewTaskRepository interface {
	Create(ctx context.Context, reviewTask *model.ReviewTask) error
	Update(ctx context.Context, reviewTask *model.ReviewTask) error
	GetProgress(ctx context.Context, reviewTaskID, clientID string) (float64, error)
}

type reviewTaskRepository struct {
	db *gorm.DB
}

// NewReviewTaskRepositoryFunc 定义创建ReviewTaskRepository的函数类型
type NewReviewTaskRepositoryFunc func() ReviewTaskRepository

// NewReviewTaskRepository 创建ReviewTaskRepository的函数变量
var NewReviewTaskRepository NewReviewTaskRepositoryFunc = func() ReviewTaskRepository {
	return &reviewTaskRepository{
		db: db.GetDB(),
	}
}

func (r *reviewTaskRepository) Create(ctx context.Context, reviewTask *model.ReviewTask) error {
	return r.db.WithContext(ctx).Create(reviewTask).Error
}

func (r *reviewTaskRepository) Update(ctx context.Context, reviewTask *model.ReviewTask) error {
	return r.db.WithContext(ctx).Model(reviewTask).Updates(reviewTask).Error
}

func (r *reviewTaskRepository) GetProgress(ctx context.Context, reviewTaskID, clientID string) (float64, error) {
	var reviewTask model.ReviewTask
	err := r.db.WithContext(ctx).Where("id = ? AND client_id = ?", reviewTaskID, clientID).First(&reviewTask).Error
	if err != nil {
		return 0, err
	}
	if reviewTask.TotalCount == 0 {
		return 0, nil
	}
	return float64(reviewTask.FinishedCount) / float64(reviewTask.TotalCount), nil
}
