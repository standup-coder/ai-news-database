# News4Coder - 程序员个人信息终端

> 你的高质量技术资讯大本营，兼顾断舍离与极速查询。  
> **本地优先 · 数据主权** —— 你的数据永远属于你。

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://golang.org)
[![SQLite](https://img.shields.io/badge/SQLite-3-003B57?style=flat&logo=sqlite)](https://sqlite.org)

一个优雅、简单、强大的 Go 语言命令行工具，专为程序员设计。采用 **LLM-Native** 架构和**本地优先**理念，支持多源 API 采集、智能网页抓取、LLM 内容增强（摘要/标签/评分）、智能策展和基于本地知识库的 RAG 问答。

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

## 🎯 官方信息源

news4coder 内置 8 个精选高质量技术源，覆盖中英文：

| 别名 | 名称 | 说明 |
|------|------|------|
| `hn` | Hacker News | YC 技术新闻聚合（Algolia API） |
| `github` | GitHub Blog | GitHub 官方博客（Jina Reader） |
| `lobsters` | lobste.rs | 友好程序员社区（HTML 解析） |
| `reddit` | Reddit r/programming | Reddit 编程热门（JSON API） |
| `ruanyf` | 阮一峰的网络日志 | 中文技术博客标杆（Jina Reader） |
| `coolshell` | 酷壳 CoolShell | 左耳朵耗子的技术博客（Jina Reader） |
| `v2ex` | V2EX | 中文技术社区（API） |
| `infoq` | InfoQ 中文站热点清单 | 热点文章列表（专用抓取器） |

使用 `news4coder sources` 查看完整列表。

---

## 快速开始

### 安装

确保已安装 Go 1.22 或更高版本：

```bash
# 克隆仓库
git clone <repository-url>
cd news4coder

# 编译项目
go build -o news4coder

# 或直接运行
go run main.go
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

### 日常使用工作流

```bash
# 1. 早上：一键同步所有官方源
news4coder sync

# 2. 调用 LLM 增强内容（生成摘要、标签、评分）
news4coder enrich

# 3. 查看智能策展的今日必读
news4coder curate --top 10 --auto-discard

# 4. 在 TUI 中快速处理文章
news4coder inbox

# 5. 随时向知识库提问
news4coder ask "最近关于 Rust 内存安全有哪些新观点？"

# 6. 晚上：归档已读 + 清理旧文章
news4coder archive
news4coder cleanup
```

---

## 命令参考

### 核心工作流

| 命令 | 说明 | 示例 |
|------|------|------|
| `sync` | 同步官方源到本地数据库 | `news4coder sync --source hn` |
| `enrich` | LLM 内容增强 | `news4coder enrich --limit 10` |
| `curate` | 智能策展 | `news4coder curate --top 10` |
| `ask` | RAG 问答 | `news4coder ask "问题"` |
| `inbox` | TUI 收件箱 | `news4coder inbox` |

### 文章管理

| 命令 | 说明 | 示例 |
|------|------|------|
| `list` | 列出订阅/文章 | `news4coder list --articles --status unread` |
| `read` | 标记已读 | `news4coder read 1 2 3` |
| `star` | 收藏 | `news4coder star 42` |
| `discard` | 丢弃 | `news4coder discard 5 6` |
| `archive` | 批量归档 | `news4coder archive` |

### 搜索与导出

| 命令 | 说明 | 示例 |
|------|------|------|
| `search` | 全文搜索 | `news4coder search "golang"` |
| `stats` | 订阅健康度 | `news4coder stats` |
| `export` | 导出文章 | `news4coder export --status starred --output fav.md` |
| `cleanup` | 清理旧文章 | `news4coder cleanup` |

### TUI 收件箱快捷键

```
j/↓      向下移动
k/↑      向上移动
r        标记为已读
s        收藏
d        丢弃
a        归档
1/2/3    切换视图（全部/未读/收藏）
q/Esc    退出
```

---

## 项目结构

```
news4coder/
├── cmd/                    # CLI 命令定义
│   ├── root.go            # 根命令
│   ├── sync.go            # 同步文章
│   ├── enrich.go          # LLM 内容增强
│   ├── curate.go          # 智能策展
│   ├── ask.go             # RAG 问答
│   ├── inbox.go           # TUI 收件箱
│   ├── list.go            # 列出订阅/文章
│   ├── read.go            # 标记已读
│   ├── star.go            # 收藏
│   ├── discard.go         # 丢弃
│   ├── archive.go         # 批量归档
│   ├── search.go          # 全文搜索
│   ├── stats.go           # 订阅健康度
│   ├── note.go            # 添加笔记
│   ├── tag.go             # 添加标签
│   ├── export.go          # 导出文章
│   ├── cleanup.go         # 清理旧文章
│   └── sources.go         # 官方源列表
├── internal/              # 内部模块
│   ├── article/           # 文章模型
│   ├── config/            # 应用配置管理
│   ├── crawler/           # 多源内容采集器
│   ├── curator/           # 智能策展
│   ├── db/                # SQLite 数据库
│   ├── enricher/          # LLM 内容增强
│   ├── llm/               # LLM 统一客户端
│   ├── rag/               # 检索增强生成
│   ├── official/          # 官方信息源注册表
│   ├── subscription/      # 订阅管理
│   └── tui/               # TUI 组件
├── local/                 # 本地运行文件（不提交到仓库）
│   ├── debug/             # 调试输出
│   ├── logs/              # 本地日志
│   └── qoder/             # 开发任务配置
├── web/                   # Web 界面
│   ├── index.html
│   ├── css/
│   ├── js/
│   └── assets/
├── main.go               # 程序入口
├── go.mod                # 依赖管理
├── go.sum                # 依赖校验
└── README.md             # 项目说明
```

> **注意**: `local/` 目录用于存放本地运行时生成的文件（调试输出、日志等），不会被提交到 Git 仓库。详见 `local/README.md`。

---

## 技术栈

- **语言**: Go 1.22+
- **CLI 框架**: [Cobra](https://github.com/spf13/cobra)
- **数据库**: SQLite (modernc.org/sqlite，纯 Go，无 CGO)
- **全文搜索**: SQLite FTS5
- **TUI**: [bubbletea](https://github.com/charmbracelet/bubbletea) + [lipgloss](https://github.com/charmbracelet/lipgloss)
- **LLM Client**: 标准库 HTTP + OpenAI 兼容 API
- **网页提取**: [Jina AI Reader](https://r.jina.ai)（免费，返回 Markdown）
- **HTML 解析**: [goquery](https://github.com/PuerkitoBio/goquery)
- **彩色输出**: [color](https://github.com/fatih/color)

---

## 数据与配置

### 数据存储位置

应用数据和配置保存在用户主目录：

| 平台 | 路径 |
|------|------|
| macOS/Linux | `~/.news4coder/` |
| Windows | `%USERPROFILE%\.news4coder\` |

### 文件说明

```
~/.news4coder/
├── config.json          # LLM 配置
├── news4coder.db        # SQLite 数据库（文章、阅读状态、全文索引）
└── subscriptions.json   # 旧订阅配置（已弃用）
```

### 数据主权原则

1. **完全本地存储**: 所有数据保存在本地 SQLite，不上传云端
2. **无需账号**: 没有注册、登录、授权流程
3. **零遥测**: 不收集任何使用数据
4. **可导出**: 支持 Markdown/JSON/CSV 导出
5. **可备份**: 数据库是标准 SQLite 文件，可随时备份

---

## 开发

### 构建

```bash
# 编译为可执行文件
go build -o news4coder

# 跨平台编译
# Windows
GOOS=windows GOARCH=amd64 go build -o news4coder.exe

# macOS
GOOS=darwin GOARCH=amd64 go build -o news4coder

# Linux
GOOS=linux GOARCH=amd64 go build -o news4coder
```

### 运行测试

```bash
go test ./...
```

---

## 文档

- [README.md](README.md) - 项目主文档
- [EXAMPLES.md](EXAMPLES.md) - 使用示例和常见问题
- [local/README.md](local/README.md) - 本地目录说明
- [web/README.md](web/README.md) - Web 界面说明

---

## 贡献

欢迎提交 Issue 和 PR！

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
