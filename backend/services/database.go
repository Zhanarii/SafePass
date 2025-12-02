package services

import (
	"fmt"
	"log"

	"safepass-backend/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB() (*gorm.DB, error) {
	// Используем SQLite для простоты (можно заменить на PostgreSQL)
	db, err := gorm.Open(sqlite.Open("safepass.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// Auto-migrate models
	if err := db.AutoMigrate(
		&models.Employee{},
		&models.Violation{},
	); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	seedData(db)

	log.Println("Database connected successfully")
	return db, nil
}

func seedData(db *gorm.DB) {
	var count int64
	db.Model(&models.Employee{}).Count(&count)
	if count > 0 {
		return
	}

	employees := []models.Employee{
		{Name: "Арайлым Жазыбаева", Position: "Инженер безопасности", BadgeID: "B001", Status: "active", PhotoURL: "https://randomuser.me/api/portraits/women/65.jpg"},
		{Name: "Данияр Сулейменов", Position: "Прораб", BadgeID: "B002", Status: "inactive", PhotoURL: "https://randomuser.me/api/portraits/men/12.jpg"},
		{Name: "Айшан Тлеуберген", Position: "Электрик", BadgeID: "B003", Status: "active", PhotoURL: "https://randomuser.me/api/portraits/women/22.jpg"},
	}

	db.Create(&employees)
	log.Println("Database seeded with test data")
}
