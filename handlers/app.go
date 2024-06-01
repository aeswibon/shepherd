package handlers

import (
	"log"
	"net/http"
	"time"

	conf "github.com/aeswibon/shepherd/config"
	"github.com/aeswibon/shepherd/k8"
	"github.com/aeswibon/shepherd/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// AppHandler handles application related operations
type AppHandler struct {
	k8sClient  k8.KubernetesClient
	helmClient k8.HelmClient
}

// NewAppHandler creates a new AppHandler
func NewAppHandler(k8sClient k8.KubernetesClient, helmClient k8.HelmClient) *AppHandler {
	return &AppHandler{
		k8sClient:  k8sClient,
		helmClient: helmClient,
	}
}

// GetApps returns all applications
func (ah *AppHandler) GetApps(c *gin.Context) {
	var apps []models.Application
	session := sessions.Default(c)
	username := session.Get("username").(string)
	log.Println("Username: ", username)
	if err := conf.DB.Where("created_by = ?", username).Find(&apps).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, apps)
}

// GetAppLogs returns logs for a specific application
func (ah *AppHandler) GetAppLogs(c *gin.Context) {
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

	logs, err := ah.k8sClient.GetLogs(foundApp.Namespace, foundApp.AppName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Store logs in the database
	log := models.Log{
		AppID:     app.ID,
		Namespace: app.Namespace,
		LogType:   "success",
		Message:   string(logs),
		CreatedAt: time.Now(),
	}

	if err := conf.DB.Create(&log).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, string(logs))
}

// DeleteApp deletes an application
func (ah *AppHandler) DeleteApp(c *gin.Context) {
	id := c.Param("id")
	session := sessions.Default(c)
	username := session.Get("username").(string)

	var app models.Application
	if err := conf.DB.Where("id = ? AND created_by = ?", id, username).First(&app).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	// Delete Helm release
	if err := ah.helmClient.UninstallRelease(app.Namespace, app.AppName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Delete namespace
	if err := ah.k8sClient.DeleteNs(app.Namespace); err != nil {
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

// Deploy deploys an application to the Kubernetes cluster
func (ah *AppHandler) Deploy(c *gin.Context) {
	var req models.DeployRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate request
	if req.Namespace == "" || req.ChartName == "" || req.Values == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Note: namespace, chart name and values are required fields"})
		return
	}

	// Create namespace if it doesn't exist
	if err := ah.k8sClient.CreateNs(req.Namespace); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Install Helm chart
	output, err := ah.helmClient.InstallChart(req.Namespace, req.ChartName, req.AppName, req.Values)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "details": output})
		return
	}

	// Check Helm release status
	status, err := ah.helmClient.GetReleaseStatus(req.Namespace, req.AppName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	session := sessions.Default(c)
	username := session.Get("username")
	// Store application info in the database
	app := models.Application{
		Namespace:    req.Namespace,
		AppName:      req.AppName,
		DeployedAt:   time.Now(),
		HealthStatus: status,
		CreatedBy:    username.(string),
	}

	if err := conf.DB.Create(&app).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Application deployed successfully", "status": status})
}
