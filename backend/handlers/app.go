package handlers

import (
	"context"
	"net/http"

	conf "github.com/aeswibon/helmdeploy/backend/config"
	"github.com/aeswibon/helmdeploy/backend/k8"
	"github.com/aeswibon/helmdeploy/backend/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// GetApps returns all applications
func GetApps(c *gin.Context) {
	var apps []models.Application
	session := sessions.Default(c)
	username := session.Get("username").(string)
	if err := conf.DB.Where("created_by = ?", username).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, apps)
}

// GetAppLogs returns logs for a specific application
func GetAppLogs(c *gin.Context) {
	id := c.Param("id")
	var app, foundApp models.Application
	if err := conf.DB.First(&app, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	session := sessions.Default(c)
	username := session.Get("username").(string)
	if err := conf.DB.Where("id = ? AND created_by = ?", id, username).First(&foundApp).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
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

	la := k8.NewLogger(clientset, app.Namespace, "Success", "Error")

	// Retrieve logs for the specific application
	logs, ok := la.LogIndex[app.AppName]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Logs not found for application"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"logs": logs})
}

// DeleteApp deletes an application
func DeleteApp(c *gin.Context) {
	id := c.Param("id")
	session := sessions.Default(c)
	username := session.Get("username").(string)

	var app models.Application
	if err := conf.DB.Where("id = ? AND created_by = ?", id, username).First(&app).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
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

	// Delete Helm release
	if err := kubeClient.DeleteRelease(app.Namespace, app.AppName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Delete namespace
	if err := clientset.CoreV1().Namespaces().Delete(context.Background(), app.Namespace, metav1.DeleteOptions{}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Delete application record from the database
	if err := conf.DB.Delete(&app).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Application deleted successfully"})
}
