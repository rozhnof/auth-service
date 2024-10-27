package http_app

import (
	"auth/internal/user/handlers"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRoutes(router gin.IRoutes, userHandler *handlers.AuthHandler) {
	router.POST("/auth/register", userHandler.Register)
	router.POST("/auth/login", userHandler.Login)
	router.POST("/auth/refresh", userHandler.Refresh)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
