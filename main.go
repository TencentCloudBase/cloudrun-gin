package main

import (
	"cloudrun-gin/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// 加载全局视图文件
	router.LoadHTMLGlob("static/html/*")

	// 初始化路由
	routes.InitRoutes(router)

	// 运行服务器
	router.Run(":8080") // 监听并在 0.0.0.0:8080 上启动服务
}
