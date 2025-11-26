package controllers

import (
	"github.com/gin-gonic/gin"
)

// InitWelcomeRouters 初始化路由信息
func InitWelcomeRouters(e *gin.Engine) {
	e.GET("/", RenderWelcomePage)
}

// RenderWelcomePage 渲染首页
func RenderWelcomePage(c *gin.Context) {
	c.HTML(200, "index.html", gin.H{
		"title": "Gin Demo - Powered By CloudBase",
	})
}
