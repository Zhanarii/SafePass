package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"safepass-backend/models"
	"safepass-backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func DetectPPE(detector *services.Detector, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		file, header, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
			return
		}
		defer file.Close()

		tempDir := "./temp"
		os.MkdirAll(tempDir, 0755)

		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), header.Filename)
		filepath := filepath.Join(tempDir, filename)

		out, err := os.Create(filepath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}
		defer out.Close()

		if _, err := io.Copy(out, file); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}

		result, err := detector.DetectFromFile(filepath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Detection failed: %v", err)})
			return
		}

		if result.HasViolations {
			for _, v := range result.Violations {
				violation := models.Violation{
					ViolationType: v.Class,
					Confidence:    v.Confidence,
					ImageURL:      filepath,
					Status:        "new",
					CameraID:      1,
					Timestamp:     time.Now(),
				}
				db.Create(&violation)
			}
		}

		c.JSON(http.StatusOK, result)
	}
}
