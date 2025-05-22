package model

// import (
// 	"gorm.io/gorm"
// )

type ReviewTask struct {
}

func (ReviewTask) TableName() string {
	return "review_task"
}
