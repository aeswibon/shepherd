package utils

import (
	"time"

	"github.com/aeswibon/shepherd/config"
	"github.com/aeswibon/shepherd/models"
	"github.com/robfig/cron/v3"
)

func evictOldLogs() {
	threshold := time.Now().AddDate(0, 0, -7)
	config.DB.Where("created_at < ?", threshold).Delete(&models.Log{})
}

// InitCron initializes cron jobs
func InitCron() {
	c := cron.New()
	c.AddFunc("@daily", func() { evictOldLogs() })
	c.Start()
}
