# Product

<!-- impeccable:product-schema 1 -->

## Platform

web

## Stack

纯静态单页 HTML/CSS（用户确认）：`GTM/index.html` + 本地资源，零构建链，GitHub Pages 可直接部署。

## Users

- **主受众：程序员读者**——跟进 AI 行业动态的开发者，需要高质量、可信、结构化的中文信息源；在页面上决定「是否收藏/关注这个库」。
- **次受众：内容贡献者**——愿意按模板提交线索/新闻的开发者；在页面上决定「是否参与共建」。

## Product Purpose

AI News Database（AI 新闻数据库）是一个持续沉淀 AI 高质量新闻与重点事件的中文 Markdown 内容库：按 17 个 AI 全景主题组织，一条线索 = 一个持续追加的时间线文档，另有按月归档的年度大事记。存在意义：中文世界缺少「宁可漏收、不可错收」的高信噪比 AI 新闻沉淀。成功 = 开发者把它当作可长期依赖的 AI 行业信息底座，并持续有人按规范贡献。

## Positioning

相邻产品（资讯站、Newsletter、聚合器）无法 truthfully 复制的机制组合：
1. **时间线式线索库**——每个故事线一个倒序累积的文档，看得见演进脉络，而非信息流；
2. **真实性门**——官方一手源 + ≥2 家主流媒体交叉印证才入库，存疑事件公开登记待核实清单；
3. **事实/观点严格分离**——「概述」只写事实，「分析」单列观点；
4. **纯 Markdown + git**——人可读、可 fork、可被工具消费。

## Operating Context

- 内容即仓库本身：GitHub 托管（standup-coder/ai-news-database，URL 为从工作目录推断，待确认），git 跟踪，人可读；
- 中心入口：`_topics.md`（主题索引）、`_2026大事记.md`（年度大事记，含「真实性约定」）、`docs/线索模板.md`（贡献模板）、`CONTRIBUTING.md`；
- 辅助工具在 `tools/cmd/`（自包含 Go module）：CLI、Web 界面、浏览器扩展，仅可选辅助非核心入口；含 FTS5 中文搜索；
- `aihot-mirror/` 为镜像数据，非主库内容。

## Capabilities and Constraints

- 内容库手动整理为主，质量优先于速度；
- GTM 页面为纯静态单页，不得引入 JS 框架/构建链；
- 页面内链接指向仓库真实文件（_topics.md、_2026大事记.md、docs/线索模板.md 等）；
- 未决事实：GitHub 仓库最终 URL、是否已有线上域名——页面先用相对/占位链接并集中可替换。

## Brand Commitments

- 名称口径以 README 为准：「AI News Database（AI 新闻数据库）」，全仓库统一此名（历史代号 News4Coder 已于 2026-08 全量更名）。

## Evidence on Hand

- README.md（产品口径与目录结构）、_topics.md（13 主题清单）、_2026大事记.md（2026 按月大事记，真实性约定原文）、docs/线索模板.md、CONTRIBUTING.md；
- 真实内容样本：各主题线索文件（如 基础模型/Claude.md、开源模型与生态/DeepSeek.md）；
- **不存在**：logo/品牌视觉资产、用户证言、流量/Star 数据、媒体报道——页面严禁虚构上述内容。

## Product Principles

1. **可信高于一切**：宁漏勿错；页面上每一句主张都必须有仓库内证据支撑，禁止夸大。
2. **结构即价值**：展示「时间线/大事记/事实观点分离」这些结构本身，而不是空洞的宣传语。
3. **为程序员写作**：中文为主、专有名词保留英文，密度高、无营销废话。
4. **内容即界面**：能直接展示真实库内片段（大事记条目、线索时间线）就不做假示意图。
