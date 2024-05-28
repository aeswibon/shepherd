package main

import (
	"crypto/rand"
	"fmt"
	"log"

	"github.com/aeswibon/helmdeploy/backend/config"
	"github.com/aeswibon/helmdeploy/backend/handlers"
	"github.com/aeswibon/helmdeploy/backend/middleware"
	"github.com/aeswibon/helmdeploy/backend/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func keyGen() []byte {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		fmt.Println("Error generating session key:", err)
		return nil
	}
	return key
}

func main() {
	// Connect to the database
	config.ConnectDB()
	// Initialize cron job
	utils.InitCron()

	r := gin.Default()
	r.Use(sessions.Sessions("helmdeployer", cookie.NewStore(keyGen())))

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

	if err := r.Run(fmt.Sprintf(":%s", config.GetEnv("PORT"))); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}
