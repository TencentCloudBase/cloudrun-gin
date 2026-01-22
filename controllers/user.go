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
