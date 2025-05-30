package main

import (
	"log"

	"github.com/AltSumpreme/Medistream.git/config"
	"github.com/AltSumpreme/Medistream.git/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	// Initialize the database connection
	config.ConnectDB()

	router := gin.Default()

	routes.RegisterRoutes(router)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
