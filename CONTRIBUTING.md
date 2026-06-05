# 贡献指南

感谢你对 News4Coder 感兴趣！我们欢迎各种形式的贡献，包括但不限于：

- 提交 Bug 报告
- 提出新功能建议
- 改进文档
- 提交代码修复或新功能
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

## 开发环境

### 前置要求

- Go 1.25 或更高版本
- Git

### 本地运行

```bash
# 克隆仓库
git clone https://github.com/YOUR_USERNAME/news4coder.git
cd news4coder

# 运行测试
go test ./...

# 编译
go build -o news4coder

# 运行
./news4coder --help
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

如果你想为 News4Coder 添加一个新的官方技术源，通常需要修改以下文件：

1. `internal/official/registry.go` - 在 `registerDefaultSources()` 中注册源信息
2. `internal/crawler/factory.go` - 在 `NewCrawler()` 中返回对应的采集器
3. `cmd/sources.go`（通常无需修改，会自动读取注册表）
4. 如需专用采集器，在 `internal/crawler/` 下新建实现 `Crawler` 接口的文件

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

- 在 [GitHub Discussions](https://github.com/news4coder/news4coder/discussions) 发起讨论
- 在 [GitHub Issues](https://github.com/news4coder/news4coder/issues) 搜索类似问题

再次感谢你的贡献！🚀
