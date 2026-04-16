package rag

import (
	"fmt"
	"news4coder/internal/article"
	"news4coder/internal/config"
	"news4coder/internal/db"
	"news4coder/internal/llm"
	"strings"
)

// RAG 检索增强生成器
type RAG struct {
	db        *db.DB
	llmClient llm.LLMClient
	cfg       *config.LLMConfig
}

// New 创建 RAG 实例
func New(database *db.DB, cfg *config.LLMConfig) *RAG {
	return NewWithDeps(database, cfg, llm.NewClient(cfg))
}

// NewWithDeps 创建 RAG 实例（支持依赖注入，用于测试）
func NewWithDeps(database *db.DB, cfg *config.LLMConfig, llmClient llm.LLMClient) *RAG {
	return &RAG{
		db:        database,
		llmClient: llmClient,
		cfg:       cfg,
	}
}

// Answer 基于本地知识库回答问题
func (r *RAG) Answer(question string) (string, []SourceRef, error) {
	// 1. 召回相关文章（先用 FTS 搜索）
	articles, err := r.db.SearchArticles(question, 10)
	if err != nil {
		// FTS 搜索失败时 fallback 到简单 LIKE
		articles, err = r.fallbackSearch(question, 10)
		if err != nil {
			return "", nil, fmt.Errorf("知识库检索失败: %w", err)
		}
	}

	if len(articles) == 0 {
		return "本地知识库中没有找到与这个问题相关的文章。建议先执行 `news4coder sync` 和 `news4coder enrich` 拉取并增强内容。", nil, nil
	}

	// 2. 组装上下文
	var contextParts []string
	var refs []SourceRef
	for i, a := range articles {
		summary := a.LLMSummary
		if summary == "" {
			summary = a.Summary
		}
		ctx := fmt.Sprintf("[%d] 标题: %s\n来源: %s\n摘要: %s", i+1, a.Title, a.Source, summary)
		contextParts = append(contextParts, ctx)
		refs = append(refs, SourceRef{
			Index:  i + 1,
			Title:  a.Title,
			Source: a.Source,
			URL:    a.URL,
		})
	}

	context := strings.Join(contextParts, "\n\n")

	prompt := fmt.Sprintf(`你是一个专业的技术顾问，请根据以下从本地知识库中检索到的文章摘要，回答用户的问题。
如果知识库中的信息不足以回答问题，请明确说明。
回答时请尽量引用相关文章的编号（如 [1]、[2]），并在最后列出引用来源。

=== 知识库文章 ===
%s

=== 用户问题 ===
%s

请用中文回答。`, context, question)

	answer, err := r.llmClient.SimpleChat(prompt, r.cfg.AskMaxTokens)
	if err != nil {
		return "", nil, fmt.Errorf("LLM 回答生成失败: %w", err)
	}

	return answer, refs, nil
}

func (r *RAG) fallbackSearch(keyword string, limit int) ([]article.Article, error) {
	return r.db.SearchByKeyword(keyword, limit)
}

// SourceRef 引用来源
type SourceRef struct {
	Index  int
	Title  string
	Source string
	URL    string
}
