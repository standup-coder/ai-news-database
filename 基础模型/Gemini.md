---
线索: Gemini
主题: 基础模型
别名: [Google Gemini, Gemini 1.5, Gemini 2.0, Gemini 2.5]
状态: 活跃
创建: 2026-07-30
更新: 2026-08-03
关键角色: [Google, Google DeepMind]
---

# Gemini

> Google DeepMind 的原生多模态旗舰，凭超长上下文与全家桶分发完成从追赶者到并跑者的逆转。

## 概述

Gemini 是 Google DeepMind 的旗舰模型系列，2023 年 12 月发布，定位「原生多模态」（从头以文本/图像/音频/视频联合训练）。其标志性能力是百万级 token 超长上下文（Gemini 1.5）与深度整合 Google 生态（搜索、Workspace、Android）。

在经历 Bard 时期的仓促追赶后，Gemini 2.0/2.5 世代已稳定处于第一梯队，该线索追踪 Google 在基础模型上的技术路线与产品整合。

## 时间线

### 2026-07

- **2026-07-21** · [Gemini 3.6 Flash + 3.5 Flash-Lite + 3.5 Flash Cyber 发布](https://techcrunch.com/2026/07/21/google-releases-three-new-gemini-models-but-no-3-5-pro/)
  三型号同日发布：3.6 Flash（工作马，$1.5/M 输入、$7.5/M 输出，DeepSWE 49% vs 3.5 的 37%）、3.5 Flash-Lite（最低成本）、3.5 Flash Cyber（仅政府专用的网络安全模型）。3.5 Pro 仍未发布。
- **2026-07-17** · Gemini 3.5 Pro 推迟
  DeepMind 放弃原基础、采用更长预训练周期，3.5 Pro 推迟（具体日为半证实单源，Google 官方仅称「仍在测试」）。

### 2026-05

- **2026-05-19** · [Google I/O 2026：Gemini 3.5 系列启动 + Gemini Omni 取代 Veo](https://blog.google/innovation-and-ai/sundar-pichai-io-2026/)
  I/O 2026 发布 Gemini 3.5 Flash（已上线）并预告 3.5 Pro；同日发布 Gemini Omni（取代 Gemini App 中的 Veo），首版 Omni Flash 当日可用，可将任意参考（图像/文本/视频/音频）转为统一输出。2026 年 5 月起 Gemini App 后端默认调用 Omni 而非 Veo。→ [[多模态大模型/视频生成竞赛]]
- **2026-05–06** · Gemini 3 Deep Think 升级
  面向科学/研究/工程的推理升级，面向 AI Ultra 订阅者。

### 2026-04

- **2026-04-02** · [Gemma 4 发布](https://en.wikipedia.org/wiki/Gemini_(language_model))
  开源轻量模型系列迭代。

### 2026-03

- **2026-03-03** · Gemini 3.1 Flash Lite 开发者发布
  通过 Google API 向开发者发布。

### 2026-01

- **2026-01-29** · [Project Genie（Genie 3 世界模型）向公众开放](https://mashable.com/article/google-launches-project-genie-3-how-to-try)
  Genie 3 世界模型通过 Project Genie 网页界面向 Google AI Ultra 美国订阅用户推出，可从文本生成可探索的交互式 3D 环境。→ [[世界模型/Genie]]
- **2026-01** · [AlphaGenome 论文发表于 Nature](https://cen.acs.org/biological-chemistry/genomics/Googles-AlphaGenome-predicts-function-DNA/104/web/2026/01)
  DNA 功能预测模型正式发表于 Nature（模型本身自 2025-06 起非商用可用），解码人类「暗基因组」。→ [[行业应用/AlphaFold与AI科研]]
- **2026-01-05** · Gemini 3 Grounding 计费开始
  Gemini 3 的 Google Search Grounding 开始计费。

### 2025-03

- **2025-03** · [Gemini 2.5 Pro 发布](https://blog.google/technology/google-deepmind/gemini-model-thinking-updates-march-2025/)
  内置思维链的推理模型，在 LMArena 等多个榜单登顶，编码与数学能力大幅跃升，被广泛视为 Google 首次拿到「单点最强」称号。

### 2024-12

- **2024-12** · [Gemini 2.0 发布，进入 Agent 时代](https://blog.google/technology/google-deepmind/google-gemini-ai-update-december-2024/)
  Gemini 2.0 Flash 主打原生工具调用与多模态输出，同步展示 Project Astra（通用助手）与 Project Mariner（浏览器 Agent），Google 正式宣布「Agentic 时代」路线。

### 2024-05

- **2024-05** · AI Overviews 全量上线美国
  I/O 大会宣布 Gemini 驱动的 AI Overviews 进入 Google 搜索主结果页，这是搜索四分之一个世纪以来最大的界面变革（初期因错误答案风波紧急收缩）。

### 2024-02

- **2024-02** · [Gemini 1.5 发布，百万 token 上下文](https://blog.google/technology/ai/google-gemini-next-generation-model-february-2024/)
  MoE 架构 + 100 万（实验室 1000 万）token 上下文窗口，长上下文成为 Gemini 相对 GPT/Claude 的第一个明确技术优势；同月 Bard 更名 Gemini。

### 2023-12

- **2023-12** · [Gemini 1.0 发布](https://blog.google/technology/ai/google-gemini-ai/)
  Google DeepMind 合并后的首个旗舰成果，分 Ultra/Pro/Nano 三档，宣称 MMLU 超越 GPT-4（后因演示视频剪辑争议口碑受损）。

## 分析

1. **从追赶到并跑的教科书**：Bard 的仓促（2023）→ Gemini 1.5 找到长上下文差异化（2024）→ 2.5 登顶（2025），Google 用两年时间证明其研究底蕴（Transformer 发源地）+ TPU 自有算力 + 数据生态的综合优势难以被长期压制。

2. **分发是终极武器**：搜索 AI Overviews、Android 内置、Workspace 集成——Gemini 不需要赢得「主动选择」，只需成为数十亿用户的默认选项。这与 Meta 策略同构，且监管是最大变数（搜索垄断案）。

3. **TPU 垂直整合**：Gemini 全系在自研 TPU 上训练，Google 是唯一不依赖 NVIDIA 的前沿实验室，算力成本优势将在推理模型时代（token 消耗爆炸）被放大。

## 关联线索

- [[基础模型/GPT系列]]
- [[基础模型/Claude]]
- [[世界模型/Genie]]
- [[推理与基础设施/AI芯片竞争]]
