# Gin 项目创建指南

本指南详细介绍如何从零开始创建一个适用于 CloudBase 部署的 Gin 项目。

## 📋 目录导航

- [环境准备](#环境准备)
- [创建项目](#创建项目)
- [基础配置](#基础配置)
- [创建模型](#创建模型)
- [创建控制器](#创建控制器)
- [配置路由](#配置路由)
- [本地测试](#本地测试)
- [下一步](#下一步)

---

## 环境准备

### 1. 检查 Go 版本

```bash
# 检查 Go 版本（推荐 1.19+）
go version
```

### 2. 创建项目目录

```bash
# 创建项目根目录
mkdir cloudrun-gin && cd cloudrun-gin

# 初始化 Go 模块
go mod init cloudrun-gin
```

## 创建项目

### 1. 安装 Gin 框架

```bash
# 安装 Gin 框架
go get github.com/gin-gonic/gin@v1.9.1

# 验证安装
go mod tidy
```

### 2. 创建项目结构

```bash
# 创建目录结构
mkdir -p controllers models

# 查看项目结构
tree .
# cloudrun-gin/
# ├── controllers/
# ├── models/
# ├── go.mod
# └── go.sum
```

## 基础配置

### 1. 创建数据模型

创建 `models/user.go`：

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

## 创建控制器

### 1. 用户控制器

创建 `controllers/user.go`：

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

## 配置路由

### 1. 创建主应用文件

创建 `main.go`：

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

## 本地测试

### 1. 编译和运行

```bash
# 下载依赖
go mod tidy

# 编译项目
go build -o main .

# 运行项目
./main
# 或者直接运行
go run main.go
```

### 2. API 测试

```bash
# 测试基础接口
curl http://localhost:8080/
# 返回: {"message":"欢迎使用 Gin CloudBase 应用！","status":"running"}

curl http://localhost:8080/health
# 返回: {"status":"healthy","timestamp":"...","framework":"Gin",...}

# 测试用户 API
# 获取用户列表
curl http://localhost:8080/api/users
curl "http://localhost:8080/api/users?page=1&limit=2"

# 获取单个用户
curl http://localhost:8080/api/users/1

# 创建用户
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"新用户","email":"newuser@example.com"}'
```

### 3. 生成 go.mod 和 go.sum

确保项目依赖正确：

```bash
# 查看 go.mod 内容
cat go.mod
# module cloudrun-gin
# 
# go 1.23.10
# 
# require github.com/gin-gonic/gin v1.9.1

# 查看依赖
go list -m all
```

## 下一步

项目创建完成后，根据您的部署需求选择相应的部署指南：

### 🚀 部署选择

| 部署方式 | 适用场景 | 详细指南 |
|----------|----------|----------|
| **HTTP 云函数** | 轻量级 API、间歇性访问 | [HTTP 云函数部署指南](./http-function.md) |
| **云托管** | 企业应用、高并发、持续运行 | [云托管部署指南](./cloud-run.md) |

### 📚 相关文档

- [返回主文档](../README.md)
- [HTTP 云函数部署指南](./http-function.md)
- [云托管部署指南](./cloud-run.md)

### 🔧 进一步开发

1. **添加更多路由**：扩展 API 功能
2. **数据库集成**：使用 GORM 连接数据库
3. **中间件开发**：认证、限流等中间件
4. **配置管理**：使用 Viper 管理配置
5. **日志系统**：集成结构化日志

---

**提示**：确保在部署前测试所有功能，特别是 API 接口和错误处理。