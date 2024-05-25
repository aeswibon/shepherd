package models

import (
	"time"

	"gorm.io/gorm"
)

// Application is the model for the app
type Application struct {
	gorm.Model
	Namespace    string    `json:"namespace"`
	AppName      string    `json:"app_name"`
	DeployedAt   time.Time `json:"deployed_at"`
	HealthStatus string    `json:"health_status"`
	Logs         string    `json:"logs"`
}
