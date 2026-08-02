---
线索: ClaudeAgentSDK
主题: 智能体平台
别名: [Claude Agent SDK, Claude Code SDK]
状态: 活跃
创建: 2026-07-30
更新: 2026-07-30
关键角色: [Anthropic]
---

# ClaudeAgentSDK

> Anthropic 把 Claude Code 的 agentic 内核抽成通用 SDK，意图成为 Agent 开发的默认底座。

## 概述

Claude Agent SDK 前身是 2025 年 6 月发布的 Claude Code SDK，同年 9 月更名，提供 Python 与 TypeScript 实现。它将 Claude Code 内部打磨的 Agent 循环（工具调用、文件操作、上下文管理、子代理编排）封装为可嵌入任意应用的开发框架，支持 subagents、hooks、MCP 工具接入。

该线索追踪 Anthropic 的 Agent 平台化战略及其与 LangChain、OpenAI Agents SDK 等框架的竞争。

## 时间线

### 2025-09

- **2025-09** · Claude Code SDK 更名为 Claude Agent SDK
  更名标志能力定位从「编码」扩展到通用 Agent 编排：官方示例覆盖客服、研究、数据处理等非编码场景，与 Claude Sonnet 4.5 同期发布。

### 2025-06

- **2025-06-16** · Claude Code SDK 发布
  面向「自定义 Agent」与集成场景，开发者可在自己的工作流中复用 Claude Code 的 agentic 能力（工具执行、权限控制、会话管理）。

### 2025-05

- **2025-05** · Claude Code GA 并开放扩展机制
  随 Claude 4 发布，Claude Code 正式可用并陆续开放 Hooks、Subagents 等机制，为 SDK 化铺路。

### 2024-11

- **2024-11-25** · [MCP（Model Context Protocol）开源发布](https://www.anthropic.com/news/model-context-protocol)
  Anthropic 开源模型上下文协议，统一 LLM 连接外部工具/数据源的接口标准，后成为 Agent SDK 的工具生态基座。

## 分析

1. **产品反哺框架**：与 LangChain「先框架后场景」相反，Agent SDK 是从 Claude Code 这个被验证的产品中反向提炼的框架——「在生产中打磨过的 Agent 循环」是其对开发者的核心说服力。

2. **协议 + SDK 双层卡位**：MCP 负责工具连接标准（已被 OpenAI、Google 采纳），Agent SDK 负责编排层。Anthropic 试图同时定义 Agent 生态的「USB 接口」与「主板」。

3. **与模型绑定的风险**：SDK 深度绑定 Claude 模型，在企业多模型策略盛行的背景下，这既是商业护城河，也是被中立框架（LangGraph 等）替代的隐患。

## 关联线索

- [[Agentic编码/ClaudeCode]]
- [[智能体平台/MCP协议]]
- [[基础模型/Claude]]
