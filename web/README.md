# News4Coder Web UI

采用 CNCF 顶级开源项目设计美学的现代 Web 界面。

## 设计特点

- **配色方案**: 参考 Kubernetes、Prometheus 等 CNCF 项目
  - 主色: `#326CE5` (Kubernetes Blue)
  - 辅助色: `#00A3C4` (Cloud Native Teal)
  - 中性色: Slate 灰阶系统

- **字体**: Inter (Google Fonts)
- **设计风格**: 
  - 卡片式布局
  - 微妙阴影和渐变
  - 充足的留白
  - 现代圆角
  - 响应式设计

## 核心亮点

### 数据主权 (Data Sovereignty)

页面重点突出 News4Coder 的**本地优先**理念：

1. **Hero 区域** - 强调"你的数据你做主"，展示 100% 本地存储、0 云端依赖
2. **数据主权专页** - 详细介绍四大核心优势：
   - 完全本地存储
   - 隐私优先设计
   - 数据可导出
   - 永久可用
3. **对比表格** - News4Coder vs 云端阅读器的直观对比
4. **特性更新** - 强调本地知识库、本地 RAG 问答

## 文件结构

```
web/
├── index.html          # 主页面（含数据主权板块）
├── css/
│   └── style.css       # 样式文件（含数据主权样式）
├── js/
│   └── main.js         # 交互脚本
├── assets/
│   ├── favicon.svg     # 网站图标
│   └── images/         # 图片资源
└── README.md           # 本文件
```

## 本地预览

```bash
# 使用 Python 简单 HTTP 服务器
cd web && python3 -m http.server 8080

# 或使用 Node.js 的 serve
npx serve web
```

然后访问 http://localhost:8080

## 响应式断点

- Desktop: > 1024px
- Tablet: 768px - 1024px
- Mobile: < 768px
- Small Mobile: < 480px

## 更新日志

### 2024-04-02
- 新增「数据主权」专页板块
- Hero 区域增加数据主权徽章
- 添加产品对比表格
- 更新 CTA 区域，突出本地优先特性
- 页脚增加数据主权标语
