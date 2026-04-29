# dddd-gui

[![License](https://img.shields.io/github/license/renbon-wu/dddd-gui)](https://github.com/renbon-wu/dddd-gui/blob/main/LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/renbon-wu/dddd-gui)](https://github.com/renbon-wu/dddd-gui)
[![Release](https://img.shields.io/github/release/renbon-wu/dddd-gui)](https://github.com/renbon-wu/dddd-gui/releases)
[![Stars](https://img.shields.io/github/stars/renbon-wu/dddd-gui)](https://github.com/renbon-wu/dddd-gui)

基于 [dddd](https://github.com/SleepingBag945/dddd) 二次开发的漏洞扫描工具，带 GUI 界面。

## 特点

- 🖥️ **图形化界面** - 基于 Wails + Vue 3 + Element Plus 的现代化 GUI
- 🔍 **自动识别** - 自动识别输入类型（IP/域名/URL/CIDR），无需手动分类
- 🎯 **指纹识别** - 便于拓展的主动/被动指纹识别，支持复杂逻辑运算
- 🛡️ **漏洞扫描** - Nuclei v3 支持，便于拓展的指纹漏洞映射数据库
- 📊 **HTML 报表** - 高效的 HTML 报表，包含漏洞请求响应
- 🔐 **安全审计** - 审计日志，敏感环境必备
- 🌐 **子域名枚举** - 高效的子域名枚举/爆破，精准的泛解析过滤
- 🔎 **资产测绘** - Hunter、Fofa、Quake 支持

## 截图

> 截图待添加

## 安装

### 从 Release 下载

前往 [Releases](https://github.com/renbon-wu/dddd-gui/releases) 页面下载对应平台的版本：

- `dddd-gui-windows-amd64.exe` - Windows 64位
- `dddd-gui-linux-amd64` - Linux 64位
- `dddd-gui-darwin-amd64` - macOS 64位
- `config.zip` - 配置文件（首次使用需下载）

### 从源码构建

#### 环境要求

- Go 1.22+
- Node.js 18+
- Wails v2.12.0+

#### 构建步骤

```bash
# 克隆仓库
git clone https://github.com/renbon-wu/dddd-gui.git
cd dddd-gui

# 安装前端依赖
cd frontend && npm install && cd ..

# 构建
go build -o dddd-gui .

# 或使用 Wails 构建
wails build
```

## 使用

### GUI 模式

直接运行程序即可启动 GUI 界面：

```bash
./dddd-gui
```

### 功能模块

1. **扫描任务** - 创建和管理扫描任务
2. **目标管理** - 导入和管理扫描目标
3. **指纹管理** - 管理指纹识别规则
4. **POC 管理** - 管理漏洞检测脚本
5. **结果分析** - 查看和分析扫描结果
6. **系统设置** - 配置 API 密钥和字典

## 配置

### API 配置

在 `config/api-config.yaml` 中配置以下 API 密钥：

- Hunter API Key
- Fofa Email & API Key
- Quake API Key

### 字典配置

支持自定义字典文件：
- 子域名爆破字典

## 开发

### 技术栈

- **后端**: Go 1.22, Wails v2, Nuclei v3
- **前端**: Vue 3, Element Plus, Vite 3

### 项目结构

```
dddd-gui/
├── main.go              # 入口文件
├── internal/api/        # API 层
├── common/              # 公共功能模块
├── gopocs/              # Go POC 引擎
├── lib/                 # 第三方库修改版
├── frontend/            # Vue 3 前端
├── structs/             # 数据结构定义
├── utils/               # 工具函数
└── config/              # 配置文件
```

## 更新历史

查看 [CHANGELOG.md](CHANGELOG.md) 了解更新历史。

## 贡献

欢迎提交 Issue 和 Pull Request！请查看 [CONTRIBUTING.md](CONTRIBUTING.md) 了解详情。

## 安全

如果您发现安全漏洞，请查看 [SECURITY.md](SECURITY.md) 了解如何报告。

## 致谢

本项目基于以下开源项目开发：

- [dddd](https://github.com/SleepingBag945/dddd) - 原始项目
- [Wails](https://wails.io/) - GUI 框架
- [Nuclei](https://github.com/projectdiscovery/nuclei) - 漏洞扫描引擎
- [Vue.js](https://vuejs.org/) - 前端框架
- [Element Plus](https://element-plus.org/) - UI 组件库

## 免责声明

本工具仅面向**合法授权**的企业安全建设行为，如您需要测试本工具的可用性，请自行搭建靶机环境。

在使用本工具进行检测时，您应确保该行为符合当地的法律法规，并且已经取得了足够的授权。**请勿对非授权目标进行扫描。**

如您在使用本工具的过程中存在任何非法行为，您需自行承担相应后果，我们将不承担任何法律及连带责任。

## License

MIT License

Copyright (c) 2023 SleepingBag945
Copyright (c) 2024 renbon-wu
