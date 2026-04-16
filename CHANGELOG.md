# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

## [Unreleased]

### Added
- 灵感模式 (`inspire`) 支持自动保存 Show HN 内容到本地数据库
- 完整的自媒体推广方案文档 (`PROMOTION.md`)
- 贡献指南 (`CONTRIBUTING.md`)
- 项目变更日志 (`CHANGELOG.md`)

### Changed
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
