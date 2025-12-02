package main

import (
	"log"
	"time"

	"safepass-backend/handlers"
	"safepass-backend/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Инициализация базы данных
	db, err := services.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Инициализация детектора
	detector, err := services.NewDetector("../ml_model/runs/detect/train4/weights/best.pt")
	if err != nil {
		log.Fatal("Failed to load YOLO model:", err)
	}

	// Gin router
	r := gin.Default()

	// CORS для React
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// API routes
	api := r.Group("/api")
	{
		api.GET("/statistics", handlers.GetStatistics(db))
		api.GET("/violations", handlers.GetViolations(db))
		api.GET("/employees", handlers.GetEmployees(db))
		api.GET("/employees/:id", handlers.GetEmployee(db))
		api.POST("/detect", handlers.DetectPPE(detector, db))
	}

	// WebSocket для камер
	r.GET("/ws/camera/:id", handlers.HandleCameraWebSocket(detector, db))

	// Health check
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "SafePass API is running",
			"status":  "ok",
		})
	})

	// Запуск сервера
	log.Println("Server starting on :8000")
	if err := r.Run(":8000"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
