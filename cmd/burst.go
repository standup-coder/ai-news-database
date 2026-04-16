package cmd

import (
	"encoding/json"
	"fmt"
	"news4coder/internal/article"
	"news4coder/internal/config"
	"news4coder/internal/db"
	"news4coder/internal/llm"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var burstCount int
var burstFocus string
var burstMode string
var burstSelect string

type burstIdea struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Inspiration string `json:"inspiration"`
}

var burstModePrompts = map[string]struct {
	System string
	Prompt string
}{
	"cross-domain": {
		System: "你是一位技术创业顾问和产品创新专家，擅长跨界联想和组合创新。请用中文回答，返回纯 JSON。",
		Prompt: `你是一位极具创造力的产品经理和技术创业顾问。以下是从 Hacker News Show HN 收集到的新产品和项目列表。请基于这些产品进行深度分析和跨界联想，生成 %d 个全新的、有创意的、可执行的产品或项目想法。

要求：
- 每个想法必须是全新的，不是简单复制已有产品
- 融合多个已有产品的核心理念进行创新，特别关注跨领域组合
- 考虑技术可行性和市场需求
- 适合独立开发者或小团队快速启动
- 用中文描述

请严格以 JSON 数组格式返回，不要包含 markdown 代码块：
[
  {
    "title": "创意名称（简洁有力）",
    "description": "详细描述：解决了什么问题、核心功能、目标用户、技术方案要点（3-5句话）",
    "inspiration": "灵感来源：来自哪些产品的启发"
  }
]

已有产品列表：
%s%s`,
	},
	"problem": {
		System: "你是一位产品战略专家，擅长从用户痛点出发设计解决方案。请用中文回答，返回纯 JSON。",
		Prompt: `以下是从 Hacker News Show HN 收集到的新产品和项目。请分析这些产品所揭示的用户痛点和未被满足的需求，然后生成 %d 个全新的产品创意。

要求：
- 从真实用户需求出发，而不是从技术出发
- 每个想法要明确说明解决了什么痛点
- 考虑现有产品的不足之处，提出更好的方案
- 目标用户要具体，不能是"所有人"
- 用中文描述

请严格以 JSON 数组格式返回，不要包含 markdown 代码块：
[
  {
    "title": "创意名称",
    "description": "痛点分析 + 解决方案 + 目标用户 + 商业模式（3-5句话）",
    "inspiration": "灵感来源：观察到什么问题或需求"
  }
]

已有产品列表：
%s%s`,
	},
	"techstack": {
		System: "你是一位全栈技术架构师，擅长将新技术栈组合成实用产品。请用中文回答，返回纯 JSON。",
		Prompt: `以下是从 Hacker News Show HN 收集到的新产品和项目。请分析这些产品使用的技术栈和架构，然后生成 %d 个全新的技术组合创意。

要求：
- 每个想法要明确说明使用的核心技术栈
- 技术选型要有创意，不是简单照搬
- 考虑新兴技术（AI、Edge Computing、WebAssembly 等）的应用
- 给出简要的技术架构描述
- 用中文描述

请严格以 JSON 数组格式返回，不要包含 markdown 代码块：
[
  {
    "title": "创意名称",
    "description": "技术方案 + 架构设计 + 核心功能 + 部署策略（3-5句话）",
    "inspiration": "灵感来源：哪些技术组合给了启发"
  }
]

已有产品列表：
%s%s`,
	},
}

var burstCmd = &cobra.Command{
	Use:   "burst",
	Short: "灵感爆发：基于 HN 产品信息，用 LLM 生成新创意",
	Long: `分析已保存的灵感模式产品信息，通过 LLM 进行跨界联想，
生成全新的产品创意和项目想法。灵感来源于你本地收藏的 Show HN 数据。

支持三种模式：
  cross-domain  跨界联想（默认）—— 融合不同领域的产品创意
  problem       问题驱动 —— 从用户痛点出发设计解决方案
  techstack     技术栈组合 —— 基于新兴技术栈组合创新`,
	Example: `  # 跨界联想模式（默认）
  news4coder burst
  news4coder burst --count 5

  # 问题驱动模式
  news4coder burst --mode problem

  # 技术栈组合模式
  news4coder burst --mode techstack

  # 聚焦某个方向
  news4coder burst --focus "AI + 开发者工具"

  # 基于特定文章生成
  news4coder burst --select 1,3,7

  # 查看历史
  news4coder burst history`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cyan := color.New(color.FgCyan).SprintFunc()
		bold := color.New(color.Bold).SprintFunc()
		yellow := color.New(color.FgYellow).SprintFunc()
		magenta := color.New(color.FgMagenta).SprintFunc()
		gray := color.New(color.FgHiBlack).SprintFunc()

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("加载配置失败: %w", err)
		}
		if cfg.LLM.APIKey == "" {
			return fmt.Errorf("LLM API Key 未配置。请编辑 ~/.news4coder/config.json")
		}

		database, err := db.New()
		if err != nil {
			return fmt.Errorf("初始化数据库失败: %w", err)
		}
		defer database.Close()

		articles, err := database.GetArticles("", "inspire", 30)
		if err != nil {
			return fmt.Errorf("查询灵感数据失败: %w", err)
		}

		if len(articles) == 0 {
			fmt.Printf("%s 暂无灵感数据。请先运行 news4coder inspire 获取。\n", yellow("!"))
			return nil
		}

		if burstSelect != "" {
			selected, parseErr := parseSelectIndices(burstSelect, len(articles))
			if parseErr != nil {
				return parseErr
			}
			var filtered []article.Article
			for _, idx := range selected {
				if idx >= 0 && idx < len(articles) {
					filtered = append(filtered, articles[idx])
				}
			}
			if len(filtered) == 0 {
				return fmt.Errorf("选择的索引无效")
			}
			articles = filtered
		}

		mode := burstMode
		if _, ok := burstModePrompts[mode]; !ok {
			mode = "cross-domain"
		}

		modeNames := map[string]string{
			"cross-domain": "跨界联想",
			"problem":      "问题驱动",
			"techstack":    "技术栈组合",
		}

		fmt.Printf("%s %s\n", cyan("💥"), bold("灵感爆发模式"))
		fmt.Printf("   模式：%s", magenta(modeNames[mode]))
		fmt.Printf(" · 基于 %d 条灵感数据", len(articles))
		if burstFocus != "" {
			fmt.Printf(" · 聚焦：%s", magenta(burstFocus))
		}
		fmt.Println()
		fmt.Printf("%s 正在调用 LLM...\n\n", cyan("⟳"))

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
		if burstFocus != "" {
			focusClause = fmt.Sprintf("\n\n用户希望聚焦的方向是：%s。请围绕这个方向展开联想。", burstFocus)
		}

		template, _ := burstModePrompts[mode]
		prompt := fmt.Sprintf(template.Prompt, burstCount, productsText, focusClause)

		client := llm.NewClient(&cfg.LLM)
		resp, err := client.Chat([]llm.Message{
			{Role: "system", Content: template.System},
			{Role: "user", Content: prompt},
		}, 4000)
		if err != nil {
			return fmt.Errorf("LLM 请求失败: %w", err)
		}

		clean := strings.TrimSpace(resp)
		clean = strings.TrimPrefix(clean, "```json")
		clean = strings.TrimPrefix(clean, "```")
		clean = strings.TrimSuffix(clean, "```")
		clean = strings.TrimSpace(clean)

		var ideas []burstIdea
		if err := json.Unmarshal([]byte(clean), &ideas); err != nil {
			fmt.Printf("%s LLM 返回内容：\n%s\n", yellow("⚠"), resp)
			return fmt.Errorf("解析 LLM 结果失败: %w", err)
		}

		now := time.Now().Format("2006-01-02 15:04")
		fmt.Println(bold("━━━ 💥 灵感爆发 ━━━"))
		fmt.Printf("   %s · %s · 基于 %d 条灵感数据\n\n", gray(now), modeNames[mode], len(articles))

		for i, idea := range ideas {
			fmt.Printf("%s %s\n", cyan(fmt.Sprintf("%2d.", i+1)), bold(idea.Title))
			fmt.Printf("   %s\n", wrapText(idea.Description, 76, "   "))
			fmt.Printf("   %s %s\n", gray("←"), gray(idea.Inspiration))
			fmt.Println()
		}

		fmt.Println(bold(fmt.Sprintf("━━━ 共 %d 个创意想法 ━━━", len(ideas))))

		ideasJSON, _ := json.Marshal(ideas)
		burstID, saveErr := database.SaveBurstResult(mode, burstFocus, string(ideasJSON), len(articles))
		if saveErr == nil {
			fmt.Printf("%s 已保存到历史记录 (#%d)\n", gray("💾"), burstID)
		}

		return nil
	},
}

var burstHistoryCmd = &cobra.Command{
	Use:   "history",
	Short: "查看灵感爆发历史记录",
	RunE: func(cmd *cobra.Command, args []string) error {
		cyan := color.New(color.FgCyan).SprintFunc()
		bold := color.New(color.Bold).SprintFunc()
		gray := color.New(color.FgHiBlack).SprintFunc()
		yellow := color.New(color.FgYellow).SprintFunc()

		database, err := db.New()
		if err != nil {
			return fmt.Errorf("初始化数据库失败: %w", err)
		}
		defer database.Close()

		results, err := database.GetBurstResults(20)
		if err != nil || len(results) == 0 {
			fmt.Printf("%s 暂无灵感爆发历史记录。\n", yellow("!"))
			return nil
		}

		modeNames := map[string]string{
			"cross-domain": "跨界联想",
			"problem":      "问题驱动",
			"techstack":    "技术栈组合",
		}

		fmt.Println(bold("━━━ 💥 灵感爆发历史 ━━━"))
		fmt.Println()

		for _, r := range results {
			var ideas []burstIdea
			json.Unmarshal([]byte(r.Ideas), &ideas)

			modeName := modeNames[r.Mode]
			if modeName == "" {
				modeName = r.Mode
			}

			focusInfo := ""
			if r.Focus != "" {
				focusInfo = fmt.Sprintf(" · 聚焦：%s", r.Focus)
			}

			fmt.Printf("%s #%d · %s · %s模式 · %d 条数据%s · %d 个创意\n",
				cyan(">"),
				r.ID,
				gray(r.CreatedAt),
				modeName,
				r.BasedOn,
				focusInfo,
				len(ideas),
			)

			for j, idea := range ideas {
				if j >= 2 {
					fmt.Printf("   %s ...还有 %d 个创意\n", gray("·"), len(ideas)-2)
					break
				}
				fmt.Printf("   %s %s\n", gray(fmt.Sprintf("%d.", j+1)), idea.Title)
			}
			fmt.Println()
		}

		fmt.Printf("共 %d 条记录。使用 %s 查看详情。\n", len(results), bold("news4coder burst show <id>"))
		return nil
	},
}

var burstShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "查看指定灵感爆发记录详情",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cyan := color.New(color.FgCyan).SprintFunc()
		bold := color.New(color.Bold).SprintFunc()
		gray := color.New(color.FgHiBlack).SprintFunc()

		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("无效的 ID: %s", args[0])
		}

		database, err := db.New()
		if err != nil {
			return fmt.Errorf("初始化数据库失败: %w", err)
		}
		defer database.Close()

		result, err := database.GetBurstResult(id)
		if err != nil {
			return fmt.Errorf("未找到记录 #%d", id)
		}

		var ideas []burstIdea
		json.Unmarshal([]byte(result.Ideas), &ideas)

		modeNames := map[string]string{
			"cross-domain": "跨界联想",
			"problem":      "问题驱动",
			"techstack":    "技术栈组合",
		}

		modeName := modeNames[result.Mode]
		focusInfo := ""
		if result.Focus != "" {
			focusInfo = fmt.Sprintf(" · 聚焦：%s", result.Focus)
		}

		fmt.Println(bold("━━━ 💥 灵感爆发 ━━━"))
		fmt.Printf("   #%d · %s · %s模式 · %d 条数据%s\n\n",
			result.ID, gray(result.CreatedAt), modeName, result.BasedOn, focusInfo)

		for i, idea := range ideas {
			fmt.Printf("%s %s\n", cyan(fmt.Sprintf("%2d.", i+1)), bold(idea.Title))
			fmt.Printf("   %s\n", wrapText(idea.Description, 76, "   "))
			fmt.Printf("   %s %s\n", gray("←"), gray(idea.Inspiration))
			fmt.Println()
		}

		fmt.Println(bold(fmt.Sprintf("━━━ 共 %d 个创意想法 ━━━", len(ideas))))
		return nil
	},
}

func parseSelectIndices(s string, maxLen int) ([]int, error) {
	parts := strings.Split(s, ",")
	var indices []int
	for _, p := range parts {
		p = strings.TrimSpace(p)
		idx, err := strconv.Atoi(p)
		if err != nil {
			return nil, fmt.Errorf("无效的索引: %s", p)
		}
		if idx < 1 || idx > maxLen {
			return nil, fmt.Errorf("索引 %d 超出范围 (1-%d)", idx, maxLen)
		}
		indices = append(indices, idx-1)
	}
	return indices, nil
}

func init() {
	rootCmd.AddCommand(burstCmd)
	burstCmd.Flags().IntVarP(&burstCount, "count", "n", 3, "生成创意数量（默认 3）")
	burstCmd.Flags().StringVarP(&burstFocus, "focus", "f", "", "聚焦方向，如 \"AI + 开发者工具\"")
	burstCmd.Flags().StringVarP(&burstMode, "mode", "M", "cross-domain", "模式：cross-domain / problem / techstack")
	burstCmd.Flags().StringVarP(&burstSelect, "select", "S", "", "选择特定文章（如 1,3,7）")

	burstCmd.AddCommand(burstHistoryCmd)
	burstCmd.AddCommand(burstShowCmd)
}
