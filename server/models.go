package main

import (
	"errors" // Import added for error handling
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"github.com/my/project/models"
)

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	err = db.AutoMigrate(&models.User{}, &models.File{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
	
	user := models.User{Username: "example", Email: "user@example.com"}
	if result := db.Create(&user); result.Error != nil {
		log.Printf("failed to create user: %v", result.Error)
	}
	
	var file models.File
	if result := db.First(&file, "name = ?", "example.txt"); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("file not found: %v", result.Error)
		} else {
			log.Printf("failed to find file: %v", result.Error)
		}
	}
}