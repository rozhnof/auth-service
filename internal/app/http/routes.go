package http_app

import (
	http_handlers "auth/internal/auth/presentation/handlers/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitAuthRoutes(router gin.IRoutes, userHandler *http_handlers.AuthHandler) {
	router.POST("/auth/register", userHandler.Register)
	router.POST("/auth/login", userHandler.Login)
	router.POST("/auth/refresh", userHandler.Refresh)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
