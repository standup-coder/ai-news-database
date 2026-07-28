package deep_research

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"ai-news-database/internal/article"
	"ai-news-database/internal/config"
	"ai-news-database/internal/crawler"
	"ai-news-database/internal/db"
	"ai-news-database/internal/llm"
	"ai-news-database/internal/search"
	"github.com/PuerkitoBio/goquery"
)

type Researcher struct {
	db          *db.DB
	llmClient   llm.LLMClient
	cfg         *config.LLMConfig
	searchEng   *search.Engine
	contentRead crawler.ContentReader
	cache       *ResearchCache
	rateLimiter *RateLimiter
	metrics     *Metrics
}

type ResearchConfig struct {
	Days             int
	Limit            int
	SubQueryCount    int
	WebSearchEnabled bool
	MaxSources       int
	Timeout          time.Duration
	MaxRetries       int
	CacheEnabled     bool
	FetchContent     bool
	MinCredibility   float64
}

type ResearchResult struct {
	Topic        string
	Plan         *ResearchPlan
	Trace        ResearchTrace
	Summary      string
	Findings     []Finding
	Perspectives []Perspective
	Sources      []Source
	Evidence     []Evidence
	GapAnalysis  GapAnalysis
	Confidence   ConfidenceLevel
	GeneratedAt  time.Time
	CacheHit     bool
}

type ResearchPlan struct {
	Hypotheses    []Hypothesis  `json:"hypotheses"`
	SearchQueries []SearchQuery `json:"search_queries"`
	KeyAspects    []string      `json:"key_aspects"`
	Depth         string        `json:"depth"`
	Scope         string        `json:"scope"`
}

type Hypothesis struct {
	Statement      string `json:"statement"`
	Priority       int    `json:"priority"`
	EvidenceNeeded string `json:"evidence_needed"`
}

type SearchQuery struct {
	Query      string `json:"query"`
	SourceType string `json:"source_type"` // "local", "web", "both"
	Purpose    string `json:"purpose"`
	Priority   int    `json:"priority"`
}

type ResearchTrace struct {
	Phases    []PhaseTrace  `json:"phases"`
	TotalTime time.Duration `json:"total_time"`
}

type PhaseTrace struct {
	Name       string        `json:"name"`
	StartTime  time.Time     `json:"start_time"`
	EndTime    time.Time     `json:"end_time"`
	Duration   time.Duration `json:"duration"`
	ItemsFound int           `json:"items_found"`
	Error      string        `json:"error,omitempty"`
}

type Source struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	URL         string  `json:"url"`
	SourceName  string  `json:"source_name"`
	Domain      string  `json:"domain"`
	Content     string  `json:"content,omitempty"`
	Snippet     string  `json:"snippet"`
	Credibility float64 `json:"credibility"`
	PublishedAt string  `json:"published_at"`
	Type        string  `json:"type"` // "local_article", "web_search", "fetched_content"
}

type Evidence struct {
	ID           int     `json:"id"`
	SourceIDs    []int   `json:"source_ids"`
	Claim        string  `json:"claim"`
	Quote        string  `json:"quote,omitempty"`
	Category     string  `json:"category"` // "fact", "opinion", "statistic", "trend"
	Verified     bool    `json:"verified"`
	Corroborated int     `json:"corroborated"` // number of sources supporting
	Relevance    float64 `json:"relevance"`
}

type Finding struct {
	ID          int             `json:"id"`
	Category    string          `json:"category"`
	Title       string          `json:"title"`
	Content     string          `json:"content"`
	EvidenceIDs []int           `json:"evidence_ids"`
	Confidence  ConfidenceLevel `json:"confidence"`
}

type Perspective struct {
	Viewpoint   string `json:"viewpoint"`
	Summary     string `json:"summary"`
	EvidenceIDs []int  `json:"evidence_ids"`
}

type GapAnalysis struct {
	UncoveredAspects []string        `json:"uncovered_aspects"`
	WeakEvidence     []string        `json:"weak_evidence"`
	Contradictions   []Contradiction `json:"contradictions"`
}

type Contradiction struct {
	EvidenceA  int    `json:"evidence_a_id"`
	EvidenceB  int    `json:"evidence_b_id"`
	Resolution string `json:"resolution,omitempty"`
}

type ConfidenceLevel string

const (
	ConfidenceHigh   ConfidenceLevel = "high"
	ConfidenceMedium ConfidenceLevel = "medium"
	ConfidenceLow    ConfidenceLevel = "low"
)

type ReportFormat int

const (
	ReportFormatMarkdown ReportFormat = iota
	ReportFormatJSON
	ReportFormatDetailed
)

var DefaultResearchConfig = &ResearchConfig{
	Days:             30,
	Limit:            20,
	SubQueryCount:    5,
	WebSearchEnabled: true,
	MaxSources:       50,
	Timeout:          5 * time.Minute,
	MaxRetries:       3,
	CacheEnabled:     true,
	FetchContent:     true,
	MinCredibility:   0.3,
}

func New(database *db.DB, cfg *config.LLMConfig) *Researcher {
	return NewWithDeps(database, cfg, llm.NewClient(cfg), search.NewEngine(), crawler.NewJinaReader())
}

func NewWithDeps(database *db.DB, cfg *config.LLMConfig, llmClient llm.LLMClient, searchEng *search.Engine, contentReader crawler.ContentReader) *Researcher {
	return &Researcher{
		db:          database,
		llmClient:   llmClient,
		cfg:         cfg,
		searchEng:   searchEng,
		contentRead: contentReader,
		cache:       NewResearchCache(30 * time.Minute),
		rateLimiter: NewRateLimiter(10, time.Minute),
		metrics:     &Metrics{},
	}
}

func (r *Researcher) Research(ctx context.Context, topic string, cfg *ResearchConfig) (*ResearchResult, error) {
	if cfg == nil {
		cfg = DefaultResearchConfig
	}

	ctx, cancel := context.WithTimeout(ctx, cfg.Timeout)
	defer cancel()

	r.metrics.incrementTotal()

	cacheKey := r.cacheKey(topic, cfg)
	if cfg.CacheEnabled {
		if cached := r.cache.Get(cacheKey); cached != nil {
			slog.Info("research cache hit", "topic", topic)
			r.metrics.incrementCacheHit()
			cached.CacheHit = true
			return cached, nil
		}
	}

	startTime := time.Now()
	result := &ResearchResult{
		Topic:       topic,
		GeneratedAt: startTime,
	}

	if err := r.executeResearchPhases(ctx, topic, cfg, result); err != nil {
		r.metrics.recordError(err.Error())
		return nil, err
	}

	result.Trace.TotalTime = time.Since(startTime)

	if cfg.CacheEnabled {
		r.cache.Set(cacheKey, result)
	}

	return result, nil
}

func (r *Researcher) executeResearchPhases(ctx context.Context, topic string, cfg *ResearchConfig, result *ResearchResult) error {
	slog.Info("starting deep research", "topic", topic)

	phase := func(name string) (func(int, string), time.Time) {
		start := time.Now()
		return func(items int, err string) {
			result.Trace.Phases = append(result.Trace.Phases, PhaseTrace{
				Name:       name,
				StartTime:  start,
				EndTime:    time.Now(),
				Duration:   time.Since(start),
				ItemsFound: items,
				Error:      err,
			})
		}, start
	}

	planDone, _ := phase("Planning")
	plan, err := r.planningPhase(ctx, topic, cfg)
	result.Plan = plan
	planDone(len(plan.SearchQueries), errorToString(err))

	if err != nil {
		slog.Warn("planning phase had issues", "error", err)
	}

	searchDone, _ := phase("Search")
	sources, err := r.searchPhase(ctx, plan, cfg)
	result.Sources = sources
	searchDone(len(sources), errorToString(err))

	if len(sources) == 0 {
		return fmt.Errorf("no sources found for topic: %s", topic)
	}

	var fetchDone func(int, string)
	if cfg.FetchContent {
		fetchDone, _ = phase("ContentFetching")
		sources, err = r.contentFetchingPhase(ctx, sources, cfg)
		result.Sources = sources
		fetchDone(len(sources), errorToString(err))
	}

	extractDone, _ := phase("EvidenceExtraction")
	evidence, err := r.evidenceExtractionPhase(ctx, topic, sources, cfg)
	result.Evidence = evidence
	extractDone(len(evidence), errorToString(err))

	analysisDone, _ := phase("Analysis")
	findings, perspectives, gapAnalysis, err := r.analysisPhase(ctx, topic, evidence, sources, cfg)
	result.Findings = findings
	result.Perspectives = perspectives
	result.GapAnalysis = gapAnalysis
	analysisDone(len(findings), errorToString(err))

	summaryDone, _ := phase("Synthesis")
	summary, err := r.synthesisPhase(ctx, topic, findings, perspectives, gapAnalysis, cfg)
	result.Summary = summary
	summaryDone(0, errorToString(err))

	result.Confidence = r.calculateConfidence(result.Findings, result.Evidence, result.GapAnalysis)

	return nil
}

func errorToString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func (r *Researcher) planningPhase(ctx context.Context, topic string, cfg *ResearchConfig) (*ResearchPlan, error) {
	r.metrics.incrementLLM()

	prompt := fmt.Sprintf(`你是一个专业的深度研究员。请为以下研究主题制定详细的研究计划。

主题：%s

研究计划应包含：

1. **研究假设**：2-3个关于该主题的核心假设或预期发现
2. **搜索查询**：5-8个针对不同角度的搜索查询，覆盖：
   - 技术细节和原理
   - 最新动态和产品发布
   - 行业趋势和市场分析
   - 社区反馈和讨论
   - 历史背景和发展历程
3. **关键方面**：该主题需要探索的关键维度
4. **研究深度**：浅层（表面了解）、中层（详细分析）、深层（全面透彻）

请用JSON格式输出：
{
  "hypotheses": [{"statement": "假设内容", "priority": 1-3, "evidence_needed": "需要什么证据"}],
  "search_queries": [{"query": "搜索内容", "source_type": "local|web|both", "purpose": "探索目的", "priority": 1-3}],
  "key_aspects": ["方面1", "方面2"],
  "depth": "shallow|medium|deep",
  "scope": "狭窄|中等|广泛"
}`, topic)

	response, err := r.llmClient.SimpleChat(ctx, prompt, 3000)
	if err != nil {
		return nil, err
	}

	return r.parseResearchPlan(response)
}

func (r *Researcher) parseResearchPlan(response string) (*ResearchPlan, error) {
	plan := &ResearchPlan{}

	start := strings.Index(response, "{")
	end := strings.LastIndex(response, "}")
	if start == -1 || end == -1 {
		return plan, fmt.Errorf("invalid plan format")
	}

	jsonStr := response[start : end+1]

	if val := extractJSONInt(jsonStr, "priority"); val > 0 {
		plan.Hypotheses = append(plan.Hypotheses, Hypothesis{
			Statement:      "主假设",
			Priority:       1,
			EvidenceNeeded: "待验证",
		})
	}

	hypothesesStr := extractJSONSection(jsonStr, "hypotheses")
	if hypothesesStr != "" {
		hypothesesStr = strings.TrimPrefix(hypothesesStr, "[")
		hypothesesStr = strings.TrimSuffix(hypothesesStr, "]")
		entries := strings.Split(hypothesesStr, "},")
		for _, entry := range entries {
			if !strings.HasSuffix(entry, "}") {
				entry += "}"
			}
			h := Hypothesis{
				Statement:      extractJSONString(entry, "statement"),
				Priority:       extractJSONInt(entry, "priority"),
				EvidenceNeeded: extractJSONString(entry, "evidence_needed"),
			}
			if h.Statement != "" {
				plan.Hypotheses = append(plan.Hypotheses, h)
			}
		}
	}

	queriesStr := extractJSONSection(jsonStr, "search_queries")
	if queriesStr != "" {
		queriesStr = strings.TrimPrefix(queriesStr, "[")
		queriesStr = strings.TrimSuffix(queriesStr, "]")
		entries := strings.Split(queriesStr, "},")
		for _, entry := range entries {
			if !strings.HasSuffix(entry, "}") {
				entry += "}"
			}
			q := SearchQuery{
				Query:      extractJSONString(entry, "query"),
				SourceType: extractJSONString(entry, "source_type"),
				Purpose:    extractJSONString(entry, "purpose"),
				Priority:   extractJSONInt(entry, "priority"),
			}
			if q.Query == "" {
				q.Query = extractJSONString(entry, "Query")
			}
			if q.SourceType == "" {
				q.SourceType = "both"
			}
			if q.Query != "" {
				plan.SearchQueries = append(plan.SearchQueries, q)
			}
		}
	}

	if len(plan.SearchQueries) == 0 {
		plan.SearchQueries = append(plan.SearchQueries, SearchQuery{
			Query:      plan.Hypotheses[0].Statement,
			SourceType: "both",
			Purpose:    "主查询",
			Priority:   1,
		})
	}

	plan.Depth = extractJSONString(jsonStr, "depth")
	if plan.Depth == "" {
		plan.Depth = "medium"
	}
	plan.Scope = extractJSONString(jsonStr, "scope")
	if plan.Scope == "" {
		plan.Scope = "中等"
	}

	return plan, nil
}

func (r *Researcher) searchPhase(ctx context.Context, plan *ResearchPlan, cfg *ResearchConfig) ([]Source, error) {
	var allSources []Source
	var mu sync.Mutex
	var wg sync.WaitGroup
	errChan := make(chan error, len(plan.SearchQueries)*2)
	sourceID := 1

	addSource := func(s Source) {
		mu.Lock()
		s.ID = sourceID
		sourceID++
		allSources = append(allSources, s)
		mu.Unlock()
	}

	priorityQueries := make([]SearchQuery, 0)
	for _, q := range plan.SearchQueries {
		if q.Priority <= 2 {
			priorityQueries = append(priorityQueries, q)
		}
	}
	if len(priorityQueries) == 0 {
		priorityQueries = plan.SearchQueries
	}

	for _, sq := range priorityQueries {
		if len(allSources) >= cfg.MaxSources {
			break
		}

		select {
		case <-ctx.Done():
			return allSources, ctx.Err()
		default:
		}

		if sq.SourceType == "local" || sq.SourceType == "both" {
			wg.Add(1)
			go func(query SearchQuery) {
				defer wg.Done()

				articles, err := r.searchLocalWithRetry(query.Query, cfg.Limit, cfg.MaxRetries)
				if err == nil {
					mu.Lock()
					for _, a := range articles {
						if len(allSources) >= cfg.MaxSources {
							break
						}
						addSource(Source{
							Title:       a.Title,
							URL:         a.URL,
							SourceName:  a.Source,
							Domain:      extractDomain(a.URL),
							Snippet:     a.Summary,
							Credibility: a.QualityScore / 10.0,
							PublishedAt: formatTime(a.PublishedAt),
							Type:        "local_article",
						})
						r.metrics.incrementLocalSearch()
					}
					mu.Unlock()
				} else {
					errChan <- fmt.Errorf("local search for '%s': %w", query.Query, err)
				}
			}(sq)
		}

		if sq.SourceType == "web" || sq.SourceType == "both" {
			if !r.rateLimiter.Allow() {
				slog.Warn("rate limit reached, skipping web search")
				continue
			}

			wg.Add(1)
			go func(query SearchQuery) {
				defer wg.Done()

				results, err := r.searchWebWithRetry(query.Query, cfg.MaxRetries)
				if err == nil {
					mu.Lock()
					for _, res := range results {
						if len(allSources) >= cfg.MaxSources {
							break
						}
						addSource(Source{
							Title:       res.Title,
							URL:         res.URL,
							SourceName:  extractDomain(res.URL),
							Domain:      extractDomain(res.URL),
							Snippet:     res.Snippet,
							Credibility: r.domainCredibility(extractDomain(res.URL)),
							PublishedAt: res.PublishedDate,
							Type:        "web_search",
						})
						r.metrics.incrementWebSearch()
					}
					mu.Unlock()
				} else {
					errChan <- fmt.Errorf("web search for '%s': %w", query.Query, err)
				}
			}(sq)
		}
	}

	wg.Wait()
	close(errChan)

	var errors []string
	for err := range errChan {
		errors = append(errors, err.Error())
	}
	if len(errors) > 0 {
		slog.Warn("some searches failed", "count", len(errors))
	}

	return allSources, nil
}

func (r *Researcher) contentFetchingPhase(ctx context.Context, sources []Source, cfg *ResearchConfig) ([]Source, error) {
	if r.contentRead == nil {
		return sources, nil
	}

	type result struct {
		index   int
		content string
		err     error
	}

	resultChan := make(chan result, len(sources))
	var wg sync.WaitGroup

	for i, source := range sources {
		if source.Type != "web_search" || source.Content != "" {
			continue
		}
		if source.Credibility < cfg.MinCredibility {
			continue
		}

		wg.Add(1)
		go func(idx int, s Source) {
			defer wg.Done()

			content, err := r.contentRead.Fetch(s.URL)
			resultChan <- result{idx, content, err}
		}(i, source)

		if len(resultChan)%5 == 0 {
			wg.Wait()
		}
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for res := range resultChan {
		if res.err == nil && res.content != "" {
			sources[res.index].Content = truncate(res.content, 5000)
			sources[res.index].Type = "fetched_content"
		}
	}

	return sources, nil
}

func (r *Researcher) evidenceExtractionPhase(ctx context.Context, topic string, sources []Source, cfg *ResearchConfig) ([]Evidence, error) {
	r.metrics.incrementLLM()

	var sourceTexts []string
	for i, s := range sources {
		content := s.Content
		if content == "" {
			content = s.Snippet
		}
		if content == "" {
			continue
		}
		sourceTexts = append(sourceTexts, fmt.Sprintf(`[来源%d] %s
URL: %s
内容: %s`, i+1, s.Title, s.URL, content))
	}

	if len(sourceTexts) == 0 {
		return nil, nil
	}

	prompt := fmt.Sprintf(`你是一个专业的证据提取专家。请从以下研究来源中提取关键证据、声明和事实。

研究主题：%s

=== 来源材料 ===
%s

请提取以下类型的证据：
1. **事实 (fact)**：可验证的客观信息
2. **观点 (opinion)**：主观看法和分析
3. **统计数据 (statistic)**：具体数字和数据
4. **趋势 (trend)**：发展方向和演变

每个证据应包含：
- claim: 证据/声明内容
- quote: 原文引用（如果有）
- category: 证据类型
- relevance: 与主题的相关度 (0-1)

输出JSON数组格式：
[
  {"claim": "证据内容", "quote": "原文引用", "category": "fact|opinion|statistic|trend", "relevance": 0.9},
  ...
]

请提取10-20个最重要的证据。`, topic, strings.Join(sourceTexts, "\n\n---\n\n"))

	response, err := r.llmClient.SimpleChat(ctx, prompt, 4000)
	if err != nil {
		return nil, err
	}

	return r.parseEvidence(response, len(sources))
}

func (r *Researcher) parseEvidence(response string, sourceCount int) ([]Evidence, error) {
	var evidence []Evidence

	start := strings.Index(response, "[")
	end := strings.LastIndex(response, "]")
	if start == -1 || end == -1 {
		return evidence, nil
	}

	jsonStr := response[start : end+1]
	jsonStr = strings.TrimPrefix(jsonStr, "[")
	jsonStr = strings.TrimSuffix(jsonStr, "]")

	entries := strings.Split(jsonStr, "},")
	for i, entry := range entries {
		if !strings.HasSuffix(entry, "}") {
			entry += "}"
		}

		e := Evidence{
			ID:        i + 1,
			Claim:     extractJSONString(entry, "claim"),
			Quote:     extractJSONString(entry, "quote"),
			Category:  extractJSONString(entry, "category"),
			Relevance: extractJSONFloat(entry, "relevance"),
		}

		if e.Category == "" {
			e.Category = "fact"
		}
		if e.Relevance == 0 {
			e.Relevance = 0.5
		}

		e.SourceIDs = r.findSourceIDs(e.Claim, sourceCount)

		if e.Claim != "" {
			evidence = append(evidence, e)
		}
	}

	return evidence, nil
}

func (r *Researcher) findSourceIDs(claim string, sourceCount int) []int {
	var ids []int
	for i := 0; i < sourceCount; i++ {
		ids = append(ids, i+1)
	}
	if len(ids) > 5 {
		ids = ids[:5]
	}
	return ids
}

func (r *Researcher) analysisPhase(ctx context.Context, topic string, evidence []Evidence, sources []Source, cfg *ResearchConfig) ([]Finding, []Perspective, GapAnalysis, error) {
	r.metrics.incrementLLM()

	var gapAnalysis GapAnalysis

	verifiedCount := 0
	for i := range evidence {
		corroborated := 0
		for j := range evidence {
			if i != j && r.claimsRelated(evidence[i].Claim, evidence[j].Claim) {
				corroborated++
			}
		}
		evidence[i].Corroborated = corroborated
		evidence[i].Verified = corroborated >= 1
		if evidence[i].Verified {
			verifiedCount++
		}
	}

	findings := r.extractFindings(evidence, sources)

	perspectives := r.extractPerspectives(evidence, sources)

	topicWords := strings.Fields(topic)
	var uncoveredAspects []string
	if len(evidence) < 5 {
		uncoveredAspects = append(uncoveredAspects, "需要更多证据支撑")
	}
	for _, word := range topicWords {
		found := false
		for _, e := range evidence {
			if strings.Contains(strings.ToLower(e.Claim), strings.ToLower(word)) {
				found = true
				break
			}
		}
		if !found {
			uncoveredAspects = append(uncoveredAspects, word)
		}
	}
	gapAnalysis.UncoveredAspects = uncoveredAspects

	return findings, perspectives, gapAnalysis, nil
}

func (r *Researcher) claimsRelated(a, b string) bool {
	aWords := strings.Fields(strings.ToLower(a))
	bWords := strings.Fields(strings.ToLower(b))

	matchCount := 0
	for _, aw := range aWords {
		for _, bw := range bWords {
			if len(aw) > 3 && len(bw) > 3 && aw == bw {
				matchCount++
				break
			}
		}
	}

	return matchCount >= 2
}

func (r *Researcher) extractFindings(evidence []Evidence, sources []Source) []Finding {
	var findings []Finding

	categories := map[string][]Evidence{}
	for _, e := range evidence {
		if e.Relevance < 0.3 {
			continue
		}
		categories[e.Category] = append(categories[e.Category], e)
	}

	categoryNames := map[string]string{
		"fact":      "技术事实",
		"opinion":   "专家观点",
		"statistic": "数据统计",
		"trend":     "趋势分析",
	}

	categoryEmojis := map[string]string{
		"fact":      "📊",
		"opinion":   "💭",
		"statistic": "📈",
		"trend":     "📉",
	}

	id := 1
	for cat, catEvidence := range categories {
		if len(catEvidence) < 1 {
			continue
		}

		var claims []string
		var evidenceIDs []int
		for _, e := range catEvidence {
			claims = append(claims, e.Claim)
			evidenceIDs = append(evidenceIDs, e.ID)
		}

		content := strings.Join(claims, "；")
		if len(content) > 500 {
			content = content[:500] + "..."
		}

		confidence := ConfidenceMedium
		verifiedCount := 0
		for _, e := range catEvidence {
			if e.Verified {
				verifiedCount++
			}
		}
		if verifiedCount >= len(catEvidence)/2 {
			confidence = ConfidenceHigh
		} else if verifiedCount == 0 {
			confidence = ConfidenceLow
		}

		findings = append(findings, Finding{
			ID:          id,
			Category:    categoryNames[cat],
			Title:       fmt.Sprintf("%s发现", categoryEmojis[cat]),
			Content:     content,
			EvidenceIDs: evidenceIDs,
			Confidence:  confidence,
		})
		id++
	}

	return findings
}

func (r *Researcher) extractPerspectives(evidence []Evidence, sources []Source) []Perspective {
	var perspectives []Perspective

	verifiedEvidence := []Evidence{}
	unverifiedEvidence := []Evidence{}
	for _, e := range evidence {
		if e.Verified {
			verifiedEvidence = append(verifiedEvidence, e)
		} else {
			unverifiedEvidence = append(unverifiedEvidence, e)
		}
	}

	if len(verifiedEvidence) > 0 {
		var claims []string
		var ids []int
		for _, e := range verifiedEvidence {
			claims = append(claims, e.Claim)
			ids = append(ids, e.ID)
		}
		perspectives = append(perspectives, Perspective{
			Viewpoint:   "主流观点",
			Summary:     strings.Join(claims, "；"),
			EvidenceIDs: ids,
		})
	}

	if len(unverifiedEvidence) > 0 {
		var claims []string
		var ids []int
		for _, e := range unverifiedEvidence {
			claims = append(claims, e.Claim)
			ids = append(ids, e.ID)
		}
		perspectives = append(perspectives, Perspective{
			Viewpoint:   "待验证观点",
			Summary:     strings.Join(claims, "；"),
			EvidenceIDs: ids,
		})
	}

	return perspectives
}

func (r *Researcher) synthesisPhase(ctx context.Context, topic string, findings []Finding, perspectives []Perspective, gapAnalysis GapAnalysis, cfg *ResearchConfig) (string, error) {
	r.metrics.incrementLLM()

	var findingsSummary string
	for _, f := range findings {
		findingsSummary += fmt.Sprintf("- [%s] %s\n", f.Category, f.Content) + "\n"
	}

	var perspectivesSummary string
	for _, p := range perspectives {
		perspectivesSummary += fmt.Sprintf("- [%s] %s\n", p.Viewpoint, p.Summary) + "\n"
	}

	prompt := fmt.Sprintf(`你是一位专业的技术分析师。请根据以下研究材料，为主题「%s」生成一份执行摘要。

研究材料：

关键发现：
%s

多角度分析：
%s

信息缺口：
- 未覆盖方面：%s
- 弱证据：%s

请用2-3句话生成执行摘要，概括：
1. 该主题的核心趋势
2. 最主要的发现
3. 需要进一步关注的方向

请用中文回答，直接输出摘要内容，不要额外格式。`, topic, findingsSummary, perspectivesSummary, strings.Join(gapAnalysis.UncoveredAspects, "、"), strings.Join(gapAnalysis.WeakEvidence, "、"))

	return r.llmClient.SimpleChat(ctx, prompt, 1500)
}

func (r *Researcher) calculateConfidence(findings []Finding, evidence []Evidence, gap GapAnalysis) ConfidenceLevel {
	if len(evidence) == 0 {
		return ConfidenceLow
	}

	verifiedCount := 0
	for _, e := range evidence {
		if e.Verified {
			verifiedCount++
		}
	}
	verificationRate := float64(verifiedCount) / float64(len(evidence))

	highConfidenceFindings := 0
	for _, f := range findings {
		if f.Confidence == ConfidenceHigh {
			highConfidenceFindings++
		}
	}
	findingScore := float64(highConfidenceFindings) / float64(len(findings)+1)

	gapScore := 1.0 - float64(len(gap.UncoveredAspects))/10.0
	if gapScore < 0 {
		gapScore = 0
	}

	confidenceScore := (verificationRate*0.4 + findingScore*0.4 + gapScore*0.2)

	if confidenceScore > 0.7 {
		return ConfidenceHigh
	} else if confidenceScore > 0.4 {
		return ConfidenceMedium
	}
	return ConfidenceLow
}

func (r *Researcher) ResearchWithFormat(ctx context.Context, topic string, cfg *ResearchConfig, format ReportFormat) (string, error) {
	result, err := r.Research(ctx, topic, cfg)
	if err != nil {
		return "", err
	}

	switch format {
	case ReportFormatJSON:
		return r.toJSON(result)
	case ReportFormatDetailed:
		return r.toDetailedMarkdown(result), nil
	default:
		return r.toMarkdown(result), nil
	}
}

func (r *Researcher) FormatResult(result *ResearchResult, format ReportFormat) string {
	switch format {
	case ReportFormatJSON:
		out, _ := r.toJSON(result)
		return out
	case ReportFormatDetailed:
		return r.toDetailedMarkdown(result)
	default:
		return r.toMarkdown(result)
	}
}

func (r *Researcher) toMarkdown(result *ResearchResult) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# 🔬 深度研究报告：%s\n\n", result.Topic))
	sb.WriteString(fmt.Sprintf("> 生成时间: %s | 来源: %d | 证据: %d | 置信度: %s\n\n",
		result.GeneratedAt.Format("2006-01-02 15:04"), len(result.Sources), len(result.Evidence), result.Confidence))
	sb.WriteString("---\n\n")

	sb.WriteString("## 📋 执行摘要\n\n")
	sb.WriteString(result.Summary + "\n\n")

	if len(result.Findings) > 0 {
		sb.WriteString("---\n\n")
		sb.WriteString("## 🔍 关键发现\n\n")
		for _, f := range result.Findings {
			emoji := getCategoryEmoji(f.Category)
			sb.WriteString(fmt.Sprintf("### %s %s\n\n", emoji, f.Title))
			sb.WriteString(fmt.Sprintf("**置信度**: %s\n\n", f.Confidence))
			sb.WriteString(f.Content + "\n\n")

			if len(f.EvidenceIDs) > 0 {
				sb.WriteString("**相关证据**: ")
				for i, eid := range f.EvidenceIDs {
					if i > 0 {
						sb.WriteString(", ")
					}
					sb.WriteString(fmt.Sprintf("[%d]", eid))
				}
				sb.WriteString("\n\n")
			}
		}
	}

	if len(result.Perspectives) > 0 {
		sb.WriteString("---\n\n")
		sb.WriteString("## 🌐 多角度分析\n\n")
		for _, p := range result.Perspectives {
			sb.WriteString(fmt.Sprintf("### %s\n\n%s\n\n", p.Viewpoint, p.Summary))
		}
	}

	if len(result.GapAnalysis.UncoveredAspects) > 0 {
		sb.WriteString("---\n\n")
		sb.WriteString("## ⚠️ 信息缺口\n\n")
		sb.WriteString("以下方面信息不足：\n")
		for _, aspect := range result.GapAnalysis.UncoveredAspects {
			sb.WriteString(fmt.Sprintf("- %s\n", aspect))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("---\n\n")
	sb.WriteString("## 📚 参考来源\n\n")
	for _, s := range result.Sources {
		cred := fmt.Sprintf("%.0f%%", s.Credibility*100)
		typeLabel := map[string]string{
			"local_article":   "本地",
			"web_search":      "网络",
			"fetched_content": "全文",
		}[s.Type]
		if typeLabel == "" {
			typeLabel = s.Type
		}
		sb.WriteString(fmt.Sprintf("%d. **[%s](%s)** - %s | 可信度: %s | 类型: %s\n",
			s.ID, s.Title, s.URL, s.SourceName, cred, typeLabel))
		if s.Snippet != "" {
			sb.WriteString(fmt.Sprintf("   > %s\n", truncate(s.Snippet, 150)))
		}
	}

	sb.WriteString("\n---\n\n")
	sb.WriteString(fmt.Sprintf("*由 AI News Database Deep Research 生成 | 研究用时: %v*\n", result.Trace.TotalTime))

	return sb.String()
}

func (r *Researcher) toDetailedMarkdown(result *ResearchResult) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# 🔬 深度研究报告：%s\n\n", result.Topic))
	sb.WriteString(fmt.Sprintf("> 生成时间: %s | 来源: %d | 证据: %d | 置信度: %s\n\n",
		result.GeneratedAt.Format("2006-01-02 15:04"), len(result.Sources), len(result.Evidence), result.Confidence))

	sb.WriteString("---\n\n")
	sb.WriteString("## 🔎 研究过程\n\n")
	for _, phase := range result.Trace.Phases {
		errStr := ""
		if phase.Error != "" {
			errStr = fmt.Sprintf(" ⚠️ %s", phase.Error)
		}
		sb.WriteString(fmt.Sprintf("- **[%s]** %v | 找到 %d 项%s\n",
			phase.Name, phase.Duration, phase.ItemsFound, errStr))
	}

	if result.Plan != nil && len(result.Plan.Hypotheses) > 0 {
		sb.WriteString("\n### 研究假设\n\n")
		for _, h := range result.Plan.Hypotheses {
			sb.WriteString(fmt.Sprintf("- %s (优先级: %d)\n", h.Statement, h.Priority))
		}
	}

	sb.WriteString("\n---\n\n")
	sb.WriteString("## 📋 执行摘要\n\n")
	sb.WriteString(result.Summary + "\n\n")

	if len(result.Evidence) > 0 {
		sb.WriteString("---\n\n")
		sb.WriteString("## 📊 证据摘要\n\n")
		for _, e := range result.Evidence {
			verified := "✓"
			if !e.Verified {
				verified = "✗"
			}
			corroborated := ""
			if e.Corroborated > 0 {
				corroborated = fmt.Sprintf(" (被 %d 个来源证实)", e.Corroborated)
			}
			sb.WriteString(fmt.Sprintf("%d. [%s] %s %s%s\n", e.ID, verified, e.Category, e.Claim, corroborated))
			if e.Quote != "" {
				sb.WriteString(fmt.Sprintf("   > \"%s\"\n", truncate(e.Quote, 200)))
			}
		}
	}

	if len(result.Findings) > 0 {
		sb.WriteString("\n---\n\n")
		sb.WriteString("## 🔍 关键发现\n\n")
		for _, f := range result.Findings {
			emoji := getCategoryEmoji(f.Category)
			sb.WriteString(fmt.Sprintf("### %s %s\n\n", emoji, f.Title))
			sb.WriteString(fmt.Sprintf("**分类**: %s | **置信度**: %s\n\n", f.Category, f.Confidence))
			sb.WriteString(f.Content + "\n\n")

			if len(f.EvidenceIDs) > 0 {
				sb.WriteString("**相关证据**: ")
				for i, eid := range f.EvidenceIDs {
					if i > 0 {
						sb.WriteString(", ")
					}
					sb.WriteString(fmt.Sprintf("[%d]", eid))
				}
				sb.WriteString("\n\n")
			}
		}
	}

	if len(result.Perspectives) > 0 {
		sb.WriteString("---\n\n")
		sb.WriteString("## 🌐 多角度分析\n\n")
		for _, p := range result.Perspectives {
			sb.WriteString(fmt.Sprintf("### %s\n\n%s\n\n", p.Viewpoint, p.Summary))
			if len(p.EvidenceIDs) > 0 {
				sb.WriteString("来源证据: ")
				for _, eid := range p.EvidenceIDs {
					sb.WriteString(fmt.Sprintf("[%d] ", eid))
				}
				sb.WriteString("\n\n")
			}
		}
	}

	if len(result.GapAnalysis.UncoveredAspects) > 0 || len(result.GapAnalysis.WeakEvidence) > 0 {
		sb.WriteString("---\n\n")
		sb.WriteString("## ⚠️ 信息缺口分析\n\n")
		if len(result.GapAnalysis.UncoveredAspects) > 0 {
			sb.WriteString("**未覆盖方面**：\n")
			for _, aspect := range result.GapAnalysis.UncoveredAspects {
				sb.WriteString(fmt.Sprintf("- %s\n", aspect))
			}
			sb.WriteString("\n")
		}
		if len(result.GapAnalysis.WeakEvidence) > 0 {
			sb.WriteString("**弱证据**：\n")
			for _, weak := range result.GapAnalysis.WeakEvidence {
				sb.WriteString(fmt.Sprintf("- %s\n", weak))
			}
		}
	}

	sb.WriteString("---\n\n")
	sb.WriteString("## 📚 参考来源\n\n")
	for _, s := range result.Sources {
		cred := fmt.Sprintf("%.0f%%", s.Credibility*100)
		typeLabel := map[string]string{
			"local_article":   "本地",
			"web_search":      "网络",
			"fetched_content": "全文",
		}[s.Type]
		if typeLabel == "" {
			typeLabel = s.Type
		}
		sb.WriteString(fmt.Sprintf("%d. **[%s](%s)** - %s | 可信度: %s | 类型: %s\n",
			s.ID, s.Title, s.URL, s.SourceName, cred, typeLabel))
		if s.Snippet != "" {
			sb.WriteString(fmt.Sprintf("   > %s\n", truncate(s.Snippet, 150)))
		}
	}

	sb.WriteString("\n---\n\n")
	sb.WriteString(fmt.Sprintf("*由 AI News Database Deep Research 生成 | 总用时: %v*\n", result.Trace.TotalTime))

	return sb.String()
}

func (r *Researcher) toJSON(result *ResearchResult) (string, error) {
	data := map[string]interface{}{
		"topic":        result.Topic,
		"generated_at": result.GeneratedAt.Format(time.RFC3339),
		"confidence":   result.Confidence,
		"cache_hit":    result.CacheHit,
		"summary":      result.Summary,
		"findings":     result.Findings,
		"perspectives": result.Perspectives,
		"evidence":     result.Evidence,
		"sources":      result.Sources,
		"gap_analysis": result.GapAnalysis,
		"plan":         result.Plan,
		"trace": map[string]interface{}{
			"total_time": result.Trace.TotalTime.String(),
			"phases":     result.Trace.Phases,
		},
	}

	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

func (r *Researcher) cacheKey(topic string, cfg *ResearchConfig) string {
	data := fmt.Sprintf("%s:%d:%d:%d:%v:%v", topic, cfg.Days, cfg.Limit, cfg.SubQueryCount, cfg.WebSearchEnabled, cfg.FetchContent)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func (r *Researcher) domainCredibility(domain string) float64 {
	credibility := map[string]float64{
		"github.com":           0.9,
		"arxiv.org":            0.95,
		"stackoverflow.com":    0.8,
		"medium.com":           0.6,
		"dev.to":               0.6,
		"techcrunch.com":       0.75,
		"theverge.com":         0.7,
		"wired.com":            0.75,
		"hackernews.com":       0.8,
		"news.ycombinator.com": 0.8,
		"reddit.com":           0.5,
		"twitter.com":          0.4,
		"x.com":                0.4,
		"wikipedia.org":        0.7,
		"docs.google.com":      0.6,
	}
	if cred, ok := credibility[domain]; ok {
		return cred
	}
	if strings.HasSuffix(domain, ".gov") {
		return 0.8
	}
	if strings.HasSuffix(domain, ".edu") {
		return 0.85
	}
	return 0.5
}

func (r *Researcher) searchLocalWithRetry(query string, limit, maxRetries int) ([]article.Article, error) {
	var lastErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		articles, err := r.searchLocal(query, limit)
		if err == nil {
			return articles, nil
		}
		lastErr = err
		if attempt < maxRetries-1 {
			time.Sleep(time.Duration(attempt+1) * 100 * time.Millisecond)
		}
	}
	return nil, lastErr
}

func (r *Researcher) searchLocal(query string, limit int) ([]article.Article, error) {
	articles, err := r.db.SearchArticles(query, limit)
	if err != nil {
		articles, err = r.db.SearchByKeyword(query, limit)
		if err != nil {
			return nil, err
		}
	}
	return articles, nil
}

func (r *Researcher) searchWebWithRetry(query string, maxRetries int) ([]search.SearchResult, error) {
	var lastErr error
	backoff := time.Second

	for attempt := 0; attempt < maxRetries; attempt++ {
		results, err := r.searchWeb(query)
		if err == nil {
			return results, nil
		}
		lastErr = err
		if attempt < maxRetries-1 {
			select {
			case <-time.After(backoff):
				backoff *= 2
				if backoff > 30*time.Second {
					backoff = 30 * time.Second
				}
			}
		}
	}
	return nil, lastErr
}

func (r *Researcher) searchWeb(query string) ([]search.SearchResult, error) {
	searchURL := fmt.Sprintf("https://duckduckgo.com/html/?q=%s", strings.ReplaceAll(query, " ", "+"))
	resp, err := r.searchWebRequest(searchURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("search failed: %d", resp.StatusCode)
	}

	return r.parseSearchResults(resp.Body)
}

func (r *Researcher) searchWebRequest(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")
	req.Header.Set("Accept", "text/html")

	client := &http.Client{Timeout: 15 * time.Second}
	return client.Do(req)
}

func (r *Researcher) parseSearchResults(body io.Reader) ([]search.SearchResult, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}

	var results []search.SearchResult
	index := 1

	doc.Find(".result").Each(func(i int, s *goquery.Selection) {
		if index > 10 {
			return
		}

		result := search.SearchResult{Index: index}

		titleElem := s.Find(".result__a")
		if titleElem.Length() > 0 {
			result.Title = strings.TrimSpace(titleElem.Text())
			if href, exists := titleElem.Attr("href"); exists {
				if strings.HasPrefix(href, "http") {
					result.URL = href
				}
			}
		}

		snippetElem := s.Find(".result__snippet")
		if snippetElem.Length() > 0 {
			result.Snippet = strings.TrimSpace(snippetElem.Text())
		}

		if result.Title != "" && result.URL != "" {
			results = append(results, result)
			index++
		}
	})

	return results, nil
}

func (r *Researcher) GetMetrics() *Metrics {
	return r.metrics
}

func (r *Researcher) ResetMetrics() {
	r.metrics = &Metrics{}
}

type ResearchCache struct {
	mu    sync.RWMutex
	items map[string]*CachedResult
	ttl   time.Duration
}

type CachedResult struct {
	Result    *ResearchResult
	CreatedAt time.Time
	ExpiresAt time.Time
}

type RateLimiter struct {
	mu             sync.Mutex
	requests       []time.Time
	maxRequests    int
	windowDuration time.Duration
}

type Metrics struct {
	mu                 sync.Mutex
	TotalResearchCount int
	CacheHits          int
	CacheMisses        int
	WebSearchCount     int
	LocalSearchCount   int
	LLMCallCount       int
	TotalDuration      time.Duration
	Errors             []string
}

func NewResearchCache(ttl time.Duration) *ResearchCache {
	return &ResearchCache{
		items: make(map[string]*CachedResult),
		ttl:   ttl,
	}
}

func (c *ResearchCache) Get(key string) *ResearchResult {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if item, ok := c.items[key]; ok {
		if time.Now().Before(item.ExpiresAt) {
			return item.Result
		}
	}
	return nil
}

func (c *ResearchCache) Set(key string, result *ResearchResult) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = &CachedResult{
		Result:    result,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(c.ttl),
	}
}

func NewRateLimiter(maxRequests int, windowDuration time.Duration) *RateLimiter {
	return &RateLimiter{
		maxRequests:    maxRequests,
		windowDuration: windowDuration,
		requests:       make([]time.Time, 0),
	}
}

func (r *RateLimiter) Allow() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-r.windowDuration)

	var validRequests []time.Time
	for _, t := range r.requests {
		if t.After(cutoff) {
			validRequests = append(validRequests, t)
		}
	}
	r.requests = validRequests

	if len(r.requests) < r.maxRequests {
		r.requests = append(r.requests, now)
		return true
	}
	return false
}

func (m *Metrics) incrementTotal() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.TotalResearchCount++
}

func (m *Metrics) incrementCacheHit() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.CacheHits++
}

func (m *Metrics) incrementWebSearch() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.WebSearchCount++
}

func (m *Metrics) incrementLocalSearch() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.LocalSearchCount++
}

func (m *Metrics) incrementLLM() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.LLMCallCount++
}

func (m *Metrics) recordError(err string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Errors = append(m.Errors, err)
	if len(m.Errors) > 100 {
		m.Errors = m.Errors[len(m.Errors)-100:]
	}
}

func getCategoryEmoji(category string) string {
	switch category {
	case "技术事实":
		return "📊"
	case "专家观点":
		return "💭"
	case "数据统计":
		return "📈"
	case "趋势分析":
		return "📉"
	case "技术趋势":
		return "📈"
	case "产品动态":
		return "🚀"
	case "行业影响":
		return "🌐"
	case "安全警示":
		return "⚠️"
	case "社区热点":
		return "🔥"
	default:
		return "📌"
	}
}

func extractDomain(urlStr string) string {
	parts := strings.Split(urlStr, "/")
	if len(parts) >= 3 {
		host := parts[2]
		host = strings.TrimPrefix(host, "www.")
		return host
	}
	return urlStr
}

func formatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02")
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

func extractJSONString(s, key string) string {
	pattern := fmt.Sprintf(`"%s"`, key)
	idx := strings.Index(s, pattern)
	if idx == -1 {
		return ""
	}
	start := idx + len(pattern) + 1
	for start < len(s) && (s[start] == ' ' || s[start] == ':') {
		if s[start] == '"' {
			start++
			break
		}
		start++
	}
	end := start
	for end < len(s) && s[end] != '"' {
		if s[end] == '\\' && end+1 < len(s) {
			end += 2
			continue
		}
		end++
	}
	result := s[start:end]
	result = strings.ReplaceAll(result, `\"`, `"`)
	result = strings.ReplaceAll(result, `\\`, `\`)
	return strings.TrimSpace(result)
}

func extractJSONInt(s, key string) int {
	pattern := fmt.Sprintf(`"%s"`, key)
	idx := strings.Index(s, pattern)
	if idx == -1 {
		return 0
	}
	start := idx + len(pattern)
	for start < len(s) && (s[start] == ' ' || s[start] == ':' || s[start] == '"') {
		start++
	}
	end := start
	for end < len(s) && s[end] >= '0' && s[end] <= '9' {
		end++
	}
	result := s[start:end]
	if result == "" {
		return 0
	}
	var num int
	fmt.Sscanf(result, "%d", &num)
	return num
}

func extractJSONFloat(s, key string) float64 {
	pattern := fmt.Sprintf(`"%s"`, key)
	idx := strings.Index(s, pattern)
	if idx == -1 {
		return 0
	}
	start := idx + len(pattern)
	for start < len(s) && (s[start] == ' ' || s[start] == ':' || s[start] == '"') {
		start++
	}
	end := start
	for end < len(s) && (s[end] >= '0' && s[end] <= '9' || s[end] == '.') {
		end++
	}
	result := s[start:end]
	if result == "" {
		return 0
	}
	var num float64
	fmt.Sscanf(result, "%f", &num)
	return num
}

func extractJSONSection(s, key string) string {
	start := strings.Index(s, fmt.Sprintf(`"%s"`, key))
	if start == -1 {
		return ""
	}
	bracketStart := strings.Index(s[start:], "[")
	if bracketStart == -1 {
		return ""
	}
	bracketStart += start

	depth := 0
	for i := bracketStart; i < len(s); i++ {
		if s[i] == '[' {
			depth++
		} else if s[i] == ']' {
			depth--
			if depth == 0 {
				return s[bracketStart : i+1]
			}
		}
	}
	return ""
}
