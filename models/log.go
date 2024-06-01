package models

import (
	"time"

	"gorm.io/gorm"
)

// Log is the model for the log
type Log struct {
	gorm.Model
	AppID     uint      `json:"app_id"`
	Namespace string    `json:"namespace"`
	LogType   string    `json:"log_type"` // "error" or "success"
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}
