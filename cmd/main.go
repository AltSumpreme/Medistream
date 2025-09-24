package main

import (
	"log"

	"time"

	"github.com/AltSumpreme/Medistream.git/config"
	"github.com/AltSumpreme/Medistream.git/metrics"
	"github.com/AltSumpreme/Medistream.git/routes"
	"github.com/AltSumpreme/Medistream.git/services/cache"
	"github.com/AltSumpreme/Medistream.git/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	// Initialize the logger
	utils.InitLogger()
	// Initialize metrics
	metrics.MetricsInit()
	// Initialize the database connection
	config.ConnectDB()

	// Initialize Redis
	config.InitRedis()

	// Set Gin to release mode in production
	if gin.Mode() != gin.DebugMode {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowCredentials: true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		MaxAge:           12 * time.Hour,
	}))
	router.RedirectTrailingSlash = false

	appointmentCache := cache.NewCache(config.Rdb, config.Ctx)
	medicalrecordsCache := cache.NewCache(config.Rdb, config.Ctx)
	prescriptionsCache := cache.NewCache(config.Rdb, config.Ctx)
	reportsCache := cache.NewCache(config.Rdb, config.Ctx)

	routes.RegisterRoutes(router, appointmentCache, medicalrecordsCache, prescriptionsCache, reportsCache)

	log.Println("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
