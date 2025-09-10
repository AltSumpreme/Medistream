package config

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/AltSumpreme/Medistream.git/models"
)

var DB *gorm.DB

func ConnectDB() {

	dsn := os.Getenv("DATABASE_URL")

	gormLogger := logger.Default.LogMode(logger.Silent)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: gormLogger})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v\n", err)

	}
	DB = db
	log.Println("Database connection established successfully")

	// Auto-migrate the models
	if err := DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("Failed to auto-migrate database: %v\n", err)
	} else {
		log.Println("Database auto-migration completed successfully")
	}
}
