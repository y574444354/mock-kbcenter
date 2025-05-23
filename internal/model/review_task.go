package model

import (
	"time"
)

type Target struct {
	Type      string `json:"type"`                                  // file | folder | code
	FilePath  string `json:"file_path"`                             // 文件路径
	LineRange []int  `json:"line_range,omitempty" gorm:"type:json"` // 可选的行范围 [start, end]
}

type ReviewTask struct {
	ID            string    `json:"id" gorm:"primaryKey"`
	Status        int       `json:"status" gorm:"default:0"`             // 任务状态: 0 未开始, 1 进行中, 2 已完成, 3 已取消
	ErrMsg        string    `json:"err_msg" gorm:"type:varchar(50)"`     // 错误信息
	ClientId      string    `json:"client_id" gorm:"type:varchar(50)"`   // 客户端标识
	Workspace     string    `json:"workspace" gorm:"type:varchar(255)"`  // 工作空间
	TotalCount    int       `json:"total_count" gorm:"default:0"`        // 子任务总数量
	FinishedCount int       `json:"finished_count" gorm:"default:0"`     // 子任务已完成数量
	RunTaskID     string    `json:"run_task_id" gorm:"type:varchar(50)"` // 异步任务ID
	Targets       []Target  `json:"targets" gorm:"type:json"`            // 目标列表
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (ReviewTask) TableName() string {
	return "review_task"
}
