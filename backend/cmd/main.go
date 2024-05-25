package main

import (
	"github.com/aeswibon/helmdeploy/backend/config"
	"github.com/aeswibon/helmdeploy/backend/handlers"
	"github.com/aeswibon/helmdeploy/backend/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to the database
	config.ConnectDB()

	r := gin.Default()

	// Public routes
	r.POST("/auth/signup", handlers.Signup)
	r.POST("/auth/login", handlers.Login)

	// Protected routes
	protected := r.Group("/app", middleware.AuthMiddleware())
	{
		protected.GET("/", handlers.GetApps)
		protected.GET("/:id/logs", handlers.GetAppLogs)
		protected.POST("/deploy", handlers.Deploy)
		protected.DELETE("/:id", handlers.DeleteApp)
	}

	r.Run(config.GetEnv("PORT")) // listen and serve on
}
