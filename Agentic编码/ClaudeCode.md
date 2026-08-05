---
线索: ClaudeCode
主题: Agentic编码
别名: [Claude Code, Anthropic 编码助手]
状态: 活跃
创建: 2026-07-21
更新: 2026-08-03
关键角色: [Anthropic]
---

# ClaudeCode

> Anthropic 推出的 agentic 命令行编码助手，在终端内自主理解代码库、执行多步编程任务。

## 概述

Claude Code 是 Anthropic 于 2025 年 2 月推出的 agentic 编码工具，定位为「终端内的自主编程代理」。区别于传统的代码补全或对话式助手，它直接在命令行运行，能浏览整个代码库、运行命令、编辑多文件、执行测试，并围绕「任务」而非「单次问答」组织交互。

它是 Agentic Coding 这一细分赛道最具代表性的产品之一，其演进（研究预览 → GA → SDK → Subagents/Skills/Hooks/Plugins 生态）清晰折射出「AI 编码从助手走向自主工程师」的产业趋势。

## 时间线

### 2026-07

- **2026-07-28** · [MCP 2026-07-28 规范发布](https://blog.modelcontextprotocol.io/posts/2026-07-28/)
  MCP 转向 stateless、可缓存、可路由、可全球扩展的 Web 式架构；与旧版本 wire-incompatible；Anthropic 在全 Claude 产品线推进支持。→ [[智能体平台/MCP协议]]

### 2026-06

- **2026-06-29 至 07-03** · Claude Code Week 27 更新
  内置 Explore agent 改为继承主会话模型（上限 Opus）而非 Haiku；background agents 持续更新。

### 2026-05

- **2026-05** · Code with Claude 2026 大会
  旧金山举办，发布 managed offerings 与 Claude API 平台相关公告。

### 2026-03

- **2026-03** · Claude Code 月度更新
  新增 computer use、auto mode、remote control、scheduled tasks、visuals 等。

### 2025-11

- **2025-11-24** · [Introducing Claude Opus 4.5](https://www.anthropic.com/news/claude-opus-4-5)
  Anthropic 发布 Claude Opus 4.5，称其为当时全球最强的编码/Agent/计算机使用模型，定价 $5/$25 每百万 tokens。Claude Code 随之获得更强的长程任务执行能力。
  > 官方表述：在编码、agents、computer use 上达到业界最佳。

### 2025-09

- **2025-09** · Claude Code SDK 更名为 Claude Agent SDK
  原 6 月发布的 Claude Code SDK 更名为 Claude Agent SDK，强调其能力已超出「编码」场景，扩展到通用 Agent 编排。提供 Python 与 TypeScript 实现，内嵌 CLI 二进制并支持 subagents。

### 2025-06

- **2025-06-16** · Claude Code SDK 发布
  Anthropic 推出 Claude Code SDK，面向「自定义 Agent」与集成场景，使开发者能在自己的工作流中嵌入 Claude Code 的 agentic 能力。

### 2025-05

- **2025-05-22** · [Introducing Claude 4](https://www.anthropic.com/news/claude-4)
  随 Claude 4 系列发布，Claude Code 进入正式可用阶段，成为 Anthropic 在编码与 Agent 赛道的旗舰入口。

### 2025-02

- **2025-02-24** · Claude Code 研究预览版发布
  随 Claude 3.7 Sonnet（Anthropic 首个混合推理模型）一同推出研究预览版，定位为终端内 agentic 命令行编码工具，强调「简洁与直接的终端内效用」。

## 分析

Claude Code 2025 全年迭代 176 次更新，演进脉络可归纳为三条主线：

1. **从助手到自主工程师**：起步是「终端里能跑工具的助手」，到年底已具备长程自主任务执行、上下文管理（Microcompact 延长会话）、多 Agent 协作能力。产品定位从「辅助编码」升级为「承担工程任务」。

2. **能力外溢为平台**：SDK 的发布与更名（Code SDK → Agent SDK）是关键信号——Anthropic 把 Claude Code 的 agentic 内核抽成可嵌入能力，意图让其成为通用 Agent 底座，而非只是一个编码工具。

3. **生态化**：Subagents（支持 @-mention 与模型选择）、Skills、Hooks、Plugins 等机制陆续加入，标志其从单一产品走向可扩展平台，社区围绕它生长出 subagents 合集、教程生态。

竞争格局上，Agentic Coding 赛道（Cursor、Cline、Windsurf、GitHub Copilot Workspace 等）高度拥挤，Claude Code 的差异化在于「终端原生 + Anthropic 自研最强编码模型 + SDK 平台化」三者结合。

## 关联线索

- [[基础模型/Claude]]
- [[智能体平台/ClaudeAgentSDK]]
