# Gin HTTP 云函数部署指南

本指南详细介绍如何将 Gin 应用部署到 CloudBase HTTP 云函数。

> **📋 前置要求**：如果您还没有创建 Gin 项目，请先阅读 [Gin 项目创建指南](./project-setup.md)。

## 📋 目录导航

- [部署特性](#部署特性)
- [准备部署文件](#准备部署文件)
- [项目结构](#项目结构)
- [部署步骤](#部署步骤)
- [访问应用](#访问应用)
- [常见问题](#常见问题)
- [最佳实践](#最佳实践)
- [性能优化](#性能优化)

---

## 部署特性

HTTP 云函数适合以下场景：

- **轻量级应用**：API 服务、微服务
- **间歇性访问**：不需要持续运行的应用
- **成本敏感**：按请求次数和执行时间计费
- **快速部署**：无需容器化配置

### 技术特点

| 特性 | 说明 |
|------|------|
| **计费方式** | 按请求次数和执行时间 |
| **启动方式** | 冷启动，按需启动 |
| **端口要求** | 固定 9000 端口 |
| **扩缩容** | 自动按请求扩缩 |
| **Go 环境** | 预配置 Go 运行时 |

## 准备部署文件

### 1. 编译 Go 应用

```bash
# 设置 Linux 环境编译
export GOOS=linux
export GOARCH=amd64

# 编译生成二进制文件
go build -o main .

# 验证编译结果
ls -la main
```

### 2. 创建启动脚本

创建 `scf_bootstrap` 文件（无扩展名）：

```bash
#!/bin/bash
export PORT=9000
./main
```

为启动脚本添加执行权限：

```bash
chmod +x scf_bootstrap
```

### 3. 项目结构

部署前的项目结构：

```
cloudrun-gin/
├── main                    # 🔑 编译后的二进制文件
├── scf_bootstrap          # 🔑 云函数启动脚本
├── controllers/
│   └── user.go
├── models/
│   └── user.go
├── main.go
├── go.mod
└── go.sum
```

## 部署步骤

### 方式一：CloudBase 控制台部署

1. 登录 [CloudBase 控制台](https://console.cloud.tencent.com/tcb)
2. 选择环境，进入「云函数」页面
3. 点击「新建云函数」
4. 选择「HTTP 云函数」
5. 填写函数名称（如：`cloudrun-gin`）
6. 选择运行环境：`Go 1.x`
7. 上传方式选择「本地上传文件夹」
8. 函数代码选择 `cloudrun-gin` 目录进行上传
9. 点击「完成」等待部署

### 方式二：CloudBase CLI 部署(敬请期待)

### 方式三：本地打包上传

```bash
# 创建部署包
zip -r cloudrun-gin.zip main scf_bootstrap controllers/ models/ go.mod go.sum

# 通过控制台上传 zip 文件
```

## 访问应用

部署成功后，您可以参考[通过 HTTP 访问云函数](https://docs.cloudbase.net/service/access-cloud-function)设置自定义域名访问 HTTP 云函数。

访问地址格式：`https://your-function-url/`

### 测试接口

```bash

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
  -d '{"name":"云函数用户","email":"scf@example.com"}'
```

## 常见问题

### 1. 编译问题

**问题**：本地编译的二进制文件无法在云函数中运行

**解决方案**：
```bash
# 确保使用 Linux 环境编译
export GOOS=linux
export GOARCH=amd64
go build -o main .
```

### 2. 端口配置问题

**问题**：应用无法正常访问

**解决方案**：
- 确保应用监听 9000 端口
- 检查 `scf_bootstrap` 中的 PORT 环境变量

```bash
# scf_bootstrap
#!/bin/bash
export PORT=9000
./main
```

### 3. 文件权限问题

**问题**：启动脚本无执行权限

**解决方案**：
```bash
chmod +x scf_bootstrap
```

### 4. 依赖问题

**问题**：缺少 Go 模块依赖

**解决方案**：
```bash
# 确保包含 go.mod 和 go.sum
go mod tidy
go mod download
```

### 5. 冷启动优化

**问题**：首次访问响应较慢

**解决方案**：
- 减少二进制文件大小
- 优化初始化逻辑
- 使用预热机制

## 最佳实践

### 1. 二进制文件优化

```bash
# 减少二进制文件大小
go build -ldflags="-s -w" -o main .

# 查看文件大小
ls -lh main
```

### 2. 日志配置

```go
// main.go
import (
    "log"
    "os"
)

func init() {
    // 配置日志输出
    log.SetOutput(os.Stdout)
    log.SetFlags(log.LstdFlags | log.Lshortfile)
}
```

### 3. 错误处理

```go
// 统一错误处理中间件
func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        
        if len(c.Errors) > 0 {
            err := c.Errors.Last()
            log.Printf("Request error: %v", err)
            
            c.JSON(http.StatusInternalServerError, gin.H{
                "success": false,
                "message": "Internal server error",
            })
        }
    }
}
```

### 4. 环境变量管理

```go
// 环境配置
type Config struct {
    Port     string
    GinMode  string
    LogLevel string
}

func LoadConfig() *Config {
    return &Config{
        Port:     getEnv("PORT", "9000"),
        GinMode:  getEnv("GIN_MODE", "release"),
        LogLevel: getEnv("LOG_LEVEL", "info"),
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
```

## 性能优化

### 1. 减少冷启动时间

```go
// 预初始化全局变量
var (
    router *gin.Engine
    once   sync.Once
)

func initRouter() {
    once.Do(func() {
        router = gin.New()
        // 配置路由...
    })
}

func main() {
    initRouter()
    router.Run(":" + os.Getenv("PORT"))
}
```

### 2. 内存优化

```go
// 设置合理的内存限制
func init() {
    // 设置 GC 目标百分比
    debug.SetGCPercent(20)
}
```

### 3. 并发控制

```go
// 限制并发数
var semaphore = make(chan struct{}, 100)

func ConcurrencyLimit() gin.HandlerFunc {
    return func(c *gin.Context) {
        select {
        case semaphore <- struct{}{}:
            defer func() { <-semaphore }()
            c.Next()
        default:
            c.JSON(http.StatusTooManyRequests, gin.H{
                "success": false,
                "message": "Too many requests",
            })
            c.Abort()
        }
    }
}
```

### 4. 监控和日志

```go
// 请求监控中间件
func RequestMonitor() gin.HandlerFunc {
    return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
        return fmt.Sprintf("[%s] %s %s %d %s\n",
            param.TimeStamp.Format("2006-01-02 15:04:05"),
            param.Method,
            param.Path,
            param.StatusCode,
            param.Latency,
        )
    })
}
```

---

**提示**：
- 确保二进制文件针对 Linux 环境编译
- 监控函数执行时间和内存使用
- 定期检查日志排查问题
- 考虑使用 CDN 加速静态资源访问