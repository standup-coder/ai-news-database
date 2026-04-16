package cmd

import (
	"encoding/json"
	"fmt"
	"news4coder/internal/article"
	"news4coder/internal/config"
	"news4coder/internal/crawler"
	"news4coder/internal/curator"
	"news4coder/internal/db"
	"news4coder/internal/enricher"
	"news4coder/internal/llm"
	"news4coder/internal/official"
	"news4coder/internal/rag"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"net/http"
)

var webPort int

type inspireItem struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	Summary     string `json:"summary"`
	PublishedAt string `json:"published_at,omitempty"`
	FetchedAt   string `json:"fetched_at"`
	ReadStatus  string `json:"read_status"`
	Points      int    `json:"points"`
}

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "启动本地 Web 服务",
	Long:  `启动一个本地 Web 服务器，提供工作台和灵感页面浏览。`,
	Example: `  news4coder web
  news4coder web --port 3000`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cyan := color.New(color.FgCyan).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()
		bold := color.New(color.Bold).SprintFunc()

		database, err := db.New()
		if err != nil {
			return fmt.Errorf("初始化数据库失败: %w", err)
		}
		defer database.Close()

		webDir := findWebDir()

		mux := http.NewServeMux()

		// API: /api/stats
		mux.HandleFunc("/api/stats", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			stats, err := database.Stats()
			if err != nil {
				json.NewEncoder(w).Encode(map[string]any{})
				return
			}
			total, unread, read := 0, 0, 0
			for _, s := range stats {
				if t, ok := s["total"].(int); ok {
					total += t
				}
				if u, ok := s["unread"].(int); ok {
					unread += u
				}
				if r, ok := s["read"].(int); ok {
					read += r
				}
			}
			articles, _ := database.GetArticles("", "", 0)
			enriched := 0
			for _, a := range articles {
				if a.EnrichedAt != nil {
					enriched++
				}
			}
			json.NewEncoder(w).Encode(map[string]any{
				"total": total, "unread": unread, "read": read, "enriched": enriched,
			})
		})

		// API: /api/sync
		mux.HandleFunc("/api/sync", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			var body struct {
				Source string `json:"source"`
			}
			json.NewDecoder(r.Body).Decode(&body)
			registry := official.GetRegistry()
			sources := registry.List()
			var saved, skipped int
			for _, source := range sources {
				if body.Source != "" && source.Alias != body.Source {
					continue
				}
				c, err := crawler.NewCrawler(source)
				if err != nil {
					continue
				}
				items, err := c.Fetch()
				if err != nil {
					continue
				}
				for _, item := range items {
					a := article.Article{
						Title:       item.Title,
						URL:         item.URL,
						Source:      item.Source,
						SourceAlias: item.SourceAlias,
						Summary:     item.RawContent,
						RawContent:  item.RawContent,
						ReadStatus:  article.StatusUnread,
					}
					if item.PublishedAt != nil {
						a.PublishedAt = item.PublishedAt
					}
					if err := database.SaveArticle(&a); err != nil {
						skipped++
					} else {
						saved++
					}
				}
			}
			json.NewEncoder(w).Encode(map[string]any{
				"message": fmt.Sprintf("新增 %d 条，更新 %d 条", saved, skipped),
				"saved":   saved, "skipped": skipped,
			})
		})

		// API: /api/enrich
		mux.HandleFunc("/api/enrich", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			cfg, _ := config.Load()
			if cfg.LLM.APIKey == "" {
				json.NewEncoder(w).Encode(map[string]string{"error": "LLM API Key 未配置"})
				return
			}
			articles, _ := database.GetUnenrichedArticles(10)
			if len(articles) == 0 {
				json.NewEncoder(w).Encode(map[string]any{"message": "没有需要增强的文章", "success": 0, "total": 0})
				return
			}
			enr := enricher.New(database, &cfg.LLM)
			success := 0
			for _, a := range articles {
				_, err := enr.EnrichArticle(&a)
				if err == nil {
					success++
				}
			}
			json.NewEncoder(w).Encode(map[string]any{
				"message": fmt.Sprintf("成功 %d / 总计 %d", success, len(articles)),
				"success": success, "total": len(articles),
			})
		})

		// API: /api/curate
		mux.HandleFunc("/api/curate", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			c := curator.New(database)
			picks, err := c.GetTopPicks(10)
			if err != nil {
				json.NewEncoder(w).Encode(map[string]any{"items": []any{}})
				return
			}
			items := make([]map[string]any, len(picks))
			for i, p := range picks {
				items[i] = map[string]any{
					"id": p.ID, "title": p.Title, "url": p.URL, "source": p.Source,
					"score": p.CuratorScore, "reason": p.Reason, "tags": p.LLMTags,
				}
			}
			json.NewEncoder(w).Encode(map[string]any{"items": items})
		})

		// API: /api/ask
		mux.HandleFunc("/api/ask", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			var body struct {
				Question string `json:"question"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				json.NewEncoder(w).Encode(map[string]string{"error": "invalid body"})
				return
			}
			cfg, _ := config.Load()
			if cfg.LLM.APIKey == "" {
				json.NewEncoder(w).Encode(map[string]string{"error": "LLM API Key 未配置"})
				return
			}
			rg := rag.New(database, &cfg.LLM)
			answer, refs, err := rg.Answer(body.Question)
			if err != nil {
				json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}
			refsOut := make([]map[string]string, len(refs))
			for i, ref := range refs {
				refsOut[i] = map[string]string{
					"index": fmt.Sprintf("%d", ref.Index),
					"title": ref.Title, "source": ref.Source, "url": ref.URL,
				}
			}
			json.NewEncoder(w).Encode(map[string]any{"answer": answer, "refs": refsOut})
		})

		// API: /api/articles
		mux.HandleFunc("/api/articles", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			filter := r.URL.Query().Get("filter")
			var status article.ReadStatus
			switch filter {
			case "unread":
				status = article.StatusUnread
			case "read":
				status = article.StatusRead
			case "starred":
				status = article.StatusStarred
			}
			articles, err := database.GetArticles(status, "", 50)
			if err != nil {
				json.NewEncoder(w).Encode(map[string]any{"items": []any{}})
				return
			}
			items := make([]map[string]any, len(articles))
			for i, a := range articles {
				items[i] = map[string]any{
					"id": a.ID, "title": a.Title, "url": a.URL,
					"source": a.Source, "source_alias": a.SourceAlias,
					"summary": a.Summary, "llm_summary": a.LLMSummary,
					"llm_tags": a.LLMTags, "quality_score": a.QualityScore,
					"read_status": string(a.ReadStatus),
					"fetched_at":  a.FetchedAt.Format("2006-01-02"),
				}
				if a.PublishedAt != nil {
					items[i]["published_at"] = a.PublishedAt.Format("2006-01-02")
				}
			}
			json.NewEncoder(w).Encode(map[string]any{"items": items})
		})

		// API: /api/articles/{id}/status
		mux.HandleFunc("/api/articles/", func(w http.ResponseWriter, r *http.Request) {
			re := regexp.MustCompile(`/api/articles/(\d+)/status`)
			matches := re.FindStringSubmatch(r.URL.Path)
			if len(matches) < 2 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if r.Method == http.MethodPost {
				var body struct {
					Status string `json:"status"`
				}
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				var st article.ReadStatus
				switch body.Status {
				case "read":
					st = article.StatusRead
				case "starred":
					st = article.StatusStarred
				case "discarded":
					st = article.StatusDiscarded
				case "archived":
					st = article.StatusArchived
				default:
					st = article.ReadStatus(body.Status)
				}
				var id int64
				fmt.Sscanf(matches[1], "%d", &id)
				if err := database.UpdateStatus(id, st); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
				return
			}
			w.WriteHeader(http.StatusMethodNotAllowed)
		})

		// API: /api/sources
		mux.HandleFunc("/api/sources", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			registry := official.GetRegistry()
			sources := registry.List()
			items := make([]map[string]string, len(sources))
			for i, s := range sources {
				items[i] = map[string]string{"name": s.Name, "alias": s.Alias, "url": s.URL}
			}
			json.NewEncoder(w).Encode(map[string]any{"sources": items})
		})

		// API: /api/inspire
		mux.HandleFunc("/api/inspire", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Header().Set("Access-Control-Allow-Origin", "*")

			sortBy := r.URL.Query().Get("sort")
			var articles []article.Article
			var err error
			if sortBy == "points" {
				articles, err = database.GetArticlesSorted("", "inspire", 0, "points")
			} else {
				articles, err = database.GetArticles("", "inspire", 0)
			}
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}
			var items []inspireItem
			for _, a := range articles {
				item := inspireItem{
					ID:         a.ID,
					Title:      a.Title,
					URL:        a.URL,
					Summary:    a.Summary,
					ReadStatus: string(a.ReadStatus),
					FetchedAt:  a.FetchedAt.Format("2006-01-02 15:04"),
					Points:     a.Points,
				}
				if a.PublishedAt != nil {
					item.PublishedAt = a.PublishedAt.Format("2006-01-02")
				}
				items = append(items, item)
			}
			if items == nil {
				items = []inspireItem{}
			}
			json.NewEncoder(w).Encode(map[string]any{"items": items, "total": len(items)})
		})

		// API: /api/inspire/read (existing)
		mux.HandleFunc("/api/inspire/read", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			var body struct {
				ID int64 `json:"id"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if err := database.UpdateStatus(body.ID, "read"); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
		})

		// API: /api/inspire/burst
		mux.HandleFunc("/api/inspire/burst", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			var body struct {
				Count int    `json:"count"`
				Focus string `json:"focus"`
				Mode  string `json:"mode"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": "invalid body"})
				return
			}
			if body.Count <= 0 {
				body.Count = 3
			}
			if body.Count > 10 {
				body.Count = 10
			}
			if body.Mode == "" {
				body.Mode = "cross-domain"
			}

			cfg, err := config.Load()
			if err != nil || cfg.LLM.APIKey == "" {
				w.WriteHeader(http.StatusServiceUnavailable)
				json.NewEncoder(w).Encode(map[string]string{"error": "LLM API Key 未配置"})
				return
			}

			articles, err := database.GetArticles("", "inspire", 30)
			if err != nil || len(articles) == 0 {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(map[string]string{"error": "暂无灵感数据，请先运行 inspire"})
				return
			}

			var products []string
			for i, a := range articles {
				entry := fmt.Sprintf("%d. %s", i+1, a.Title)
				if a.Summary != "" {
					if len(a.Summary) > 120 {
						entry += " — " + a.Summary[:117] + "..."
					} else {
						entry += " — " + a.Summary
					}
				}
				products = append(products, entry)
			}
			productsText := strings.Join(products, "\n")

			focusClause := ""
			if body.Focus != "" {
				focusClause = fmt.Sprintf("\n\n用户希望聚焦的方向是：%s。请围绕这个方向展开联想。", body.Focus)
			}

			template, ok := burstModePrompts[body.Mode]
			if !ok {
				template = burstModePrompts["cross-domain"]
			}
			prompt := fmt.Sprintf(template.Prompt, body.Count, productsText, focusClause)

			client := llm.NewClient(&cfg.LLM)
			resp, err := client.Chat([]llm.Message{
				{Role: "system", Content: template.System},
				{Role: "user", Content: prompt},
			}, 4000)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}
			clean := strings.TrimSpace(resp)
			clean = strings.TrimPrefix(clean, "```json")
			clean = strings.TrimPrefix(clean, "```")
			clean = strings.TrimSuffix(clean, "```")
			clean = strings.TrimSpace(clean)
			var ideas []burstIdea
			if err := json.Unmarshal([]byte(clean), &ideas); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"error": "LLM 返回格式异常", "raw": resp})
				return
			}

			ideasJSON, _ := json.Marshal(ideas)
			burstID, _ := database.SaveBurstResult(body.Mode, body.Focus, string(ideasJSON), len(articles))

			json.NewEncoder(w).Encode(map[string]any{
				"id": burstID, "ideas": ideas, "count": len(ideas),
				"based_on": len(articles), "focus": body.Focus,
				"mode":         body.Mode,
				"generated_at": time.Now().Format("2006-01-02 15:04"),
			})
		})

		// API: /api/inspire/mark-all-read
		mux.HandleFunc("/api/inspire/mark-all-read", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			affected, err := database.MarkAllRead("inspire")
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(map[string]any{"status": "ok", "updated": affected})
		})

		// API: /api/inspire/star
		mux.HandleFunc("/api/inspire/star", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			var body struct {
				ID int64 `json:"id"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if err := database.UpdateStatus(body.ID, "starred"); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
		})

		// API: /api/burst/history
		mux.HandleFunc("/api/burst/history", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			results, err := database.GetBurstResults(20)
			if err != nil {
				json.NewEncoder(w).Encode(map[string]any{"items": []any{}})
				return
			}
			if results == nil {
				results = []db.BurstResult{}
			}
			json.NewEncoder(w).Encode(map[string]any{"items": results})
		})

		// API: /api/burst/deep-dive
		mux.HandleFunc("/api/burst/deep-dive", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			var body struct {
				Title       string `json:"title"`
				Description string `json:"description"`
				Inspiration string `json:"inspiration"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			cfg, err := config.Load()
			if err != nil || cfg.LLM.APIKey == "" {
				w.WriteHeader(http.StatusServiceUnavailable)
				json.NewEncoder(w).Encode(map[string]string{"error": "LLM API Key 未配置"})
				return
			}

			prompt := fmt.Sprintf(`请对以下产品创意进行深入分析和扩展：

## 创意名称
%s

## 描述
%s

## 灵感来源
%s

请提供以下内容（用中文）：
1. **市场分析**：目标市场规模、竞品分析、差异化优势
2. **技术方案**：推荐的技术栈、架构概要、关键技术难点
3. **MVP 定义**：最小可行产品的核心功能（3-5个）
4. **开发计划**：一个人 4 周内的开发计划（按周分解）
5. **商业模式**：盈利模式、定价策略
6. **风险评估**：主要风险和应对方案

请以 JSON 格式返回：
{
  "market": "...",
  "tech": "...",
  "mvp": ["...", "..."],
  "timeline": ["第1周：...", "第2周：...", "第3周：...", "第4周：..."],
  "business": "...",
  "risks": "..."
}`, body.Title, body.Description, body.Inspiration)

			client := llm.NewClient(&cfg.LLM)
			resp, err := client.Chat([]llm.Message{
				{Role: "system", Content: "你是一位资深的独立开发者产品顾问，擅长将创意转化为可执行的计划。请用中文回答，返回纯 JSON。"},
				{Role: "user", Content: prompt},
			}, 4000)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}

			clean := strings.TrimSpace(resp)
			clean = strings.TrimPrefix(clean, "```json")
			clean = strings.TrimPrefix(clean, "```")
			clean = strings.TrimSuffix(clean, "```")
			clean = strings.TrimSpace(clean)

			json.NewEncoder(w).Encode(map[string]any{
				"analysis":     clean,
				"generated_at": time.Now().Format("2006-01-02 15:04"),
			})
		})

		// Static files
		if webDir != "" {
			dashboardHTML := filepath.Join(webDir, "dashboard.html")
			inspireHTML := filepath.Join(webDir, "inspire.html")
			indexHTML := filepath.Join(webDir, "index.html")
			mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case "/", "/dashboard", "/dashboard.html":
					http.ServeFile(w, r, dashboardHTML)
				case "/inspire", "/inspire.html":
					http.ServeFile(w, r, inspireHTML)
				case "/index", "/index.html":
					http.ServeFile(w, r, indexHTML)
				default:
					fs := http.FileServer(http.Dir(webDir))
					fs.ServeHTTP(w, r)
				}
			})
		} else {
			mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				fmt.Fprint(w, `<!DOCTYPE html><html><head><meta charset="UTF-8"><title>News4Coder</title></head>
<body style="font-family:system-ui;padding:2rem;text-align:center;color:#666;">
<h1>News4Coder Web</h1>
<p>请从项目根目录运行 news4coder web 以启用完整 Web 界面。</p>
<p><a href="/inspire">灵感页面</a> | <a href="/dashboard">工作台</a></p>
</body></html>`)
			})
		}

		addr := fmt.Sprintf(":%d", webPort)
		fmt.Printf("%s %s\n", cyan("✨"), bold("News4Coder Web 服务启动中..."))
		fmt.Printf("%s 本地数据库: ~/.news4coder/news4coder.db\n", green("●"))
		fmt.Printf("%s 访问地址: %s\n\n", green("●"), bold(fmt.Sprintf("http://localhost%s", addr)))
		fmt.Printf("%s 工作台:    %s\n", green("●"), bold(fmt.Sprintf("http://localhost%s/dashboard", addr)))
		fmt.Printf("%s 灵感页面:  %s\n\n", green("●"), bold(fmt.Sprintf("http://localhost%s/inspire", addr)))

		if err := http.ListenAndServe(addr, mux); err != nil {
			return fmt.Errorf("Web 服务启动失败: %w", err)
		}
		return nil
	},
}

func findWebDir() string {
	exePath, err := os.Executable()
	if err == nil {
		candidate := filepath.Join(filepath.Dir(exePath), "web")
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return candidate
		}
	}
	cwd, err := os.Getwd()
	if err == nil {
		candidate := filepath.Join(cwd, "web")
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return candidate
		}
	}
	return ""
}

func init() {
	rootCmd.AddCommand(webCmd)
	webCmd.Flags().IntVarP(&webPort, "port", "p", 8080, "Web 服务端口（默认 8080）")
}
