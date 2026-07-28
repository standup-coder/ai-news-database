# 贡献指南

感谢你对 AI News Database 感兴趣！我们欢迎各种形式的贡献，包括但不限于：

- 提交 Bug 报告
- 提出新功能建议
- 改进文档
- 提交代码修复或新功能
- **贡献 AI 新闻线索/时间线**（本仓库当前主体，详见下文「内容库贡献」）
- 分享使用经验
- 帮助回答其他用户的问题

---

## 如何开始

1. **Fork 本仓库** 到你的 GitHub 账号下
2. **Clone 你的 Fork** 到本地
3. 创建一个新的分支：`git checkout -b feature/your-feature-name`
4. 进行修改并提交
5. 推送分支到你的 Fork：`git push origin feature/your-feature-name`
6. 在原始仓库中发起一个 **Pull Request**

---

## 内容库贡献

本仓库现在以 **AI 新闻内容库**为主体（根目录按主题组织的中文 md 文件夹）。欢迎你以「线索」「新闻条目」的方式参与贡献。

> 完整主题清单见 [`_topics.md`](_topics.md)；格式规范详见 `docs/superpowers/specs/2026-07-21-content-restructure-design.md` §5。

### A. 新建一条线索

1. 从 [`docs/线索模板.md`](docs/线索模板.md) 复制到对应主题文件夹，重命名为线索名（中文简称为主，专有名词保留英文，如 `智能体平台/Claude代码助手.md`）。
2. 填写 YAML frontmatter 字段：

   | 字段 | 说明 |
   |---|---|
   | `线索` | 线索中文名 |
   | `主题` | 所属主题文件夹名 |
   | `别名` | 英文/别名列表 |
   | `状态` | `活跃` \| `已完结` \| `观察中` |
   | `创建` / `更新` | 日期，格式 `YYYY-MM-DD` |
   | `关键角色` | 相关公司/机构/人名 |

3. 在该主题的 `_index.md`「线索列表」下追加一行登记：

   ```markdown
   - [线索名](线索名.md) — 活跃
   ```

### B. 向已有线索追加一条新闻

在「时间线」章节顶部最近的 `### YYYY-MM` 标题下插入（最新在上）：

```markdown
- **2026-07-15** · [标题](URL)
  一句话摘要。关键数据/影响。
```

- 跨月则新建一个 `### YYYY-MM` 月份标题。
- 同步更新 frontmatter 的 `更新:` 日期。

### C. 内容原则

- **概述 vs 分析分离**：frontmatter 之后的「概述」章节只写事实（是什么、为什么重要），「分析」章节才写观点判断（趋势、格局）。两者分开，避免事实与观点混淆。
- **时间线倒序**：最新在上，按 `年-月` 分组，月内按日期倒序。
- **双链**：用 `[[主题/线索名]]` 表达线索间关联（如 `[[基础模型/GPT]]`）。
- **专有名词保留英文**：Claude、GPT、Sora、Agentic 等不翻译。

### D. 引用参考

- 线索模板：[`docs/线索模板.md`](docs/线索模板.md)
- 完整主题清单：[`_topics.md`](_topics.md)
- 完整格式规范：`docs/superpowers/specs/2026-07-21-content-restructure-design.md` §5

---

## 开发环境

### 前置要求

- Go 1.25 或更高版本
- Git

### 本地运行

> 注：Go CLI 代码已迁移至 `tools/cmd/`，仓库根是 AI 新闻内容库（仓库主体）。
> 在仓库根运行 `make` 会自动转发到 `tools/cmd`。

```bash
# 克隆仓库
git clone https://github.com/standup-coder/ai-news-database.git
cd ai-news-database

# 运行测试（从仓库根用 make 转发到 tools/cmd）
make test
# 或直接进入 tools/cmd 操作
cd tools/cmd && go test ./...

# 编译（产物：tools/cmd/ai-news-database）
make build

# 运行
./tools/cmd/ai-news-database --help
```

---

## 提交 Issue

### Bug 报告

如果你发现了 Bug，请在提交 Issue 前检查是否已有类似报告。新建 Issue 时请包含以下信息：

- **操作系统和版本**（如 macOS 14.2, Ubuntu 22.04）
- **Go 版本**（`go version` 输出）
- **复现步骤**：尽可能详细
- **预期行为** vs **实际行为**
- **错误日志或截图**
- **相关配置文件**（如有，注意隐去 API Key）

### 功能建议

- 清晰描述你希望添加的功能
- 说明这个功能解决了什么痛点
- 如果可能，提供一个使用场景示例

---

## 代码规范

### Go 代码

- 使用 `gofmt` 格式化代码
- 遵循 [Effective Go](https://go.dev/doc/effective_go) 规范
- 为新功能添加单元测试（如适用）
- 保持现有代码风格一致
- 运行 `make quality` 确保所有质量检查通过后再提交

### 测试要求

| 测试类型 | 要求 | 运行命令 |
|----------|------|----------|
| 单元测试 | 必须 | `go test ./...` |
| 覆盖率 | ≥60% | `make test-coverage` |
| 基准测试 | 推荐 | `make benchmark` |

### 提交信息

我们建议使用清晰的提交信息格式：

```
<type>: <subject>

<body>
```

**Type 说明**：

- `feat`: 新功能
- `fix`: 修复 Bug
- `docs`: 文档更新
- `style`: 代码格式调整（不影响功能）
- `refactor`: 代码重构
- `test`: 测试相关
- `chore`: 构建过程或辅助工具的变动

**示例**：

```
feat: add support for Dev.to official source

- Implement DevToCrawler using RSS feed
- Register source in official registry
- Add unit tests for crawler
```

### 质量检查清单

提交 PR 前确保：

- [ ] `make fmt` 通过
- [ ] `go vet ./...` 无错误
- [ ] `make lint` 无错误
- [ ] `make test-coverage` 覆盖率 ≥60%
- [ ] `make security` 无安全警告
- [ ] 所有测试通过 `go test ./... -race`

---

## 添加新采集源

如果你想为 AI News Database 添加一个新的官方技术源，通常需要修改以下文件（注意路径在 `tools/cmd/` 下）：

1. `tools/cmd/internal/official/registry.go` - 在 `registerDefaultSources()` 中注册源信息
2. `tools/cmd/internal/crawler/factory.go` - 在 `NewCrawler()` 中返回对应的采集器
3. `tools/cmd/cmd/sources.go`（通常无需修改，会自动读取注册表）
4. 如需专用采集器，在 `tools/cmd/internal/crawler/` 下新建实现 `Crawler` 接口的文件

**采集源入选标准**：

- 内容质量高，受众为程序员/技术人员
- 允许公开访问（不要求登录即可阅读标题和摘要）
- 更新频率稳定
- 有清晰的 URL 结构或 API

---

## 文档贡献

文档和代码同等重要！如果你发现文档有：

- 拼写或语法错误
- 过时的安装步骤
- 缺失的功能说明
- 不够清晰的示例

欢迎直接提交 PR 修复。

---

## 行为准则

- 保持友善和尊重
- 接受建设性的批评
- 关注什么对社区最有利
- 对其他社区成员表示同理心

---

## 获取帮助

- 在 [GitHub Discussions](https://github.com/standup-coder/ai-news-database/discussions) 发起讨论
- 在 [GitHub Issues](https://github.com/standup-coder/ai-news-database/issues) 搜索类似问题

再次感谢你的贡献！🚀
