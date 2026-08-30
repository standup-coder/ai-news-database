---
线索: Claude
主题: 基础模型
别名: [Claude 3, Claude 3.5, Claude 4, Anthropic 模型]
状态: 活跃
创建: 2026-07-30
更新: 2026-08-28
关键角色: [Anthropic]
---

# Claude

> Anthropic 的旗舰模型家族，以编码能力与安全对齐著称，是企业与开发者市场的头号玩家。

## 概述

Claude 是 Anthropic（由前 OpenAI 研究副总裁 Dario Amodei 等人于 2021 年创办）的模型家族，按能力分为 Haiku/Sonnet/Opus 三档。2024 年 Claude 3.5 Sonnet 凭编码能力成为开发者首选；2025 年 Claude 4 系列与 Claude Code 联动，确立其在 Agentic 编码市场的统治地位。

Anthropic 以「宪法 AI（Constitutional AI）」等对齐方法著称，其模型发布节奏与安全策略（RSP 负责任扩展政策）是行业安全实践的参照系。

## 时间线

### 2026-08

- **2026-08-11** · [Claude 文本水印机制公布](https://www.anthropic.com/news/claude-text-watermark)
  Anthropic 官宣在 Claude 生成的文本中加入不可见水印，用于标识 AI 生成内容来源；TechCrunch、Business Insider、CNN 跟进报道。

### 2026-07

- **2026-07-24** · [Claude Opus 5 发布](https://www.anthropic.com/news/claude-opus-5)
  相对 Opus 4.8 的阶跃式提升，成为 Max 计划默认；定价 $5/M 输入、$25/M 输出。Axios 报道其性能接近 Anthropic 最强的 Fable 层级。

### 2026-06

- **2026-06-30** · [Claude Sonnet 5 发布；Fable 5 全球重新部署启动](https://www.anthropic.com/news/claude-sonnet-5)
  Sonnet 5 以低价带来「接近 Opus 的编程智能」（$2/M 输入、$10/M 输出）；同日宣布 Fable 5 全球重新部署启动（出口管制已解除）。
- **2026-06-12** · Fable 5 / Mythos 5 访问暂停
  因美国政府突然对该批新模型施加出口管制，Anthropic 暂停访问（Fable 5 实际仅上线约 72 小时）。→ [[AI安全与对齐/AI监管政策]]
- **2026-06-09** · [Claude Fable 5 与 Mythos 5 发布](https://www.anthropic.com/news/claude-fable-5-mythos-5)
  新增「Mythos 级」能力层（高于 Opus 级）；Fable 5 面向长程任务。这是 Claude 产品线首次引入 Opus 之上的新档位。

### 2026-05

- **2026-05-28** · [Claude Opus 4.8 发布](https://www.anthropic.com/news/claude-opus-4-8)
  主打「诚实与可靠性」，配套动态工作流编排器；同价（$5/$25）、知识截止 2026-01，成为未指定模型时的默认。

### 2026-04

- **2026-04-16** · [Claude Opus 4.7 发布](https://www.anthropic.com/news/claude-opus-4-7)
  引入「xhigh」努力等级与更高分辨率视觉处理；$5/M 输入、$25/M 输出；知识截止 2026-01。

### 2026-02

- **2026-02-17** · [Claude Sonnet 4.6 发布](https://www.cnbc.com/2026/02/17/anthropic-ai-claude-sonnet-4-6-default-free-pro.html)
  编程与计算机交互大幅改进，面向免费/Pro 用户默认提供。
- **2026-02-05** · [Claude Opus 4.6 发布](https://www.anthropic.com/news/claude-opus-4-6)
  引入 100 万 token（beta）上下文窗口与「并行自主工作」能力。与 OpenAI GPT-5.3-Codex 同日（相隔约 10 分钟）先后发布。

### 2025-11

- **2025-11** · [Claude Opus 4.5 发布](https://www.anthropic.com/news/claude-opus-4-5)
  宣称在编码、Agent、计算机使用上业界最佳，定价较前代大幅下调（$5/$25 每百万 tokens），旗舰能力开始价格战。

### 2025-05

- **2025-05** · [Claude 4 发布](https://www.anthropic.com/news/claude-4)
  Opus 4 与 Sonnet 4 上线，主打长程 Agent 任务（可连续工作数小时）；Claude Code 同期 GA。Opus 4 首次触发 Anthropic ASL-3 安全等级部署措施。

### 2024-10

- **2024-10** · [Computer Use 能力发布](https://www.anthropic.com/news/3-5-models-and-computer-use)
  升级版 Claude 3.5 Sonnet 可直接操作计算机（看屏幕、移动光标、点击输入），Anthropic 成为首个开放「计算机使用」API 的前沿实验室。

### 2024-06

- **2024-06** · [Claude 3.5 Sonnet 发布](https://www.anthropic.com/news/claude-3-5-sonnet)
  以中档价格超越 GPT-4o 的编码与推理表现，配合 Artifacts 功能成为开发者社区的默认选择，奠定 Anthropic 在编码赛道的口碑。

### 2024-03

- **2024-03** · [Claude 3 家族发布](https://www.anthropic.com/news/claude-3-family)
  Haiku/Sonnet/Opus 三档定位确立，Opus 在多项基准首次全面超越 GPT-4，Anthropic 从「安全实验室」转身为「性能竞争者」。

## 分析

1. **编码即楔子**：Anthropic 没有和 OpenAI 拼消费级入口，而是押注编码/Agent 这一高价值 B 端场景。Claude 3.5 Sonnet → Claude Code → Agent SDK 的路径使其 API 收入结构中编码占比极高，Cursor、GitHub Copilot 等均为其大客户。

2. **安全与商业的平衡术**：RSP、ASL 分级、可解释性研究让 Anthropic 保持「负责任实验室」人设，但发布节奏丝毫不慢——安全叙事本身已成为其对企业客户的销售卖点。

3. **三档产品线的定价学**：Haiku/Sonnet/Opus 的分层 + Opus 4.5 降价，反映前沿模型商品化速度：旗舰能力的溢价窗口正在缩短。

## 关联线索

- [[基础模型/GPT系列]]
- [[基础模型/Gemini]]
- [[Agentic编码/ClaudeCode]]
- [[AI安全与对齐/对齐与可解释性研究]]
