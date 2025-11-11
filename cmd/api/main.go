package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	"time"

	"github.com/AltSumpreme/Medistream.git/config"
	"github.com/AltSumpreme/Medistream.git/metrics"
	"github.com/AltSumpreme/Medistream.git/queue"
	"github.com/AltSumpreme/Medistream.git/routes"
	"github.com/AltSumpreme/Medistream.git/services/cache"
	"github.com/AltSumpreme/Medistream.git/services/mail"
	"github.com/AltSumpreme/Medistream.git/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	// Initialize the logger
	utils.InitLogger()
	// Initialize metrics
	metrics.MetricsInit()
	// Initialize the database connection
	config.ConnectDB()

	defer config.CloseDB()

	// Initialize Redis
	config.InitRedis()

	// Initialize AWS S3
	config.InitS3()

	//Initialize Job Queue
	config.InitAsynqQueue()
	jobQueue := queue.Init()
	defer queue.Close()

	// Initialize SMTP Mailer
	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		log.Fatalf("Invalid SMTP_PORT: %v", err)
	}
	mail.InitMailer(mail.MailerConfig{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     port,
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
		From:     os.Getenv("SMTP_FROM"),
	})

	// Set Gin to release mode in production
	if gin.Mode() != gin.DebugMode {
		gin.SetMode(gin.ReleaseMode)
	}

	err = utils.CreateAdminUserIfNotExists()
	if err != nil {
		utils.Log.Warnf("Failed to create admin user: %v", err)
	}

	router := gin.Default()

	origins := os.Getenv("CORS_ALLOW_ORIGINS")
	if origins == "" {
		origins = "*"
	}

	router.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Split(origins, ","),
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
	vitalsCache := cache.NewCache(config.Rdb, config.Ctx)

	routes.RegisterRoutes(router, appointmentCache, medicalrecordsCache, prescriptionsCache, reportsCache, vitalsCache, jobQueue)

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	log.Println("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
