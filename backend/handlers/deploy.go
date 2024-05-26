package handlers

import (
	"net/http"
	"time"

	conf "github.com/aeswibon/helmdeploy/backend/config"
	"github.com/aeswibon/helmdeploy/backend/k8"
	"github.com/aeswibon/helmdeploy/backend/models"
	"github.com/gin-gonic/gin"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// Deploy deploys an application to the Kubernetes cluster
func Deploy(c *gin.Context) {
	var req k8.HelmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate request
	if req.Namespace == "" || req.ReleaseName == "" || req.Chart == "" || req.Name == "" || req.URL == "" || req.Version == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Namespace, release name, chart, name, version and URL are required"})
		return
	}

	// Load Kubernetes configuration
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	kubeClient := k8.NewKubernetesClient(clientset)

	// Create namespace if it doesn't exist
	if err := kubeClient.CreateNs(req.Namespace); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Install Helm chart
	output, err := kubeClient.InstallChart(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "details": output})
		return
	}

	// Check Helm release status
	status, err := kubeClient.CheckRelease(req.Namespace, req.ReleaseName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Store application info in the database
	app := models.Application{
		Namespace:    req.Namespace,
		AppName:      req.ReleaseName,
		DeployedAt:   time.Now(),
		HealthStatus: status,
		Logs:         output,
	}

	if err := conf.DB.Create(&app).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Application deployed successfully", "status": status})
}
