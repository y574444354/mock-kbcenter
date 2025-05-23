package types

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
