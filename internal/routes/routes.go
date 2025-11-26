package routes

import (
	"cloudrun-gin/internal/controllers"

	"github.com/gin-gonic/gin"
)

type RouterOptions func(router *gin.Engine)

var opts = []RouterOptions{
	controllers.InitUserRouters,
	controllers.InitWelcomeRouters,
}

// InitRoutes 初始化路由规则
func InitRoutes(router *gin.Engine) {
	for _, opt := range opts {
		opt(router)
	}
}
