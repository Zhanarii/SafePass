package handlers

import (
	"net/http"

	"safepass-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetEmployees(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var employees []models.Employee

		if err := db.Find(&employees).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"employees": employees,
			"total":     len(employees),
		})
	}
}

func GetEmployee(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var employee models.Employee
		if err := db.First(&employee, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
			return
		}

		c.JSON(http.StatusOK, employee)
	}
}
