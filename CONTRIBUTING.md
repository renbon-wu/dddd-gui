# 贡献指南

感谢您有兴趣为 dddd-gui 做出贡献！

## 如何贡献

### 提交 Issue

如果您发现了 bug 或有新功能建议，请先查看是否已有相关的 Issue。如果没有，请新建一个 Issue 并提供以下信息：

**Bug 报告**：
- 问题描述
- 复现步骤
- 期望行为
- 实际行为
- 环境信息（操作系统、Go 版本等）

**功能建议**：
- 功能描述
- 使用场景
- 可能的实现方式

### 提交 Pull Request

1. Fork 本仓库
2. 创建功能分支：`git checkout -b feature/your-feature`
3. 提交更改：`git commit -m 'feat: add some feature'`
4. 推送分支：`git push origin feature/your-feature`
5. 提交 Pull Request

### 代码规范

#### Go 代码

- 遵循 [Effective Go](https://golang.org/doc/effective_go) 规范
- 使用 `gofmt` 格式化代码
- 添加必要的注释

#### 前端代码

- 遵循 Vue 3 风格指南
- 使用 ESLint 进行代码检查

### 提交信息规范

使用约定式提交：

- `feat:` 新功能
- `fix:` Bug 修复
- `docs:` 文档更新
- `style:` 代码格式调整
- `refactor:` 代码重构
- `perf:` 性能优化
- `test:` 测试相关
- `chore:` 构建/工具相关

## 开发环境搭建

### 环境要求

- Go 1.22+
- Node.js 18+
- Wails v2.12.0+

### 搭建步骤

```bash
# 克隆仓库
git clone https://github.com/renbon-wu/dddd-gui.git
cd dddd-gui

# 安装 Go 依赖
go mod download

# 安装前端依赖
cd frontend && npm install && cd ..

# 开发模式运行
wails dev

# 构建
wails build
```

## 许可证

通过贡献代码，您同意您的代码将按照 MIT 许可证授权。
