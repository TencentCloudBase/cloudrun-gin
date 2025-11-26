package controllers

import (
	"cloudrun-gin/internal/service"

	"github.com/gin-gonic/gin"
)

// InitUserRouters
func InitUserRouters(e *gin.Engine) {
	userGroup := e.Group("/api/users")
	{
		userGroup.GET("", service.ListUsers)
	}
}
