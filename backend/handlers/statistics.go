package handlers

import (
	"net/http"
	"time"

	"safepass-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetStatistics(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var stats models.Statistics

		// Total employees
		var totalEmployees int64
		db.Model(&models.Employee{}).Count(&totalEmployees)
		stats.TotalEmployees = int(totalEmployees)

		// Active employees
		var activeEmployees int64
		db.Model(&models.Employee{}).Where("status = ?", "active").Count(&activeEmployees)
		stats.ActiveEmployees = int(activeEmployees)

		// Total violations
		var totalViolations int64
		db.Model(&models.Violation{}).Count(&totalViolations)
		stats.TotalViolations = int(totalViolations)

		// Violations today
		today := time.Now().Truncate(24 * time.Hour)
		var violationsToday int64
		db.Model(&models.Violation{}).
			Where("timestamp >= ?", today).
			Count(&violationsToday)
		stats.ViolationsToday = int(violationsToday)

		// Compliance rate
		if stats.TotalEmployees > 0 {
			violationRate := float64(stats.TotalViolations) / float64(stats.TotalEmployees)
			stats.ComplianceRate = (1 - violationRate) * 100
			if stats.ComplianceRate < 0 {
				stats.ComplianceRate = 0
			}
		} else {
			stats.ComplianceRate = 100
		}

		c.JSON(http.StatusOK, stats)
	}
}
