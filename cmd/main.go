package main

import (
	"crypto/rand"
	"fmt"
	"log"

	conf "github.com/aeswibon/shepherd/config"
	"github.com/aeswibon/shepherd/handlers"
	"github.com/aeswibon/shepherd/k8"
	"github.com/aeswibon/shepherd/middleware"
	"github.com/aeswibon/shepherd/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
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
	conf.ConnectDB()
	// Initialize cron job
	utils.InitCron()

	// Load Kubernetes configuration
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		log.Fatalf("Failed to load Kubernetes configuration: %v", err)
		return
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create Kubernetes client: %v", err)
		return
	}

	k8Client := k8.NewKubernetesClient(clientset)
	helmClient, err := k8.NewHelmClient()
	if err != nil {
		log.Fatalf("Failed to create Helm client: %v", err)
	}

	appHandler := handlers.NewAppHandler(k8Client, helmClient)

	// Initialize the Gin engine
	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(cors.Default())
	r.Use(sessions.Sessions("helmdeployer", cookie.NewStore(keyGen())))

	api := r.Group("/api")
	{
		// Public routes
		api.POST("/auth/signup", handlers.Signup)
		api.POST("/auth/login", handlers.Login)

		// Protected routes
		protected := api.Group("/app", middleware.AuthMiddleware())
		{
			protected.GET("/", appHandler.GetApps)
			protected.GET("/:id/logs", appHandler.GetAppLogs)
			protected.POST("/deploy", appHandler.Deploy)
			protected.DELETE("/:id", appHandler.DeleteApp)
		}
	}

	if err := r.Run(fmt.Sprintf(":%s", conf.GetEnv("PORT"))); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}
