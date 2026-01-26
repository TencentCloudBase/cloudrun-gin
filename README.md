# 快速部署 Gin 应用

一个完整的 Gin 应用模板，支持快速部署到 CloudBase 平台。

## 🚀 快速开始

### 前置条件

- [Go 1.19](https://golang.org/dl/) 或更高版本
- 了解基本的 Go 模块使用
- 腾讯云账号并开通了 CloudBase 服务
- 基本的 Go 和 Gin 开发知识

### 创建应用

> **📋 详细指南**：完整的项目创建步骤请参考 [Gin 项目创建指南](./docs/project-setup.md)

```bash
# 快速创建（基础步骤）
mkdir cloudrun-gin && cd cloudrun-gin
go mod init cloudrun-gin
go get github.com/gin-gonic/gin@v1.9.1
```

### 本地测试

```bash
# 启动开发服务器
go run main.go

# 访问应用
open http://localhost:8080
```

## 📦 项目结构

```
cloudrun-gin/
├── main.go                  # 主应用文件
├── go.mod                   # Go 模块文件
├── go.sum                   # 依赖校验文件
├── controllers/             # 控制器目录
│   └── user.go             # 用户控制器
├── models/                  # 数据模型目录
│   └── user.go             # 用户模型
├── scf_bootstrap           # HTTP 云函数启动脚本
├── Dockerfile              # 云托管容器配置
└── .dockerignore           # Docker 忽略文件
```

## 🎯 部署方式

### 部署方式对比

| 特性 | HTTP 云函数 | 云托管 |
|------|------------|--------|
| **计费方式** | 按请求次数和执行时间 | 按资源使用量（CPU/内存） |
| **启动方式** | 冷启动，按需启动 | 持续运行 |
| **适用场景** | API 服务、轻量级应用 | 企业级应用、复杂 Web 应用 |
| **端口要求** | 固定 9000 端口 | 可自定义端口（默认 8080） |
| **扩缩容** | 自动按请求扩缩 | 支持自动扩缩容配置 |
| **Go 环境** | 预配置 Go 运行时 | 完全自定义 Go 环境 |

### 选择部署方式

- **选择 HTTP 云函数**：轻量级 API 服务、间歇性访问、成本敏感
- **选择云托管**：企业级应用、复杂 Web 应用、需要更多控制权

## 📚 详细部署指南

### 🔥 HTTP 云函数部署

适合轻量级应用和 API 服务，按请求计费，冷启动快。

**快速部署步骤：**
1. 编译 Go 二进制文件
2. 创建 `scf_bootstrap` 启动脚本
3. 通过 CloudBase 控制台上传部署

[📖 查看详细的 HTTP 云函数部署指南](./docs/http-function.md)

### 🐳 云托管部署

适合企业级应用，支持更复杂的部署需求，容器化部署。

**快速部署步骤：**
1. 创建 `Dockerfile` 容器配置
2. 配置 `.dockerignore` 文件
3. 通过 CloudBase 控制台或 CLI 部署

[📖 查看详细的云托管部署指南](./docs/cloud-run.md)

## 🔧 API 接口

本模板包含以下 RESTful API 接口：

### 基础接口
```bash
GET /                        # 欢迎页面
GET /health                  # 健康检查
```

### 用户管理
```bash
GET /api/users               # 获取用户列表（支持分页）
GET /api/users/{id}          # 获取单个用户
POST /api/users              # 创建用户
```

### 示例请求

```bash
# 健康检查
curl https://your-app-url/health

# 获取用户列表（分页）
curl "https://your-app-url/api/users?page=1&limit=5"

# 创建新用户
curl -X POST https://your-app-url/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"测试用户","email":"test@example.com"}'
```

## ❓ 常见问题

### 端口配置
- **HTTP 云函数**：必须使用 9000 端口
- **云托管**：推荐使用 8080 端口，支持自定义

### 文件要求
- **HTTP 云函数**：需要编译后的二进制文件和 `scf_bootstrap` 启动脚本
- **云托管**：需要 `Dockerfile` 和 `.dockerignore`

### 数据存储
- 当前使用内存存储（重启后数据丢失）
- 生产环境建议集成数据库（PostgreSQL、MySQL 等）

### 如何选择部署方式？
- **轻量级应用**：选择 HTTP 云函数
- **企业级应用**：选择云托管
- **成本敏感**：选择 HTTP 云函数
- **需要持续运行**：选择云托管

## 🛠️ 开发工具

### 推荐的开发依赖

```bash
# 核心框架
go get github.com/gin-gonic/gin@v1.9.1

# 数据库支持
go get gorm.io/gorm
go get gorm.io/driver/postgres

# 配置管理
go get github.com/spf13/viper

# 日志库
go get github.com/sirupsen/logrus
```

### 环境变量配置

```bash
# 应用配置
export GIN_MODE=release
export PORT=8080

# 数据库配置（可选）
export DATABASE_URL=postgres://user:password@localhost/dbname
```

## 📖 进阶功能

- **中间件支持**：CORS、日志、认证等中间件
- **路由分组**：模块化路由管理
- **数据绑定**：自动 JSON/XML 数据绑定
- **模板引擎**：HTML 模板渲染支持
- **静态文件**：静态资源服务
- **优雅关闭**：支持优雅关闭服务器

## 🔗 相关链接

### 📚 项目文档
- [Gin 项目创建指南](./docs/project-setup.md) - 从零开始创建项目
- [HTTP 云函数部署指南](./docs/http-function.md) - 云函数部署详细步骤
- [云托管部署指南](./docs/cloud-run.md) - 云托管部署详细步骤

### 🌐 官方文档
- [CloudBase 官方文档](https://docs.cloudbase.net/)
- [Gin 官方文档](https://gin-gonic.com/)
- [Go 官方文档](https://golang.org/doc/)

## 📄 许可证

本项目采用 MIT 许可证。详情请查看 [LICENSE](./LICENSE) 文件。

---

**需要帮助？** 

- 查看 [HTTP 云函数部署指南](./docs/http-function.md)
- 查看 [云托管部署指南](./docs/cloud-run.md)
- 访问 [CloudBase 官方文档](https://docs.cloudbase.net/)