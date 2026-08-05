---
线索: Agent框架
主题: 智能体平台
别名: [LangGraph, AutoGen, Microsoft Agent Framework, CrewAI, MAF]
状态: 活跃
创建: 2026-08-03
更新: 2026-08-03
关键角色: [Microsoft, LangChain, CrewAI]
---

# Agent框架

> 构建多步骤、多智能体应用的编排框架生态，Agent 时代的「应用服务器」层。

## 概述

Agent 框架是连接大模型能力与实际应用的编排层：负责工具调用、记忆、多步规划、多 agent 协作。2024–2025 年头部框架包括 LangChain/LangGraph（社区生态最大）、Microsoft AutoGen（多 agent 对话）、CrewAI（角色化协作）、Semantic Kernel。

2026 年的关键事件是 Microsoft 将 AutoGen 与 Semantic Kernel 合并为统一的 Microsoft Agent Framework（MAF）并 GA，标志企业级 Agent 编排进入整合期。本线索追踪主流框架的版本演进与格局变化。

## 时间线

### 2026-07

- **2026-07-16** · LangGraph v1.3.14
  持续迭代（LangGraph 1.0 GA 于 2025-10-22），仍是社区主流 agent 编排框架。

### 2026-04

- **2026-04-07** · [AutoGen 进入维护模式](https://learn.microsoft.com/en-us/agent-framework/migration-guide/from-autogen/)
  不再加新功能，仅 bug 修复与安全补丁；提供迁移指南至 MAF。
- **2026-04-02/03** · [Microsoft Agent Framework (MAF) 1.0 GA](https://devblogs.microsoft.com/agent-framework/microsoft-agent-framework-at-build-2026-announce/)
  AutoGen + Semantic Kernel 合并后的统一 SDK 正式发布，在 Build 2026 宣布。

### 2026-02

- **2026-02** · MAF Release Candidate
  API 锁定，面向生产使用。

## 分析

1. **整合是成熟信号**：AutoGen + Semantic Kernel 合并为 MAF，反映 Agent 框架从「百花齐放的研究原型」走向「企业级标准化平台」的拐点，正如当年 Web 框架的收敛。

2. **编排层 vs 工具层**：MCP 标准化了「模型-工具」连接（工具层），Agent 框架竞争的是其上的「多步编排/协作」（编排层）。两层正在分工清晰化。

3. **LangGraph 的护城河**：即便微软官方推 MAF，LangGraph 凭借开源社区惯性、模型无关性与丰富的教程生态，仍是创业团队与黑客的主流选择。

## 关联线索

- [[智能体平台/MCP协议]]
- [[智能体平台/ClaudeAgentSDK]]
- [[Agentic编码/ClaudeCode]]
