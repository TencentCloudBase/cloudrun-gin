# 快速部署 Gin 应用

## 📋 目录导航

- [部署方式对比](#部署方式对比)
- [前置条件](#前置条件)
- [第一步：创建 Gin 应用](#第一步创建-gin-应用)
- [第二步：添加 API 路由](#第二步添加-api-路由)
- [第三步：本地测试](#第三步本地测试)
- [第四步：准备部署文件](#第四步准备部署文件)
- [第五步：项目结构](#第五步项目结构)
- [第六步：部署应用](#第六步部署应用)
- [第七步：访问您的应用](#第七步访问您的应用)
- [常见问题](#常见问题)
- [最佳实践](#最佳实践)
- [进阶功能](#进阶功能)

---

[Gin](https://gin-gonic.com/) 是一个用 Go 语言编写的 HTTP Web 框架。它具有类似 Martini 的 API，但性能比 Martini 快 40 倍。Gin 使用了自定义版本的 HttpRouter，因此它不仅提供了极快的路由，还提供了中间件支持。

本指南介绍如何在 CloudBase 上部署 Gin 应用程序，支持两种部署方式：

- **HTTP 云函数**：适合轻量级应用和 API 服务，按请求计费，冷启动快
- **云托管**：适合企业级应用，支持更复杂的部署需求，容器化部署

## 部署方式对比

| 特性 | HTTP 云函数 | 云托管 |
|------|------------|--------|
| **计费方式** | 按请求次数和执行时间 | 按资源使用量（CPU/内存） |
| **启动方式** | 冷启动，按需启动 | 持续运行 |
| **适用场景** | API 服务、轻量级应用 | 企业级应用、微服务架构 |
| **部署文件** | 需要 `scf_bootstrap` 启动脚本 | 需要 `Dockerfile` 容器配置 |
| **端口要求** | 固定 9000 端口（通过环境变量） | 可自定义端口（默认 8080） |
| **扩缩容** | 自动按请求扩缩 | 支持自动扩缩容配置 |
| **编译方式** | 交叉编译或容器编译 | Docker 多阶段构建 |

## 前置条件

在开始之前，请确保您已经：

- 安装了 [Go 1.19](https://golang.org/dl/) 或更高版本
- 拥有腾讯云账号并开通了 CloudBase 服务
- 了解基本的 Go 语言和 Gin 框架开发知识

## 第一步：创建 Gin 应用

> 💡 **提示**：如果您已经有一个 Gin 应用，可以跳过此步骤。

### 创建项目目录

```bash
mkdir cloudrun-gin
cd cloudrun-gin
```

### 初始化 Go 模块

```bash
go mod init cloudrun-gin
go get -u github.com/gin-gonic/gin
```

### 创建主应用文件

在 `cloudrun-gin` 目录下创建 `main.go` 文件：

```go
package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// 设置 Gin 模式
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// 基础路由
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "欢迎使用 Gin CloudBase 应用！",
			"status":  "running",
		})
	})

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":      "healthy",
			"framework":   "Gin",
			"go_version":  "1.19+",
			"gin_version": gin.Version,
		})
	})

	// 获取端口，支持环境变量
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 启动服务器
	router.Run(":" + port)
}
```

### 本地测试应用

启动应用：

```bash
go run main.go
```

打开浏览器访问 `http://localhost:8080`，您应该能看到 JSON 响应。

## 第二步：添加 API 路由

让我们创建一个 RESTful API 来演示 Gin 的功能。

### 创建用户模型

在项目根目录创建 `models` 目录，并创建 `user.go` 文件：

```go
package models

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type ApiResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
```

### 创建用户控制器

在项目根目录创建 `controllers` 目录，并创建 `user.go` 文件：

```go
package controllers

import (
	"net/http"
	"strconv"
	"sync"

	"cloudrun-gin/models"

	"github.com/gin-gonic/gin"
)

var (
	users   []models.User
	usersMu sync.RWMutex
	nextID  = 1
)

func init() {
	// 初始化测试数据
	users = []models.User{
		{ID: 1, Name: "张三", Email: "zhangsan@example.com"},
		{ID: 2, Name: "李四", Email: "lisi@example.com"},
		{ID: 3, Name: "王五", Email: "wangwu@example.com"},
	}
	nextID = 4
}

// GetUsers 获取用户列表
func GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	usersMu.RLock()
	defer usersMu.RUnlock()

	startIndex := (page - 1) * limit
	endIndex := startIndex + limit

	if startIndex >= len(users) {
		c.JSON(http.StatusOK, models.ApiResponse{
			Success: true,
			Message: "获取成功",
			Data:    []models.User{},
		})
		return
	}

	if endIndex > len(users) {
		endIndex = len(users)
	}

	paginatedUsers := users[startIndex:endIndex]

	c.JSON(http.StatusOK, models.ApiResponse{
		Success: true,
		Message: "获取成功",
		Data: gin.H{
			"total": len(users),
			"page":  page,
			"limit": limit,
			"items": paginatedUsers,
		},
	})
}

// GetUser 根据ID获取用户
func GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Success: false,
			Message: "无效的用户ID",
		})
		return
	}

	usersMu.RLock()
	defer usersMu.RUnlock()

	for _, user := range users {
		if user.ID == id {
			c.JSON(http.StatusOK, models.ApiResponse{
				Success: true,
				Message: "获取成功",
				Data:    user,
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, models.ApiResponse{
		Success: false,
		Message: "用户不存在",
	})
}

// CreateUser 创建用户
func CreateUser(c *gin.Context) {
	var newUser models.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Success: false,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	if newUser.Name == "" || newUser.Email == "" {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Success: false,
			Message: "姓名和邮箱不能为空",
		})
		return
	}

	usersMu.Lock()
	newUser.ID = nextID
	nextID++
	users = append(users, newUser)
	usersMu.Unlock()

	c.JSON(http.StatusCreated, models.ApiResponse{
		Success: true,
		Message: "创建成功",
		Data:    newUser,
	})
}
```

### 更新主应用文件

更新 `main.go` 文件，添加路由和中间件：

```go
package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"cloudrun-gin/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	// 设置 Gin 模式
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// 添加 CORS 中间件
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// 添加日志中间件
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
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
	}))

	// 基础路由
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "欢迎使用 Gin CloudBase 应用！",
			"status":  "running",
		})
	})

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":      "healthy",
			"timestamp":   time.Now().Format(time.RFC3339),
			"framework":   "Gin",
			"go_version":  "1.19+",
			"gin_version": gin.Version,
		})
	})

	// API 路由组
	api := router.Group("/api")
	{
		users := api.Group("/users")
		{
			users.GET("", controllers.GetUsers)
			users.GET("/:id", controllers.GetUser)
			users.POST("", controllers.CreateUser)
		}
	}

	// 获取端口，支持环境变量
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 启动服务器
	router.Run(":" + port)
}
```

## 第三步：本地测试

### 更新依赖

```bash
go mod tidy
```

### 启动应用

```bash
go run main.go
```

### 测试 API 接口

```bash
# 测试健康检查
curl http://localhost:8080/health

# 测试首页
curl http://localhost:8080/

# 测试用户列表
curl http://localhost:8080/api/users

# 测试分页
curl "http://localhost:8080/api/users?page=1&limit=2"

# 测试获取单个用户
curl http://localhost:8080/api/users/1

# 测试创建用户
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"新用户","email":"newuser@example.com"}'
```

## 第四步：准备部署文件

根据您选择的部署方式，需要准备不同的配置文件：

### 📋 选择部署方式

<details>
<summary><strong>🔥 HTTP 云函数部署配置</strong></summary>

HTTP 云函数需要 `scf_bootstrap` 启动脚本和编译后的二进制文件。

#### 1. 创建启动脚本

创建 `scf_bootstrap` 文件（无扩展名）：

```bash
#!/bin/bash
export PORT=9000
export GIN_MODE=release
./main
```

为启动脚本添加执行权限：

```bash
chmod +x scf_bootstrap
```

#### 2. 编译应用

编译 Go 应用为 Linux 二进制文件：

```bash
# 交叉编译为 Linux 64位
GOOS=linux GOARCH=amd64 go build -o main .
```

#### 3. 项目结构

```
cloudrun-gin/
├── controllers/
│   └── user.go
├── models/
│   └── user.go
├── main.go
├── go.mod
├── go.sum
├── main                   # 编译后的二进制文件
└── scf_bootstrap         # 🔑 云函数启动脚本
```

> 💡 **说明**：
> - `scf_bootstrap` 是 CloudBase 云函数的启动脚本
> - 需要将 Go 应用编译为 Linux 二进制文件
> - 设置 `PORT=9000` 环境变量确保应用监听正确端口
> - 设置 `GIN_MODE=release` 优化性能

</details>

<details>
<summary><strong>🐳 云托管部署配置</strong></summary>

云托管使用 Docker 容器化部署，需要 `Dockerfile` 配置文件。

#### 1. 创建 Dockerfile

创建 `Dockerfile` 文件：

```dockerfile
# 构建阶段
FROM golang:1.22.3-alpine as builder

# 安装依赖包，选用国内镜像源以提高下载速度
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tencent.com/g' /etc/apk/repositories \
    && apk add --update --no-cache git

# 指定构建过程中的工作目录
WORKDIR /app

# 将 go.mod 和 go.sum 文件拷贝到工作目录
COPY go.mod go.sum ./

# 设置 Go 模块代理
ENV GOPROXY=https://goproxy.cn,direct

# 下载依赖
RUN go mod download

# 将当前目录下所有文件都拷贝到工作目录下
COPY . .

# 执行代码编译命令。操作系统参数为linux，编译后的二进制产物命名为main
RUN GOOS=linux go build -o main .

# 运行阶段
FROM alpine:3.18

# 安装依赖包，选用国内镜像源以提高下载速度
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tencent.com/g' /etc/apk/repositories \
    && apk add --update --no-cache ca-certificates tzdata \
    && rm -f /var/cache/apk/*

# 容器默认时区为UTC，如需使用上海时间请启用以下时区设置命令
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo Asia/Shanghai > /etc/timezone

# 指定运行时的工作目录
WORKDIR /app

# 将构建产物 main 拷贝到运行时的工作目录中
COPY --from=builder /app/main /app/

# 暴露端口
EXPOSE 8080

# 设置环境变量
ENV PORT=8080
ENV GIN_MODE=release

# 执行启动命令
CMD ["./main"]
```

#### 2. 创建 .dockerignore 文件

创建 `.dockerignore` 文件以优化构建性能：

```
.git
.gitignore
README.md
.env
.DS_Store
*.log
scf_bootstrap
main
```

#### 3. 项目结构

```
cloudrun-gin/
├── controllers/
│   └── user.go
├── models/
│   └── user.go
├── main.go
├── go.mod
├── go.sum
├── Dockerfile            # 🔑 容器配置文件
└── .dockerignore         # Docker 忽略文件
```

> 💡 **说明**：
> - 云托管支持自定义端口，默认使用 8080 端口
> - 使用多阶段构建优化镜像大小
> - 支持完整的 Go 模块和依赖管理

</details>

## 第五步：项目结构

确保您的项目目录结构包含必要的文件。根据部署方式的不同，某些文件是可选的：

```
cloudrun-gin/
├── controllers/
│   └── user.go           # 用户控制器
├── models/
│   └── user.go           # 用户模型
├── main.go               # 主应用文件
├── go.mod                # Go 模块文件
├── go.sum                # 依赖锁定文件
├── main                  # 编译后的二进制文件 (仅云函数需要)
├── scf_bootstrap         # HTTP 云函数启动脚本 (仅云函数需要)
├── Dockerfile            # 云托管容器配置 (仅云托管需要)
└── .dockerignore         # Docker 忽略文件 (仅云托管需要)
```

## 第六步：部署应用

选择您需要的部署方式：

### 🚀 部署方式选择

<details>
<summary><strong>🔥 部署到 HTTP 云函数</strong></summary>

#### 通过控制台部署

1. 登录 [CloudBase 控制台](https://console.cloud.tencent.com/tcb)
2. 选择您的环境，进入「云函数」页面
3. 点击「新建云函数」
4. 选择「HTTP 云函数」
5. 填写函数名称（如：`gin-app`）
6. 选择运行时：**Go 1.x**（或其他支持的版本）
7. 提交方法选择：**本地上传代码包**
8. 上传编译后的二进制文件和 `scf_bootstrap` 文件
9. **自动安装依赖**：关闭（Go 应用无需此选项）
10. 点击「创建」按钮等待部署完成

#### 打包部署

如果需要手动打包：

```bash
# 编译应用
GOOS=linux GOARCH=amd64 go build -o main .

# 创建部署包
zip -j gin-app.zip main scf_bootstrap
```

</details>

<details>
<summary><strong>🐳 部署到云托管</strong></summary>

#### 通过控制台部署

1. 登录 [CloudBase 控制台](https://console.cloud.tencent.com/tcb)
2. 选择您的环境，进入「云托管」页面
3. 点击「新建服务」
4. 填写服务名称（如：`gin-service`）
5. 选择「本地代码」上传方式
6. 上传包含 `Dockerfile` 的项目目录
7. 配置服务参数：
   - **端口**：8080（或您在应用中配置的端口）
   - **CPU**：0.25 核
   - **内存**：0.5 GB
   - **实例数量**：1-10（根据需求调整）
8. 点击「创建」按钮等待部署完成

#### 模板部署（快速开始）

1. 登录 [腾讯云托管控制台](https://tcb.cloud.tencent.com/dev#/platform-run/service/create?type=image)
2. 点击「通过模板部署」，选择 **Gin 模板**
3. 输入自定义服务名称，点击部署
4. 等待部署完成后，点击左上角箭头，返回到服务详情页
5. 点击概述，获取默认域名并访问

</details>

## 第七步：访问您的应用

### HTTP 云函数访问

部署成功后，您可以参考[通过 HTTP 访问云函数](https://docs.cloudbase.net/service/access-cloud-function)设置自定义域名访问 HTTP 云函数。

访问地址格式：`https://your-function-url/`

### 云托管访问

云托管部署成功后，系统会自动分配访问地址。您也可以绑定自定义域名。

访问地址格式：`https://your-service-url/`

### 测试接口

无论使用哪种部署方式，您都可以测试以下接口：

- **根路径**：`/` - Gin 欢迎页面
- **健康检查**：`/health` - 查看应用状态
- **用户列表**：`/api/users` - 获取用户列表
- **用户详情**：`/api/users/1` - 获取特定用户
- **创建用户**：`POST /api/users` - 创建新用户

### 示例请求

```bash
# 健康检查
curl https://your-app-url/health

# 获取用户列表
curl https://your-app-url/api/users

# 分页查询
curl "https://your-app-url/api/users?page=1&limit=2"

# 创建新用户
curl -X POST https://your-app-url/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"测试用户","email":"test@example.com"}'
```

## 常见问题

### ❓ 问题分类

<details>
<summary><strong>🔥 HTTP 云函数相关问题</strong></summary>

#### Q: 为什么 HTTP 云函数必须使用 9000 端口？
A: CloudBase HTTP 云函数要求应用监听 9000 端口，这是平台的标准配置。应用通过环境变量 `PORT=9000` 来控制端口，本地开发时默认使用 8080 端口。

#### Q: Go 应用如何进行交叉编译？
A: 使用以下命令进行交叉编译：
```bash
GOOS=linux GOARCH=amd64 go build -o main .
```

#### Q: 如何优化 Go 应用的冷启动时间？
A: 
- 减少依赖包数量
- 使用 `gin.ReleaseMode` 模式
- 避免在启动时进行重复的初始化操作
- 合理设置内存配置

#### Q: scf_bootstrap 文件有什么作用？
A: `scf_bootstrap` 是云函数的启动脚本，用于设置环境变量和启动编译后的二进制文件。

</details>

<details>
<summary><strong>🐳 云托管相关问题</strong></summary>

#### Q: 云托管支持哪些端口？
A: 云托管支持自定义端口，Gin 应用默认使用 8080 端口，也可以根据需要配置其他端口。

#### Q: 如何优化 Docker 镜像构建速度？
A: 
- 使用多阶段构建
- 配置 Go 模块代理 `GOPROXY=https://goproxy.cn`
- 使用国内镜像源
- 合理设置 `.dockerignore`

#### Q: Dockerfile 中为什么使用 Alpine Linux？
A: Alpine Linux 是轻量级的 Linux 发行版，镜像体积小，安全性高，适合容器化部署。

#### Q: 如何处理 Go 模块依赖？
A: 
- 在构建阶段先复制 `go.mod` 和 `go.sum`
- 使用 `go mod download` 下载依赖
- 设置合适的 Go 模块代理

</details>

<details>
<summary><strong>🔧 通用问题</strong></summary>

#### Q: 如何处理静态文件？
A: Gin 可以使用 `router.Static()` 方法来处理静态文件服务。

#### Q: 如何查看应用日志？
A: 
- **HTTP 云函数**：在 CloudBase 控制台的云函数页面查看运行日志
- **云托管**：在云托管服务详情页面查看实例日志

#### Q: 支持哪些 Go 版本？
A: CloudBase 支持 Go 1.18、1.19、1.20 等版本，建议使用最新的稳定版本。

#### Q: 如何处理 CORS 跨域问题？
A: 可以使用中间件或手动设置响应头来处理跨域请求，示例代码中已包含 CORS 处理。

#### Q: 两种部署方式如何选择？
A: 
- **选择 HTTP 云函数**：轻量级 API 服务、间歇性访问、成本敏感
- **选择云托管**：企业级应用、持续运行、需要更多控制权

</details>

## 最佳实践

### 1. 环境变量管理

在应用中使用环境变量：

```go
import "os"

func main() {
    // 端口配置
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    
    // 数据库配置
    dbURL := os.Getenv("DATABASE_URL")
    
    // 应用模式
    if os.Getenv("GIN_MODE") == "" {
        gin.SetMode(gin.ReleaseMode)
    }
}
```

### 2. 端口配置策略

为了同时支持两种部署方式，建议使用动态端口配置：

```go
// 支持多环境端口配置（默认 8080，云函数环境通过 PORT=9000 控制）
port := os.Getenv("PORT")
if port == "" {
    port = "8080"
}
router.Run(":" + port)
```

### 3. 添加中间件

使用 Gin 中间件增强功能：

```go
// 恢复中间件
router.Use(gin.Recovery())

// 自定义日志中间件
router.Use(gin.LoggerWithConfig(gin.LoggerConfig{
    Formatter: func(param gin.LogFormatterParams) string {
        return fmt.Sprintf("%s - %s \"%s %s\" %d %s\n",
            param.TimeStamp.Format("2006/01/02 - 15:04:05"),
            param.ClientIP,
            param.Method,
            param.Path,
            param.StatusCode,
            param.Latency,
        )
    },
}))

// 限流中间件
router.Use(ratelimit.RateLimiter(store, &ratelimit.Options{
    ErrorHandler: func(c *gin.Context, info ratelimit.Info) {
        c.String(429, "Too many requests. Try again in "+time.Until(info.ResetTime).String())
    },
}))
```

### 4. 错误处理

实现统一的错误处理：

```go
// 全局错误处理中间件
func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        
        if len(c.Errors) > 0 {
            err := c.Errors.Last()
            
            switch err.Type {
            case gin.ErrorTypeBind:
                c.JSON(http.StatusBadRequest, gin.H{
                    "success": false,
                    "message": "请求参数错误",
                    "error":   err.Error(),
                })
            default:
                c.JSON(http.StatusInternalServerError, gin.H{
                    "success": false,
                    "message": "服务器内部错误",
                })
            }
        }
    }
}
```

### 5. 结构化日志

使用结构化日志记录：

```go
import "github.com/sirupsen/logrus"

func init() {
    logrus.SetFormatter(&logrus.JSONFormatter{})
    logrus.SetLevel(logrus.InfoLevel)
}

func LoggerMiddleware() gin.HandlerFunc {
    return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
        logrus.WithFields(logrus.Fields{
            "status_code": param.StatusCode,
            "latency":     param.Latency,
            "client_ip":   param.ClientIP,
            "method":      param.Method,
            "path":        param.Path,
        }).Info("HTTP Request")
        
        return ""
    })
}
```

### 6. 健康检查增强

增强健康检查接口：

```go
func HealthCheck(c *gin.Context) {
    // 检查数据库连接
    dbStatus := "healthy"
    if err := db.Ping(); err != nil {
        dbStatus = "unhealthy"
    }
    
    // 检查外部服务
    externalStatus := "healthy"
    // ... 检查逻辑
    
    status := "healthy"
    if dbStatus != "healthy" || externalStatus != "healthy" {
        status = "unhealthy"
        c.Status(http.StatusServiceUnavailable)
    }
    
    c.JSON(c.Writer.Status(), gin.H{
        "status":          status,
        "timestamp":       time.Now().Format(time.RFC3339),
        "framework":       "Gin",
        "go_version":      runtime.Version(),
        "gin_version":     gin.Version,
        "database":        dbStatus,
        "external_service": externalStatus,
    })
}
```

## 进阶功能

### 数据库集成

集成 GORM 和 MySQL：

```bash
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql
```

```go
import (
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

func initDB() *gorm.DB {
    dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
        dsn = "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
    }
    
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }
    
    return db
}
```

### 身份验证

添加 JWT 身份验证：

```bash
go get -u github.com/golang-jwt/jwt/v4
```

### API 文档

使用 Swagger 生成 API 文档：

```bash
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files
```

### 缓存支持

添加 Redis 缓存：

```bash
go get -u github.com/go-redis/redis/v8
```
