# News4Coder - 程序员个人信息终端

> 你的高质量技术资讯大本营，兼顾断舍离与极速查询。  
> **本地优先 · 数据主权** —— 你的数据永远属于你。

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://golang.org)
[![SQLite](https://img.shields.io/badge/SQLite-3-003B57?style=flat&logo=sqlite)](https://sqlite.org)

一个优雅、简单、强大的 Go 语言命令行工具，专为程序员设计。采用 **LLM-Native** 架构和**本地优先**理念，支持多源 API 采集、智能网页抓取、LLM 内容增强（摘要/标签/评分）、智能策展和基于本地知识库的 RAG 问答。

---

## 目录

- [为什么选择 News4Coder？](#为什么选择-news4coder)
- [功能特性](#功能特性)
- [官方信息源](#官方信息源)
- [快速开始](#快速开始)
- [完整使用指南](#完整使用指南)
- [命令参考](#命令参考)
- [TUI 收件箱](#tui-收件箱)
- [技术架构](#技术架构)
- [数据库设计](#数据库设计)
- [模块详解](#模块详解)
- [开发指南](#开发指南)
- [故障排除](#故障排除)
- [许可证](#许可证)

---

## 为什么选择 News4Coder？

| 特性 | News4Coder | 云端阅读器 |
|------|------------|-----------|
| **数据存储** | 本地 SQLite | 云端服务器 |
| **账号注册** | 不需要 | 必需 |
| **隐私保护** | 数据不上传 | 服务商可查看 |
| **服务可用性** | 永久可用 | 依赖服务商 |
| **离线使用** | 完全支持 | 有限支持 |
| **数据导出** | 完全自由 | 受限 |

---

## 功能特性

- ✅ **本地优先**: 所有数据存储在本地 SQLite，无需云端账号
- 🔒 **数据主权**: 你的阅读记录、收藏、笔记完全属于你
- 🌐 **多源智能采集**: 直接调用 HN/Reddit/V2EX API + Jina AI Reader，告别 RSS
- 🤖 **LLM 内容增强**: 自动生成高质量摘要、技术标签、质量评分（0-10）
- 🎯 **智能策展**: 基于质量评分和阅读偏好，自动生成「今日必读」
- 💾 **本地知识库**: SQLite + FTS5 全文检索，构建你的技术资产
- 📖 **阅读状态流**: unread / read / starred / archived / discarded
- 🧹 **断舍离支持**: 批量归档、自动清理、语义去重
- 🔍 **极速搜索**: 本地全文搜索，秒级响应
- 💬 **RAG 问答**: 基于本地知识库回答，标注引用来源
- 📺 **TUI 收件箱**: 交互式极速浏览，秒级决策
- 📝 **快速输入**: 笔记、标签、Markdown 导出
- 🎨 **优雅输出**: 彩色终端输出，清晰易读
- 🚀 **快速启动**: 单二进制文件，无需额外依赖

---

## 官方信息源

news4coder 内置 8 个精选高质量技术源，覆盖中英文：

| 别名 | 名称 | 采集方式 | 说明 |
|------|------|----------|------|
| `hn` | Hacker News | Algolia API | YC 技术新闻聚合 |
| `github` | GitHub Blog | Jina Reader | GitHub 官方博客 |
| `lobsters` | lobste.rs | HTML 解析 | 友好程序员社区 |
| `reddit` | Reddit r/programming | JSON API | Reddit 编程热门 |
| `ruanyf` | 阮一峰的网络日志 | Jina Reader | 中文技术博客标杆 |
| `coolshell` | 酷壳 CoolShell | Jina Reader | 左耳朵耗子的技术博客 |
| `v2ex` | V2EX | API | 中文技术社区 |
| `infoq` | InfoQ 中文站热点清单 | 专用抓取器 | 热点文章列表 |

使用 `news4coder sources` 查看完整列表，也可以直接用别名获取内容，如 `news4coder hn`。

---

## 快速开始

### 安装

#### 方式一：一键安装脚本（推荐）

```bash
curl -sSL https://get.news4coder.dev | bash
```

#### 方式二：从源码构建

确保已安装 Go 1.25 或更高版本：

```bash
# 克隆仓库
git clone https://github.com/news4coder/news4coder.git
cd news4coder

# 编译项目
make build

# 或直接运行
make run

# 安装到 $GOPATH/bin
make install
```

#### 方式三：go install

```bash
go install github.com/news4coder/news4coder@latest
```

#### 验证安装

```bash
news4coder --version
```

### 配置 LLM（可选但推荐）

编辑 `~/.news4coder/config.json`，填入你的 LLM API 信息：

```json
{
  "llm": {
    "base_url": "https://api.openai.com/v1",
    "api_key": "sk-xxxxxxxx",
    "model": "gpt-4o-mini",
    "embedding_model": "text-embedding-3-small",
    "enrich_max_tokens": 2000,
    "ask_max_tokens": 4000
  }
}
```

> 支持任何 OpenAI 兼容接口，包括 Ollama 本地服务（`http://localhost:11434/v1`）。

### 首次运行

```bash
# 查看帮助
./news4coder --help

# 查看可用官方源
./news4coder sources

# 同步所有官方源
./news4coder sync

# 查看已同步文章
./news4coder list --articles
```

---

## 完整使用指南

### 场景 1: 早晨资讯同步（5分钟）

```bash
# 1. 同步所有官方源
./news4coder sync

# 2. 对新文章进行 LLM 增强（生成摘要、标签、评分）
./news4coder enrich --limit 20

# 3. 查看智能策展的今日必读
./news4coder curate --top 10
```

### 场景 2: 深度阅读工作流

```bash
# 1. 进入 TUI 收件箱浏览文章
./news4coder inbox
```

在 TUI 中：
- `j/k` 或 `↓/↑` - 浏览文章
- `r` - 标记为已读
- `s` - 收藏
- `d` - 丢弃
- `a` - 归档
- `1/2/3` - 切换视图（全部/未读/收藏）
- `q` - 退出

```bash
# 2. 为重要文章添加笔记
./news4coder note <article-id> "核心观点记录..."

# 3. 添加自定义标签
./news4coder tag <article-id> "golang,concurrency"
```

### 场景 3: 知识库问答

```bash
# 基于本地知识库回答问题
./news4coder ask "Go 和 Rust 在并发模型上有什么差异？"

./news4coder ask "最近关于微服务架构有哪些最佳实践？"
```

输出示例：
```
⟳ 正在检索知识库并生成回答...

回答：
Go 和 Rust 在并发模型上有以下主要差异...

━━━ 引用来源 ━━━
[1] Go 并发模式详解 (Hacker News)
    https://...
[2] Rust 异步编程指南 (GitHub Blog)
    https://...

💡 提示: 知识库内容越丰富，回答质量越高。建议定期执行 sync + enrich。
```

### 场景 4: 全文搜索

```bash
# 搜索本地知识库
./news4coder search "kubernetes"

./news4coder search "微服务架构"
```

### 场景 5: 导出与备份

```bash
# 导出收藏的文章为 Markdown
./news4coder export --status starred --output my-favorites.md

# 导出已读文章
./news4coder export --status read --output read-articles.md

# 直接备份数据库文件
cp ~/.news4coder/news4coder.db ~/backups/news4coder-$(date +%Y%m%d).db
```

### 场景 6: 定期维护

```bash
# 查看订阅健康度
./news4coder stats

# 归档所有已读文章
./news4coder archive

# 清理过期文章（丢弃超过7天、归档超过30天）
./news4coder cleanup
```

### 高级技巧

#### 组合命令

```bash
# 同步后立即增强
./news4coder sync && ./news4coder enrich

# 策展后进入收件箱
./news4coder curate --top 10 && ./news4coder inbox
```

#### 使用别名

在 `~/.bashrc` 或 `~/.zshrc` 中添加：

```bash
alias n4c='./news4coder'
alias n4c-sync='./news4coder sync'
alias n4c-inbox='./news4coder inbox'
alias n4c-ask='./news4coder ask'
```

#### 定时自动同步

使用 cron（Linux/macOS）：

```bash
# 每天早上 8 点同步
0 8 * * * cd /path/to/news4coder && ./news4coder sync >> ~/.news4coder/sync.log 2>&1
```

---

## 命令参考

### 核心工作流

| 命令 | 说明 | 示例 |
|------|------|------|
| `sync` | 同步官方源到本地数据库 | `news4coder sync --source hn` |
| `enrich` | LLM 内容增强 | `news4coder enrich --limit 10` |
| `curate` | 智能策展 | `news4coder curate --top 10 --auto-discard` |
| `ask` | RAG 问答 | `news4coder ask "问题"` |
| `inbox` | TUI 收件箱 | `news4coder inbox` |

### 文章管理

| 命令 | 说明 | 示例 |
|------|------|------|
| `list` | 列出订阅/文章 | `news4coder list -a --status unread` |
| `read` | 标记已读 | `news4coder read 1 2 3` |
| `star` | 收藏 | `news4coder star 42` |
| `discard` | 丢弃 | `news4coder discard 5 6` |
| `archive` | 批量归档 | `news4coder archive` |
| `note` | 添加笔记 | `news4coder note 42 "笔记"` |
| `tag` | 添加标签 | `news4coder tag 42 "tag1,tag2"` |

### 搜索与导出

| 命令 | 说明 | 示例 |
|------|------|------|
| `search` | 全文搜索 | `news4coder search "golang"` |
| `stats` | 订阅健康度 | `news4coder stats` |
| `export` | 导出文章 | `news4coder export --status starred` |
| `cleanup` | 清理旧文章 | `news4coder cleanup` |

### 订阅管理

| 命令 | 说明 | 示例 |
|------|------|------|
| `sources` | 查看官方源列表 | `news4coder sources` |
| `add` | 添加自定义订阅 | `news4coder add -n "MyBlog" -u "https://..."` |
| `remove` | 删除订阅 | `news4coder remove -n "MyBlog"` |
| `fetch` | 获取订阅最新内容 | `news4coder fetch -n infoq` |

### 快捷别名命令

| 命令 | 说明 | 示例 |
|------|------|------|
| `hn` | 直接获取 Hacker News | `news4coder hn` |
| `infoq` | 直接获取 InfoQ 热点 | `news4coder infoq` |
| `v2ex` | 直接获取 V2EX 热门 | `news4coder v2ex` |
| `github` | 直接获取 GitHub Blog | `news4coder github` |
| `reddit` | 直接获取 Reddit 热门 | `news4coder reddit` |
| `lobsters` | 直接获取 Lobsters | `news4coder lobsters` |
| `ruanyf` | 直接获取阮一峰博客 | `news4coder ruanyf` |
| `coolshell` | 直接获取酷壳 | `news4coder coolshell` |
| `inspire` | HN 新产品灵感（**自动保存到本地数据库**） | `news4coder inspire --days 7` |

---

## TUI 收件箱

启动命令：`news4coder inbox`

### 快捷键

```
j / ↓    向下移动
k / ↑    向上移动
r        标记为已读
s        收藏
d        丢弃
a        归档
1        显示全部文章
2        显示未读文章
3        显示收藏文章
q / Esc  退出
```

### 界面说明

TUI 界面分为三个区域：
1. **顶部标题栏** - 显示当前视图模式（全部/未读/收藏）和文章总数
2. **文章列表** - 显示文章 ID、来源、标题、质量评分和阅读状态图标
3. **底部预览区** - 显示当前选中文章的标题、链接、标签、笔记和摘要

阅读状态图标：
- `○` - 未读
- `✓` - 已读
- `★` - 收藏
- `✗` - 丢弃
- `▣` - 归档

---

## 技术架构

### 项目目录结构

```
news4coder/
├── cmd/                          # CLI 命令定义（Cobra）
│   ├── root.go                   # 根命令 + 官方源别名处理
│   ├── sync.go                   # 同步文章到本地数据库
│   ├── enrich.go                 # LLM 内容增强
│   ├── curate.go                 # 智能策展
│   ├── ask.go                    # RAG 问答
│   ├── inbox.go                  # TUI 收件箱
│   ├── list.go                   # 列出订阅/文章
│   ├── read.go                   # 标记已读
│   ├── star.go                   # 收藏
│   ├── discard.go                # 丢弃
│   ├── archive.go                # 批量归档
│   ├── search.go                 # 全文搜索
│   ├── stats.go                  # 订阅健康度
│   ├── note.go                   # 添加笔记
│   ├── tag.go                    # 添加标签
│   ├── export.go                 # 导出文章
│   ├── cleanup.go                # 清理旧文章
│   ├── sources.go                # 官方源列表
│   ├── fetch.go                  # 获取订阅内容
│   ├── add.go / remove.go        # 订阅增删
│   ├── infoq.go                  # InfoQ 专注模式
│   ├── inspire.go                # HN 灵感模式
│   └── console_windows.go        # Windows UTF-8 支持
├── internal/                     # 内部模块
│   ├── article/                  # 文章模型（状态、结构体）
│   ├── config/                   # 应用配置管理（JSON）
│   ├── crawler/                  # 多源内容采集器
│   │   ├── factory.go            # 采集器工厂
│   │   ├── hn.go                 # Hacker News API 采集
│   │   ├── reddit.go             # Reddit API 采集
│   │   ├── v2ex.go               # V2EX API 采集
│   │   ├── generic.go            # 通用 HTML 列表采集
│   │   ├── reader.go             # Jina AI Reader 网页提取
│   │   └── model.go              # 采集器接口定义
│   ├── curator/                  # 智能策展（评分 + 偏好匹配）
│   ├── db/                       # SQLite 数据库 + FTS5
│   ├── dedup/                    # 语义去重（Jaccard + 标签重叠）
│   ├── enricher/                 # LLM 内容增强客户端
│   ├── llm/                      # LLM 统一客户端 + Embedding
│   ├── official/                 # 官方信息源注册表
│   │   ├── registry.go           # 8 个内置源注册
│   │   ├── model.go              # 源结构体
│   │   ├── fetcher.go            # 抓取器工厂
│   │   └── infoq_fetcher.go      # InfoQ 专用抓取器
│   ├── rag/                      # 检索增强生成（FTS + LLM）
│   ├── search/                   # DuckDuckGo 搜索引擎
│   ├── storage/                  # 订阅配置持久化（JSON）
│   ├── subscription/             # 订阅管理
│   └── tui/                      # TUI 组件（Bubble Tea）
├── web/                          # Web 宣传界面
│   ├── index.html                # 主页面
│   ├── css/style.css             # CNCF 设计系统样式
│   ├── js/main.js                # 交互动画
│   └── assets/                   # 图标资源
├── local/                        # 本地运行文件（不提交）
├── main.go                       # 程序入口
├── go.mod                        # 依赖管理
└── README.md                     # 项目说明
```

### 核心数据流

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   官方源    │────▶│  Crawler/   │────▶│   SQLite    │
│ (HN/V2EX等) │     │   Fetcher   │     │    (db)     │
└─────────────┘     └─────────────┘     └──────┬──────┘
                                                │
                       ┌────────────────────────┘
                       ▼
              ┌─────────────────┐
              │   Enricher      │
              │ (LLM 摘要/标签)  │
              └────────┬────────┘
                       │
         ┌─────────────┼─────────────┐
         ▼             ▼             ▼
    ┌─────────┐   ┌─────────┐   ┌─────────┐
    │ Curator │   │  RAG    │   │ Search  │
    │(智能推荐)│   │(知识问答)│   │(全文检索)│
    └────┬────┘   └────┬────┘   └────┬────┘
         │             │             │
         └─────────────┴─────────────┘
                       │
                       ▼
                ┌─────────────┐
                │  TUI / CLI  │
                │  (用户交互)  │
                └─────────────┘
```

### 各层职责

1. **采集层** (`crawler/`): 通过 API 或 HTML 解析从各源获取文章元数据
   - `HNCrawler`: 调用 HN Algolia API
   - `RedditCrawler`: 调用 Reddit JSON API
   - `V2EXCrawler`: 调用 V2EX Hot API
   - `GenericCrawler`: 基于 goquery 的通用 HTML 列表解析
   - `JinaReader`: 通过 Jina AI Reader 提取网页正文为 Markdown

2. **存储层** (`db/`): SQLite 持久化，URL 去重，FTS5 全文索引
   - 自动迁移机制支持旧数据库升级
   - FTS5 虚拟表提供全文搜索能力
   - 触发器自动同步 FTS 索引

3. **增强层** (`enricher/`): 调用 LLM 生成摘要、标签、质量评分
   - 自动回退到 Jina Reader 获取正文
   - 内容截断控制 Token 消耗
   - JSON 格式解析增强结果

4. **策展层** (`curator/`): 基于评分和用户偏好（starred 标签）生成推荐
   - 分析用户收藏文章的标签分布
   - 质量评分 + 偏好匹配 = 策展得分
   - 自动生成推荐理由

5. **交互层** (`cmd/`, `tui/`): CLI 命令和 TUI 界面供用户操作
   - Cobra 构建的命令行接口
   - Bubble Tea + Lipgloss 构建的 TUI

6. **问答层** (`rag/`): 基于 FTS 召回 + LLM 生成带引用的回答
   - FTS5 全文检索召回相关文章
   - 失败时回退到 LIKE 模糊搜索
   - 组装上下文调用 LLM 生成答案

---

## 数据库设计

### 主表: `articles`

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | INTEGER PRIMARY KEY | 自增主键 |
| `title` | TEXT NOT NULL | 文章标题 |
| `url` | TEXT NOT NULL UNIQUE | 文章链接（唯一索引） |
| `source` | TEXT NOT NULL | 来源名称 |
| `source_alias` | TEXT NOT NULL | 来源别名 |
| `summary` | TEXT | 原始摘要 |
| `llm_summary` | TEXT | LLM 生成摘要 |
| `llm_tags` | TEXT | LLM 自动标签 |
| `quality_score` | REAL DEFAULT 0 | 质量评分（0-10） |
| `language` | TEXT | 文章语言 |
| `raw_content` | TEXT | 原始抓取内容 |
| `published_at` | DATETIME | 发布时间 |
| `fetched_at` | DATETIME DEFAULT CURRENT_TIMESTAMP | 抓取时间 |
| `enriched_at` | DATETIME | 增强时间 |
| `read_status` | TEXT DEFAULT 'unread' | 阅读状态 |
| `tags` | TEXT | 用户手动标签 |
| `note` | TEXT | 用户笔记 |

### 阅读状态

- `unread` - 未读
- `read` - 已读
- `starred` - 收藏
- `archived` - 归档
- `discarded` - 丢弃

### 全文搜索表: `articles_fts`

使用 SQLite FTS5 虚拟表，索引字段：`title`, `summary`, `llm_summary`。

通过触发器自动与 `articles` 表同步：
- `articles_ai` - INSERT 触发器
- `articles_ad` - DELETE 触发器
- `articles_au` - UPDATE 触发器

### 索引

- `idx_articles_status` - 按阅读状态查询
- `idx_articles_source` - 按来源查询
- `idx_articles_fetched_at` - 按时间排序

---

## 模块详解

### 1. 文章模型 (`internal/article`)

定义文章结构体和阅读状态常量：

```go
type Article struct {
    ID            int64
    Title         string
    URL           string
    Source        string
    SourceAlias   string
    Summary       string
    LLMSummary    string
    LLMTags       string
    QualityScore  float64
    Language      string
    RawContent    string
    PublishedAt   *time.Time
    FetchedAt     time.Time
    EnrichedAt    *time.Time
    ReadStatus    ReadStatus
    Tags          string
    Note          string
}
```

### 2. 配置管理 (`internal/config`)

配置文件存储在 `~/.news4coder/config.json`：

```go
type LLMConfig struct {
    BaseURL          string `json:"base_url"`
    APIKey           string `json:"api_key"`
    Model            string `json:"model"`
    EmbeddingModel   string `json:"embedding_model"`
    EnrichMaxTokens  int    `json:"enrich_max_tokens"`
    AskMaxTokens     int    `json:"ask_max_tokens"`
}
```

首次运行自动创建默认配置。支持补全缺失字段的默认值。

### 3. 爬虫系统 (`internal/crawler`)

#### 工厂模式

`crawler.NewCrawler(source)` 根据别名返回对应采集器：

| 别名 | 采集器 | 说明 |
|------|--------|------|
| `hn` | `HNCrawler` | Algolia API |
| `reddit` | `RedditCrawler` | Reddit JSON API |
| `v2ex` | `V2EXCrawler` | V2EX Hot API |
| `lobsters` | `GenericCrawler` | HTML 解析 `.story_liner` |
| `ruanyf` | `GenericCrawler` | HTML 解析 `article` |
| `coolshell` | `GenericCrawler` | HTML 解析 `article` |
| `github` | `GenericCrawler` | HTML 解析 `article` |
| `infoq` | `GenericCrawler` | HTML 解析 `article` |

#### Jina AI Reader

用于提取网页正文：

```go
reader := crawler.NewJinaReader()
content, err := reader.Fetch("https://example.com/article")
```

请求格式：`https://r.jina.ai/http://example.com/article`

### 4. LLM 客户端 (`internal/llm`)

OpenAI 兼容 API 的 HTTP 客户端：

- `Chat(messages, maxTokens)` - 多轮对话
- `SimpleChat(prompt, maxTokens)` - 单轮快速对话
- `GetEmbedding(text)` - 获取文本向量
- `CosineSimilarity(a, b)` - 计算余弦相似度

超时设置：Chat 120 秒，Embedding 120 秒。

### 5. 内容增强器 (`internal/enricher`)

对单篇文章执行 LLM 增强：

1. 检查 `raw_content` 长度，若不足 200 字符，尝试 Jina Reader 抓取正文
2. 截断内容至 12000 字符控制 Token
3. 调用 LLM 生成 JSON 格式结果：`summary`, `tags`, `score`, `language`
4. 解析并保存到数据库

Prompt 要求：
- 摘要：2-3 句话中文摘要
- 标签：3-5 个精准技术标签（如 golang, rust, ai, kubernetes）
- 评分：0-10 分，基于技术深度、实用性、时效性、可读性
- 语言：zh, en, ja 等

### 6. 智能策展 (`internal/curator`)

策展算法：

```
CuratorScore = QualityScore + PreferenceBonus - UnenrichedPenalty

PreferenceBonus = sum(weight * 0.5) for each matching tag
UnenrichedPenalty = 1.0 if EnrichedAt is nil
```

用户偏好标签从 `starred` 状态文章中提取，按出现频次计算权重。

### 7. RAG 问答 (`internal/rag`)

实现流程：

1. **检索**：使用 FTS5 `MATCH` 搜索相关问题，召回 Top 10
2. **回退**：FTS 失败时回退到 `LIKE` 模糊搜索
3. **组装上下文**：将相关文章的标题、来源、摘要组装为上下文
4. **生成答案**：调用 LLM，要求标注引用来源编号 `[1]`, `[2]`
5. **返回**：答案 + 引用列表

### 8. 灵感模式 (`cmd/inspire`)

`inspire` 命令不仅展示 Hacker News 上的 `Show HN` 新产品/新项目，还会**自动将结果持久化到本地 SQLite**，使其成为知识库的一部分，可被搜索、策展和 RAG 问答检索。

#### 数据源

- **API**: HN Algolia API (`search_by_date?tags=show_hn`)
- **筛选条件**: 最近 N 天（默认 7 天，可通过 `--days` 调整）
- **数量限制**: 默认最多 30 条（可通过 `--limit` 调整）

#### 数据处理流程

```
调用 HN API ──▶ 解析 JSON ──▶ 清洗标题/描述 ──▶ 组装 Article ──▶ 去重检查 ──▶ 写入 SQLite
```

#### 清洗规则

1. **标题处理**: 去除 `Show HN:` / `Show HN` 前缀，去除首尾空格
2. **描述处理**: 
   - 去除 HTML 标签（如 `<p>`, `<br>`）
   - HTML 实体解码（如 `&quot;` → `"`）
   - 去除首尾空格
   - 截断至 200 字符（超出部分以 `...` 结尾）
3. **URL 构建**: 使用 HN 讨论页作为文章 URL
   - 格式: `https://news.ycombinator.com/item?id=<objectID>`
4. **发布时间**: 从 `created_at`（RFC3339）解析为 `PublishedAt`

#### 本地存储字段映射

| Article 字段 | 值 |
|-------------|-----|
| `Title` | 清洗后的产品/项目名称 |
| `URL` | `https://news.ycombinator.com/item?id=<objectID>` |
| `Source` | `Show HN` |
| `SourceAlias` | `inspire` |
| `Summary` | 清洗后的描述（≤200字符） |
| `RawContent` | 同 `Summary` |
| `ReadStatus` | `unread` |
| `PublishedAt` | API 返回的创建时间 |

#### 去重与保存

1. **预检查去重**: 调用 `db.ArticleExists(url)` 检查该 HN 讨论页是否已存在于本地数据库
2. **写入**: 通过 `db.SaveArticle()` 插入新记录
3. **冲突处理**: 若 URL 已存在（或保存失败），自动跳过并计入 `dup` 统计
4. **输出提示**: 命令最后会显示 `💾 已保存 X 条新内容，Y 条已存在跳过`

#### 与其他功能的关联

- **搜索**: 保存后的灵感内容可通过 `news4coder search "Show HN"` 检索
- **阅读管理**: 可在 TUI 中标记为已读、收藏或丢弃
- **导出**: 可随其他文章一起导出为 Markdown
- **RAG 问答**: 参与本地知识库问答，作为新产品/项目灵感来源

### 9. 语义去重 (`internal/dedup`)

去重规则（按优先级）：

1. URL 完全相同
2. 标题 Jaccard 相似度 > 0.85（基于字符 bigram）
3. LLM 标签重叠率 > 80%

发现重复时，保留质量评分更高的一篇，另一篇标记为 `discarded`。

### 10. 搜索引擎 (`internal/search`)

基于 DuckDuckGo HTML 版本的站内搜索：

```go
engine := search.NewEngine()
results, err := engine.Search("https://example.com")
```

搜索语法：`site:example.com`，返回最多 10 条结果。
支持提取 DuckDuckGo 重定向链接中的真实 URL。

### 11. 订阅管理 (`internal/subscription` + `internal/storage`)

用户自定义订阅存储在 `~/.news4coder/subscriptions.json`：

```json
{
  "subscriptions": [
    {
      "name": "InfoQ中文站",
      "alias": "infoq",
      "url": "https://www.infoq.cn",
      "created_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

提供添加、删除、列表、查询功能。

---

## 路线图

### 近期（v0.x）

- [x] 8 个官方技术源采集（HN / V2EX / GitHub Blog / Reddit 等）
- [x] LLM 内容增强（摘要 / 标签 / 评分）
- [x] 智能策展（今日必读 Top N）
- [x] TUI 收件箱（Bubble Tea）
- [x] 本地 FTS5 全文搜索
- [x] RAG 问答（带引用来源）
- [x] Markdown 导出
- [x] 灵感模式自动本地存储（Show HN）

### 中期（v1.0）

- [ ] Web Dashboard（图形化管理界面）
- [ ] 自定义订阅源可视化配置
- [ ] 阅读统计可视化（周报 / 月报）
- [ ] 支持更多官方源（如 Dev.to、Medium、Twitter/X 列表）
- [ ] 自动同步 Scheduler（内置定时任务）

### 远期（v2.0+）

- [ ] 跨设备同步（可选的端到端加密同步）
- [ ] 浏览器扩展（一键保存文章到 News4Coder）
- [ ] 插件系统（允许社区扩展采集源）
- [ ] 移动端 App（基于本地 SQLite 的只读浏览）

---

## 安全说明

### API Key 存储

News4Coder 的 LLM 配置（`~/.news4coder/config.json`）仅存储在本地设备，不会上传至任何服务器。请确保：

- 不要将 `config.json` 提交到公共代码仓库
- 定期更换 LLM API Key
- 使用 `.gitignore` 排除本地配置文件

### 数据安全

- 所有文章数据存储在本地 SQLite 文件（`~/.news4coder/news4coder.db`）
- 程序不收集任何遥测、使用统计或崩溃报告
- 网络请求仅用于从公开技术源采集文章和调用用户配置的 LLM API

---

## 开发指南

### 构建

```bash
# 编译为可执行文件
make build

# 跨平台编译（生成所有平台二进制到 dist/）
make release

# 手动单平台编译
GOOS=linux GOARCH=amd64 go build -o news4coder
```

### 运行测试

```bash
# 运行全部测试（含覆盖率）
make test

# 快速测试
make test-short

# 代码格式化
make fmt

# 静态检查
make vet

# 查看所有可用命令
make help
```

### Web 界面本地预览

```bash
cd web && python3 -m http.server 8080
# 访问 http://localhost:8080
```

### 添加新命令

1. 在 `cmd/` 下新建 `xxx.go`
2. 定义 `xxxCmd` 并使用 `rootCmd.AddCommand(xxxCmd)`
3. 在 `init()` 中注册 Flag

### 添加新采集源

1. 在 `internal/official/registry.go` 的 `registerDefaultSources()` 中注册
2. 如需专用采集器，在 `internal/crawler/factory.go` 的 `NewCrawler()` 中实现分支
3. 如需网页提取，可直接复用 `GenericCrawler`

---

## 故障排除

### Q: 为什么数据要存储在本地？

**A**: News4Coder 坚持「本地优先」和「数据主权」理念：

- ✅ **隐私**: 无需账号，数据不上传云端
- ✅ **永久可用**: 不受服务商倒闭影响
- ✅ **可导出**: 随时备份或迁移
- ✅ **离线使用**: 没有网络也能访问已同步内容

### Q: 不配置 LLM 能用吗？

**A**: 可以。不配置 LLM 也能：
- 同步和阅读文章
- 管理阅读状态
- 搜索本地文章
- 导出文章

但无法使用：
- `enrich` - 自动生成摘要和标签
- `curate` - 智能策展
- `ask` - RAG 问答

### Q: 如何添加自定义订阅源？

**A**: 使用 `add` 命令：

```bash
news4coder add --name "MyBlog" --alias myblog --url "https://myblog.com"
```

### Q: 数据库会占用多大空间？

**A**: 取决于同步的文章数量。通常：
- 1000 篇文章约 10-20 MB
- 包含全文索引，支持快速搜索

### Q: 如何清理空间？

**A**: 定期执行：

```bash
# 清理旧文章
news4coder cleanup

# 手动查看丢弃的文章
news4coder list --articles --status discarded
```

### Q: `sync` 时某些源报错怎么办？

**A**: 
1. 检查网络连接
2. 尝试单独同步该源：`news4coder sync --source <alias>`
3. 对于 API 源，可能是 rate limit，稍后再试
4. 对于 HTML 解析源，可能是页面结构变更

### Q: Windows 终端中文显示乱码？

**A**: 程序已内置 Windows UTF-8 控制台支持（`cmd/console_windows.go`）。如果仍有问题，请在 PowerShell 中执行：

```powershell
[Console]::OutputEncoding = [System.Text.Encoding]::UTF8
```

---

## Web 界面

项目包含一个现代化的 Web 宣传页面（`web/` 目录），采用 CNCF 顶级开源项目设计美学：

- **配色方案**: Kubernetes Blue (`#326CE5`) + Cloud Native Teal (`#00A3C4`)
- **设计风格**: 卡片式布局、微妙阴影、现代圆角、响应式设计
- **核心页面**:
  - `index.html` — 产品官网（数据主权、特性、信息源、工作流、下载）
  - `inspire.html` — 灵感模式展示页（Show HN 产品/项目）
- **响应式断点**: Desktop (>1024px) / Tablet (768-1024px) / Mobile (<768px)

### 本地运行 Web

```bash
# 方式 1：Python（最简单，已预装）
cd web && python3 -m http.server 8080

# 方式 2：Node.js（live-reload 开发体验）
npx live-server web --port=8080

# 方式 3：Docker
docker run -d -p 8080:80 -v $(pwd)/web:/usr/share/nginx/html nginx:alpine
```

访问地址：
- 主站：http://localhost:8080
- 灵感页：http://localhost:8080/inspire.html

更多细节详见 [`web/README.md`](web/README.md)。

---

## 相关文档

- [`EXAMPLES.md`](EXAMPLES.md) - 使用示例和常见问题
- [`PROMOTION.md`](PROMOTION.md) - 自媒体推广方案（含文案、渠道、运营策略）
- [`BUSINESS.md`](BUSINESS.md) - 商业战略规划（演进计划、产品策略、定价与商业模式）
- [`CONTRIBUTING.md`](CONTRIBUTING.md) - 贡献指南
- [`CHANGELOG.md`](CHANGELOG.md) - 版本变更记录
- [`local/README.md`](local/README.md) - 本地目录说明
- [`web/README.md`](web/README.md) - Web 界面说明

---

## 贡献

欢迎提交 Issue 和 PR！

详细贡献指南请查看 [`CONTRIBUTING.md`](CONTRIBUTING.md)。

简要流程：

1. Fork 本仓库
2. 创建你的特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开一个 Pull Request

---

## 许可证

MIT License

## 作者

News4Coder 项目团队

---

**享受你的本地优先技术阅读体验！** 🚀  
*你的数据，永远属于你。*
