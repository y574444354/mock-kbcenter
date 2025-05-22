package repository

import (
	"context"

	"github.com/zgsm/review-manager/internal/model"
	"github.com/zgsm/review-manager/pkg/db"
	"gorm.io/gorm"
)

type ReviewTaskRepository interface {
	Create(ctx context.Context, reviewTask *model.ReviewTask) error
}

type reviewTaskRepository struct {
	db *gorm.DB
}

func NewReviewTaskRepository() ReviewTaskRepository {
	return &reviewTaskRepository{
		db: db.GetDB(),
	}
}

func (r *reviewTaskRepository) Create(ctx context.Context, reviewTask *model.ReviewTask) (error) {
	return r.db.WithContext(ctx).Create(reviewTask).Error
}