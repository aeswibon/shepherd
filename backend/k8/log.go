package k8

import (
	"bytes"
	"log"
	"os/exec"
	"sync"
	"time"

	"github.com/aeswibon/helmdeploy/backend/config"
	"github.com/aeswibon/helmdeploy/backend/models"
)

// FetchLogs fetches logs concurrently and store them in the database
func FetchLogs(app *models.Application) error {
	var wg sync.WaitGroup
	logTypes := []string{"success", "error"}

	for _, logType := range logTypes {
		wg.Add(1)
		go func(logType string) {
			defer wg.Done()

			cmd := exec.Command("kubectl", "logs", "-n", app.Namespace, "deployment/"+app.AppName)
			var out bytes.Buffer
			cmd.Stdout = &out
			if err := cmd.Run(); err != nil {
				log.Printf("Error: %v", err)
				config.DB.Create(&models.Log{
					AppID:     app.ID,
					LogType:   "error",
					Message:   err.Error(),
					CreatedAt: time.Now(),
				})
				return
			}
			config.DB.Create(&models.Log{
				AppID:     app.ID,
				LogType:   logType,
				Message:   out.String(),
				CreatedAt: time.Now(),
			})
		}(logType)
	}

	wg.Wait()
	return nil
}
