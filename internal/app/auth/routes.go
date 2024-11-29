package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/rozhnof/auth-service/internal/presentation/handlers"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRoutes(router gin.IRoutes, authHandler *handlers.AuthHandler, googleAuthHandler *handlers.GoogleAuthHandler) {
	router.POST("/auth/register", authHandler.Register)
	router.POST("/auth/login", authHandler.Login)
	router.POST("/auth/refresh", authHandler.Refresh)

	router.GET("/auth/google/login", googleAuthHandler.Login)
	router.GET("/auth/google/callback", googleAuthHandler.Callback)
}

func InitSwagger(router gin.IRoutes) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
