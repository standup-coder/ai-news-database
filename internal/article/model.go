package article

import "time"

// ReadStatus 阅读状态
type ReadStatus string

const (
	StatusUnread    ReadStatus = "unread"
	StatusRead      ReadStatus = "read"
	StatusStarred   ReadStatus = "starred"
	StatusArchived  ReadStatus = "archived"
	StatusDiscarded ReadStatus = "discarded"
)

// Article 表示一条抓取到的文章
type Article struct {
	ID            int64      `json:"id"`
	Title         string     `json:"title"`
	URL           string     `json:"url"`
	Source        string     `json:"source"`
	SourceAlias   string     `json:"source_alias"`
	Summary       string     `json:"summary"`        // 原始摘要/简介
	LLMSummary    string     `json:"llm_summary"`    // LLM 生成摘要
	LLMTags       string     `json:"llm_tags"`       // LLM 自动标签
	QualityScore  float64    `json:"quality_score"`  // 质量评分
	Language      string     `json:"language"`       // 语言
	RawContent    string     `json:"raw_content"`    // 原始抓取内容
	PublishedAt   *time.Time `json:"published_at"`
	FetchedAt     time.Time  `json:"fetched_at"`
	EnrichedAt    *time.Time `json:"enriched_at"`
	ReadStatus    ReadStatus `json:"read_status"`
	Tags          string     `json:"tags"` // 用户手动标签，逗号分隔
	Note          string     `json:"note"`
}

// IsValidStatus 检查阅读状态是否合法
func IsValidStatus(s ReadStatus) bool {
	switch s {
	case StatusUnread, StatusRead, StatusStarred, StatusArchived, StatusDiscarded:
		return true
	}
	return false
}
