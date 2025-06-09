package model

import (
	"time"

	"github.com/zgsm/go-webserver/pkg/types"
)

type ReviewTask struct {
	ID            string         `json:"id" gorm:"primaryKey"`
	Status        int            `json:"status" gorm:"default:0"`             // Task status: 0 not started, 1 in progress, 2 completed, 3 canceled
	ErrMsg        string         `json:"err_msg" gorm:"type:varchar(50)"`     // Error message
	ClientId      string         `json:"client_id" gorm:"type:varchar(50)"`   // Client identifier
	Workspace     string         `json:"workspace" gorm:"type:varchar(255)"`  // Workspace
	TotalCount    int            `json:"total_count" gorm:"default:0"`        // Total number of subtasks
	FinishedCount int            `json:"finished_count" gorm:"default:0"`     // Number of completed subtasks
	RunTaskID     string         `json:"run_task_id" gorm:"type:varchar(50)"` // Async task ID
	Targets       []types.Target `json:"targets" gorm:"type:json"`            // List of targets
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

func (ReviewTask) TableName() string {
	return "review_task"
}
