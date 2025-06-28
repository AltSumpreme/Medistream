package main

import (
	"log"
	"time"

	"github.com/AltSumpreme/Medistream.git/config"
	"github.com/AltSumpreme/Medistream.git/routes"
	"github.com/AltSumpreme/Medistream.git/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	// Initialize the logger
	utils.InitLogger()
	// Initialize the database connection
	config.ConnectDB()

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"}, // or "*" for all origins
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
		MaxAge:       12 * time.Hour,
	}))

	routes.RegisterRoutes(router)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
