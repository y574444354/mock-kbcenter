package types

import (
	"fmt"

	"github.com/zgsm/go-webserver/i18n"
)

// Issue 表示代码审查发现的问题
type Issue struct {
	IssueID    string   `json:"issue_id"`
	FilePath   string   `json:"file_path"`
	IssueCode  *string  `json:"issue_code,omitempty"`
	FixPatch   *string  `json:"fix_patch,omitempty"`
	StartLine  int      `json:"start_line"`
	EndLine    int      `json:"end_line"`
	Title      *string  `json:"title,omitempty"`
	Message    string   `json:"message"`
	IssueTypes []string `json:"issue_types"`
	Severity   string   `json:"severity"` // low | middle | high
	Status     int      `json:"status"`   // 0 | 1 | 2 | 3
	Confidence int      `json:"confidence"`
	CreatedAt  string   `json:"created_at"`
	UpdatedAt  string   `json:"updated_at"`
}

type IssueIncrementReviewTaskResult struct {
	IsDone     bool    `json:"is_done"`
	Progress   float64 `json:"progress"`
	Total      int     `json:"total"`
	NextOffset int     `json:"next_offset"`
	Issues     []Issue `json:"issues"`
}

type Target struct {
	Type      string `json:"type"`                 // file | folder | code
	FilePath  string `json:"file_path"`            // 文件路径
	LineRange []int  `json:"line_range,omitempty"` // 可选的行范围 [start, end]
}

func (t *Target) Validate() error {
	// 验证type
	if t.Type != "file" && t.Type != "folder" && t.Type != "code" {
		return fmt.Errorf("%s", i18n.Translate("review_task.invalid_target_type", "", nil))
	}
	// 验证file_path
	if t.FilePath == "" && t.Type != "folder" {
		return fmt.Errorf("%s", i18n.Translate("review_task.invalid_file_path", "", nil))
	}
	// 验证line_range
	if t.LineRange != nil && len(t.LineRange) == 2 && t.LineRange[0] > t.LineRange[1] {
		return fmt.Errorf("%s", i18n.Translate("review_task.invalid_line_range", "", nil))
	}
	return nil
}
