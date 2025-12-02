package handlers

import (
	"net/http"
	"strconv"

	"safepass-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetViolations(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		limitStr := c.DefaultQuery("limit", "50")
		limit, _ := strconv.Atoi(limitStr)

		var violations []models.Violation
		result := db.Preload("Employee").Order("timestamp DESC").Limit(limit).Find(&violations)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"violations": violations,
			"total":      len(violations),
		})
	}
}
