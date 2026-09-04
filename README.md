# AI News Database（AI 新闻数据库）

一个专门用于持续沉淀 AI 高质量新闻和重点事件的 Markdown 数据库。按主题组织、以线索为单位、用时间线累积，最终呈现深度分析与演进脉络。

<!-- RECENT:START -->
- **2026-09-03** · [OpenAI 发布 GPT-6 Astra 并公布安全概览](https://aihot.virxact.com/items/cmtm02uu60t2arow5vd3nq4gi) — 首个触发 Preparedness Framework「Critical」网络安全阈值的模型，ARC-AGI-3 达 99.9%，Daybreak Access 受限开放。
- **2026-09-03** · [NVIDIA 宣布以 129.3 亿美元收购 Hugging Face](https://aihot.virxact.com/items/cmtli5yd109u4row52i1xg9j4) — 芯片巨头把开源模型社区收入麾下，黄仁勋承诺保持开放。
- **2026-09-03** · [METR 发布 OpenAI/Hugging Face 智能体攻击事件独立调查报告](https://aihot.virxact.com/items/cmtl25m9c0e89roalh6qci0r5) — 夏季智能体逃逸事故潮首次获得独立验证与系统性梳理。
- **2026-09-02** · [Anthropic 发布 Claude Fable 5.1 与 Claude Mythos 5.1](https://aihot.virxact.com/items/cmtjjkmd800r4roe4wpq221bc) — 双旗舰点版本迭代，Fable 5.1 同步上线 Claude Code 并降价 75%。
- **2026-09-02** · [Google DeepMind 发布 Gemini 3.8 Flash 与 3.8 Flash Cyber](https://aihot.virxact.com/items/cmtkbdbti01n8roz5k5kt1g98) — 三周内再度点版本迭代，「Flash + Cyber」双轨延续。
- **2026-09-02** · [Qwen3.8-Max-0902 登顶 Code Arena](https://aihot.virxact.com/items/cmtjgq5z60c3vroq546ccqp7r) — 国产模型首登编码竞技场榜首，$5/MToken 处 Pareto 前沿。
- **2026-09-02** · [OpenAI 因 Tumbler Ridge 枪击案面临 30 起新诉讼](https://aihot.virxact.com/items/cmtkaa7gs01v5romp7lhz3mwj) — AI 产品责任诉讼进入规模化阶段。
<!-- RECENT:END -->

## 这是什么

- **Markdown 数据库为主体**：根目录是按 AI 全景主题组织的 Markdown 文件，git 跟踪、人可读、可长期沉淀。
- **一条线索 = 一个时间线文档**：每个 `.md` 持续追加某条故事线（产品、技术、趋势）的相关新闻，越积越厚。
- **手动整理为主**：当前阶段不依赖自动化抓取，新闻由人工筛选、整理后录入，保证质量。
- **辅助工具统一在 [`tools/`](tools/)**：CLI、Web 界面、浏览器扩展均在 `tools/cmd/` 下作为自包含 Go module 维护，仅作可选辅助，不是核心入口。

## 目录结构

```
ai-news-database/
├── <主题文件夹>/          # 17 个 AI 全景主题，见 _topics.md
│   ├── _index.md          # 主题说明 + 线索列表 + 候选线索
│   └── <线索>.md          # 线索·时间线文件
├── _topics.md             # 中心主题索引
├── _2026大事记.md         # 2026 年度重点事件索引（按月归档）
├── _2025大事记.md         # 2025 年度重点事件索引（回溯建档）
├── tools/                 # 辅助工具（CLI + Web 界面 + 浏览器扩展 + 脚本）
│   ├── cmd/               # 自包含 Go module
│   └── scripts/           # validate.py 格式校验 / build_feed.py 订阅源生成
├── docs/                  # 项目文档 + 线索模板 + 路线图
└── README.md
```

## 如何阅读

1. 从 [`_topics.md`](_topics.md) 选一个主题，或从 [`_2026大事记.md`](_2026大事记.md) / [`_2025大事记.md`](_2025大事记.md) 按月份浏览年度重点事件。
2. 进入主题文件夹，读 `_index.md` 看该主题下有哪些线索、哪些候选线索待建。
3. 打开任一线索 `.md`，「时间线」章节按时间倒序排列（最新在上），往下看即追溯演进。

## 如何订阅

RSS：`https://standup-coder.github.io/ai-news-database/feed.xml`（每周自动更新；备份地址见 `feed.xml`）。首页：`https://standup-coder.github.io/ai-news-database/`。

## 如何贡献一条新闻

1. 找到新闻所属的线索文件（如 `智能体平台/Claude代码助手.md`）。若无线索，从 [`docs/线索模板.md`](docs/线索模板.md) 复制一个新建。
2. 在「时间线」顶部最近的 `### YYYY-MM` 下插入一条：

   ```
   - **2026-07-21** · [标题](URL)
     一句话摘要。关键数据/影响。
   ```

3. 跨月则新建月份标题。
4. 更新 frontmatter 的 `更新:` 日期。
5. 若属于年度重点事件，同步在 [`_2026大事记.md`](_2026大事记.md) 对应月份下登记一行。
6. （可选）补充「概述」（事实）或「分析」（观点）。

详见 [`CONTRIBUTING.md`](CONTRIBUTING.md)。

## 如何新建一个主题

见 [`_topics.md`](_topics.md) 的「如何新增主题」。

## 许可

- **代码**（`tools/` 等）：[MIT](LICENSE)
- **内容**（主题文件夹、年度大事记、docs 文档）：[CC BY-SA 4.0](https://creativecommons.org/licenses/by-sa/4.0/deed.zh) —— 转载请署名「AI News Database」并以相同方式共享。

## 工具

所有辅助工具统一在 [`tools/`](tools/) 目录下管理（当前阶段以手动整理为主，工具仅作可选辅助）：

- `tools/cmd/`：Go CLI（`ai-news-database`），含本地阅读、筛选、导出等能力
- `tools/cmd/web/`：本地 Web 界面（由 CLI 的 `web` 子命令提供服务）
- `tools/cmd/browser-extension/`：浏览器剪藏扩展

在 `tools/cmd/` 目录下构建：

```bash
cd tools/cmd
go build ./...
go test ./...
```

或从仓库根用转发 Makefile：

```bash
make build
make test
```
