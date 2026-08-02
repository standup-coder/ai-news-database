# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

## [Unreleased]

### Added
- FTS5 中文单字分词支持：中文关键词可直接命中标题/摘要/LLM 摘要（`internal/db/fts.go`），旧版索引自动迁移重建
- `RebuildFTSIndex` 全量重建索引能力，便于手动修复索引不一致
- cmd 层 e2e 测试（`e2e_test.go`）：构建真实二进制验证 help/sources/stats/search/completion 等命令链路
- 项目架构设计文档（`docs/architecture.md`）与最佳实践指南（`docs/best-practices.md`）
- 内容库新增 4 个主题：语音与音频、AI搜索与信息获取、商业与投融资、消费级AI应用，并为全部主题填充线索文档

### Changed
- `articles_fts` 从外部内容表+触发器同步改为独立 FTS 表，由 Go 写入路径同步分词后文本
- 拆分 `internal/deep_research/researcher.go`（1655 行 → 7 个职责文件：model/phases/report/websearch/cache/jsonutil/researcher）
- 拆分 `cmd/burst.go`（418 行 → burst/burst_prompts/burst_history），去重三处 `modeNames` 映射

### Security
- 浏览器扩展 `host_permissions` 从任意 localhost 端口收窄为 8080/8081 固定端口（扩展 v1.1.0）

### Added (历史未发布项)
- 灵感模式 (`inspire`) 支持自动保存 Show HN 内容到本地数据库
- 完整的自媒体推广方案文档 (`PROMOTION.md`)
- 贡献指南 (`CONTRIBUTING.md`)
- 项目变更日志 (`CHANGELOG.md`)

### Changed (历史未发布项)
- 重构 `README.md`，增加技术架构详解、模块说明、路线图、安全说明

---

## [0.1.0] - 2024-04-16

### Added
- 本地优先技术资讯终端核心功能
- 8 个官方技术源：Hacker News、GitHub Blog、Lobsters、Reddit、阮一峰博客、酷壳、V2EX、InfoQ
- 文章本地存储（SQLite + FTS5 全文索引）
- LLM 内容增强（自动生成摘要、标签、质量评分）
- 智能策展（基于评分和阅读偏好的今日必读推荐）
- TUI 收件箱（Bubble Tea 交互界面）
- RAG 问答（基于本地知识库回答技术问题）
- 阅读状态五态流：unread / read / starred / archived / discarded
- Markdown 导出
- 订阅管理和自定义源添加
- 本地全文搜索
- 订阅健康度统计
