package service

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// User ...
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// ListUsersResp ...
type ListUsersResp struct {
	Total int    `json:"total"`
	Items []User `json:"items"`
}

// ListUsers ...
func ListUsers(c *gin.Context) {
	for k, v := range c.Request.Header {
		fmt.Printf("k=%v, v=%v\n", k, v)
	}
	users := make([]User, 0)

	// 添加第一个用户
	user_0 := User{
		ID:   0,
		Name: "zhangsan",
	}
	users = append(users, user_0)

	// 添加第二个用户
	user_1 := User{
		ID:   1,
		Name: "lisi",
	}
	users = append(users, user_1)

	// 组合返回列表
	result := ListUsersResp{
		Total: 2,
		Items: users,
	}
	c.JSON(200, result)
}
