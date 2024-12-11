package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/rozhnof/auth-service/internal/presentation/handlers"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitAuthRoutes(router gin.IRouter, authHandler *handlers.AuthHandler, googleAuthHandler *handlers.GoogleAuthHandler) {
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/refresh", authHandler.Refresh)
		authGroup.GET("/confirm", authHandler.Confirm)

		googleAuthGroup := router.Group("/google")
		{
			googleAuthGroup.GET("/login", googleAuthHandler.Login)
			googleAuthGroup.GET("/callback", googleAuthHandler.Callback)
		}
	}
}

func InitSwaggerRoutes(router gin.IRoutes) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
