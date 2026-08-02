---
线索: SWE-bench
主题: 评测与基准
别名: [SWE-bench Verified, SWE-agent]
状态: 活跃
创建: 2026-07-30
更新: 2026-07-30
关键角色: [Princeton, OpenAI, Anthropic]
---

# SWE-bench

> 用真实 GitHub issue 考模型：AI 编码能力的黄金基准，两年内从 2% 卷到 80%。

## 概述

SWE-bench 由普林斯顿团队于 2023 年 10 月提出，从 12 个流行 Python 仓库收集 2294 个真实 GitHub issue 及其修复 PR，要求模型在完整代码库上生成能通过测试的补丁。它取代 HumanEval 成为 Agentic 编码能力的事实标准，各前沿模型发布必引其分数。

该线索追踪 SWE-bench 及其变体（Verified/Lite/Multimodal/Pro）的演进与「基准饱和→升级」的循环。

## 时间线

### 2025-05

- **2025-05** · 前沿模型逼近 80%
  Claude 4 系列在 SWE-bench Verified 上报告约 80% 得分（配合并行采样等脚手架），基准接近饱和；社区推出 SWE-bench Pro、Multi-SWE-bench（多语言）等更难版本接棒。

### 2024-08

- **2024-08-13** · [OpenAI 发布 SWE-bench Verified](https://openai.com/index/introducing-swe-bench-verified/)
  OpenAI 与原作者合作，人工审核筛出 500 个题目组成 Verified 子集，剔除描述不清/测试不公的样本——「基准本身需要对齐」成为共识，此后各家报分默认用 Verified。

### 2024-06

- **2024-06** · 商业 Agent 竞速开始
  Factory、Cognition（Devin）、Amazon Q 等纷纷刷榜，SWE-bench 分数成为编码 Agent 创业公司的融资名片，也引发「过拟合基准」「测试集污染」的质疑。

### 2024-04

- **2024-04** · SWE-agent 开源发布
  普林斯顿团队发布 SWE-agent（Agent-Computer Interface 思想），GPT-4 得分从直接提示的个位数提升至 12.5%，证明「脚手架设计」与模型能力同等重要。

### 2023-10

- **2023-10** · SWE-bench 论文发布
  首发评测中最强模型（Claude 2）仅解决 1.96% 的问题，论文标题反问「语言模型能解决真实世界的 GitHub issue 吗」——当时答案接近否定。

## 分析

1. **基准定义赛道**：SWE-bench 把「编码能力」的定义从写函数（HumanEval）拉到改真实代码库，直接塑造了 2024–2025 年 Agentic 编码的研发方向。好基准不只是度量工具，更是研究议程的设置者。

2. **2% → 80% 的两年**：这条曲线是整个 AI 能力进化最直观的缩影——其中既有模型进步（推理模型），也有脚手架进步（Agent 循环）、算力堆叠（并行采样），拆分归因是评测科学的核心难题。

3. **饱和即失效**：当分数逼近上限，基准的区分度与公信力同时崩塌（数据污染疑云、私榜差异）。评测的未来在动态化（实时更新题库）与经济化（以「能完成多少美元的真实工作」计量，如 SWE-Lancer）。

## 关联线索

- [[评测与基准/LMArena]]
- [[Agentic编码/ClaudeCode]]
- [[Agentic编码/Cursor]]
