# Gin 云托管部署指南

本指南详细介绍如何将 Gin 应用部署到 CloudBase 云托管服务。

> **📋 前置要求**：如果您还没有创建 Gin 项目，请先阅读 [Gin 项目创建指南](./project-setup.md)。

## 📋 目录导航

- [部署特性](#部署特性)
- [准备部署文件](#准备部署文件)
- [项目结构](#项目结构)
- [部署步骤](#部署步骤)
- [访问应用](#访问应用)
- [常见问题](#常见问题)
- [最佳实践](#最佳实践)
- [高级配置](#高级配置)

---

## 部署特性

云托管适合以下场景：

- **企业级应用**：复杂的 Web 应用和管理系统
- **高并发**：需要处理大量并发请求
- **自定义环境**：需要特定的运行时环境
- **微服务架构**：容器化部署和管理

### 技术特点

| 特性 | 说明 |
|------|------|
| **计费方式** | 按资源使用量（CPU/内存） |
| **启动方式** | 持续运行 |
| **端口配置** | 可自定义端口（默认 8080） |
| **扩缩容** | 支持自动扩缩容配置 |
| **Go 环境** | 完全自定义 Go 环境 |

## 准备部署文件

### 1. 创建 Dockerfile

创建 `Dockerfile` 文件：

```dockerfile
# 使用官方 Go 运行时作为基础镜像
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装 git（某些 Go 模块可能需要）
RUN apk add --no-cache git

# 复制 go mod 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 编译应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# 使用轻量级镜像作为运行环境
FROM alpine:latest

# 安装 ca-certificates（HTTPS 请求需要）
RUN apk --no-cache add ca-certificates

# 创建工作目录
WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .

# 暴露端口
EXPOSE 8080

# 设置环境变量
ENV GIN_MODE=release
ENV PORT=8080

# 运行应用
CMD ["./main"]
```

### 2. 创建 .dockerignore

创建 `.dockerignore` 文件：

```dockerignore
# Git
.git
.gitignore

# Docker
Dockerfile
.dockerignore

# 文档
README.md
docs/

# 编译文件
main
*.exe

# 开发工具
.vscode/
.idea/

# 日志文件
*.log

# 临时文件
tmp/
temp/

# 云函数相关
scf_bootstrap

# 操作系统
.DS_Store
Thumbs.db
```

## 项目结构

部署前的项目结构：

```
cloudrun-gin/
├── main.go                 # 主应用文件
├── go.mod                  # Go 模块文件
├── go.sum                  # 依赖校验文件
├── controllers/            # 控制器目录
│   └── user.go
├── models/                 # 数据模型目录
│   └── user.go
├── Dockerfile             # 🔑 容器配置文件
├── .dockerignore          # Docker 忽略文件
└── scf_bootstrap          # 云函数启动脚本（云托管不需要）
```

## 部署步骤

### 方式一：CloudBase 控制台部署

1. 登录 [CloudBase 控制台](https://console.cloud.tencent.com/tcb)
2. 选择环境，进入「云托管」页面
3. 点击「新建服务」
4. 填写服务名称（如：`cloudrun-gin-service`）
5. 选择「本地代码」上传方式
6. 上传项目文件夹
7. 系统自动检测 Dockerfile 并构建镜像
8. 配置服务参数：
   - CPU：0.25 核
   - 内存：0.5 GB
   - 端口：8080
   - 环境变量：`GIN_MODE=release`
9. 点击「完成」开始部署

### 方式二：CloudBase CLI 部署

```bash
# 安装 CloudBase CLI
npm install -g @cloudbase/cli

# 登录
tcb login

# 初始化配置
tcb init

# 创建 cloudbaserc.json 配置文件
cat > cloudbaserc.json << EOF
{
  "envId": "your-env-id",
  "framework": {
    "name": "gin",
    "plugins": {
      "container": {
        "use": "@cloudbase/framework-plugin-container",
        "inputs": {
          "serviceName": "cloudrun-gin-service",
          "servicePath": "/cloudrun-gin",
          "dockerfilePath": "./Dockerfile",
          "buildDir": "./",
          "containerPort": 8080,
          "minNum": 0,
          "maxNum": 10,
          "cpu": 0.25,
          "mem": 0.5,
          "envVariables": {
            "GIN_MODE": "release",
            "PORT": "8080"
          }
        }
      }
    }
  }
}
EOF

# 部署到云托管
tcb framework deploy
```

### 方式三：Docker 本地构建

```bash
# 本地构建镜像
docker build -t cloudrun-gin .

# 测试镜像
docker run -p 8080:8080 cloudrun-gin

# 推送到镜像仓库（可选）
docker tag cloudrun-gin your-registry/cloudrun-gin:latest
docker push your-registry/cloudrun-gin:latest
```

## 访问应用

部署成功后，您将获得一个访问链接：

```
https://cloudrun-gin-service-xxx-xxx.ap-shanghai.run.tcloudbase.com
```

### 测试接口

```bash
# 替换为您的实际域名
export API_URL="https://cloudrun-gin-service-xxx-xxx.ap-shanghai.run.tcloudbase.com"

# 测试基础接口
curl $API_URL/
curl $API_URL/health

# 测试用户 API
curl $API_URL/api/users
curl "$API_URL/api/users?page=1&limit=2"
curl $API_URL/api/users/1

# 创建用户
curl -X POST $API_URL/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"云托管用户","email":"cloudrun@example.com"}'
```

## 常见问题

### 1. 构建失败

**问题**：Docker 构建过程中出现错误

**解决方案**：
```bash
# 检查 Dockerfile 语法
docker build --no-cache -t cloudrun-gin .

# 查看构建日志
docker build -t cloudrun-gin . --progress=plain
```

### 2. 端口配置问题

**问题**：服务无法正常访问

**解决方案**：
- 确保 Dockerfile 中 EXPOSE 端口与应用监听端口一致
- 检查环境变量 PORT 配置

```dockerfile
# Dockerfile
EXPOSE 8080
ENV PORT=8080
```

### 3. 依赖下载失败

**问题**：Go 模块下载超时

**解决方案**：
```dockerfile
# 在 Dockerfile 中设置 Go 代理
ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=sum.golang.google.cn
```

### 4. 内存不足

**问题**：应用运行时内存不足

**解决方案**：
- 增加容器内存配置
- 优化 Go 应用内存使用

### 5. 健康检查失败

**问题**：容器启动后健康检查失败

**解决方案**：
```go
// 添加健康检查端点
router.GET("/health", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "status": "healthy",
        "timestamp": time.Now().Format(time.RFC3339),
    })
})
```

## 最佳实践

### 1. 多阶段构建优化

```dockerfile
# 优化的 Dockerfile
FROM golang:1.21-alpine AS builder

# 设置必要的环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /build

# 复制并下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制源码并编译
COPY . .
RUN go build -ldflags="-s -w" -o app .

# 运行阶段
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

COPY --from=builder /build/app .

EXPOSE 8080
CMD ["./app"]
```

### 2. 环境变量管理

```go
// config/config.go
package config

import (
    "os"
    "strconv"
)

type Config struct {
    Port        string
    GinMode     string
    DatabaseURL string
    RedisURL    string
    LogLevel    string
}

func Load() *Config {
    return &Config{
        Port:        getEnv("PORT", "8080"),
        GinMode:     getEnv("GIN_MODE", "release"),
        DatabaseURL: getEnv("DATABASE_URL", ""),
        RedisURL:    getEnv("REDIS_URL", ""),
        LogLevel:    getEnv("LOG_LEVEL", "info"),
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
```

### 3. 优雅关闭

```go
// main.go
package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/gin-gonic/gin"
)

func main() {
    router := gin.Default()
    
    // 配置路由...
    
    srv := &http.Server{
        Addr:    ":" + os.Getenv("PORT"),
        Handler: router,
    }

    // 启动服务器
    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("listen: %s\n", err)
        }
    }()

    // 等待中断信号
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    log.Println("Shutting down server...")

    // 优雅关闭
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    if err := srv.Shutdown(ctx); err != nil {
        log.Fatal("Server forced to shutdown:", err)
    }

    log.Println("Server exiting")
}
```

### 4. 日志配置

```go
// middleware/logger.go
package middleware

import (
    "log"
    "os"
    "time"

    "github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
    return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
        return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
            param.ClientIP,
            param.TimeStamp.Format(time.RFC1123),
            param.Method,
            param.Path,
            param.Request.Proto,
            param.StatusCode,
            param.Latency,
            param.Request.UserAgent(),
            param.ErrorMessage,
        )
    })
}
```

## 高级配置

### 1. 自动扩缩容

```json
{
  "framework": {
    "plugins": {
      "container": {
        "inputs": {
          "serviceName": "cloudrun-gin-service",
          "minNum": 1,
          "maxNum": 10,
          "cpu": 0.5,
          "mem": 1,
          "policyType": "cpu",
          "policyThreshold": 70,
          "envVariables": {
            "GIN_MODE": "release"
          }
        }
      }
    }
  }
}
```

### 2. 数据库连接

```go
// database/db.go
package database

import (
    "database/sql"
    "log"
    "os"
    
    _ "github.com/lib/pq"
)

var DB *sql.DB

func Init() {
    var err error
    dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        log.Fatal("DATABASE_URL environment variable is required")
    }

    DB, err = sql.Open("postgres", dbURL)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    if err = DB.Ping(); err != nil {
        log.Fatal("Failed to ping database:", err)
    }

    log.Println("Database connected successfully")
}
```

### 3. 缓存集成

```go
// cache/redis.go
package cache

import (
    "context"
    "log"
    "os"
    "time"

    "github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

func Init() {
    redisURL := os.Getenv("REDIS_URL")
    if redisURL == "" {
        log.Println("REDIS_URL not set, skipping Redis initialization")
        return
    }

    opt, err := redis.ParseURL(redisURL)
    if err != nil {
        log.Fatal("Failed to parse Redis URL:", err)
    }

    RedisClient = redis.NewClient(opt)

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    _, err = RedisClient.Ping(ctx).Result()
    if err != nil {
        log.Fatal("Failed to connect to Redis:", err)
    }

    log.Println("Redis connected successfully")
}
```

### 4. 监控和健康检查

```go
// middleware/health.go
package middleware

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
)

func HealthCheck() gin.HandlerFunc {
    return func(c *gin.Context) {
        if c.Request.URL.Path == "/health" {
            c.JSON(http.StatusOK, gin.H{
                "status":    "healthy",
                "timestamp": time.Now().Format(time.RFC3339),
                "service":   "cloudrun-gin",
                "version":   "1.0.0",
            })
            c.Abort()
            return
        }
        c.Next()
    }
}
```

---

**提示**：
- 使用多阶段构建减少镜像大小
- 配置合适的资源限制
- 实现优雅关闭机制
- 添加健康检查端点
- 使用环境变量管理配置