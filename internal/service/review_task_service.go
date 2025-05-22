package service

import (
	"github.com/zgsm/review-manager/internal/repository"
)

type ReviewTaskService interface {

}

type reviewTaskService struct {
	reviewTaskRepo repository.ReviewTaskRepository
}

func NewReviewTaskService() ReviewTaskService {
	return &reviewTaskService{
		reviewTaskRepo: repository.NewReviewTaskRepository(),
	}
}