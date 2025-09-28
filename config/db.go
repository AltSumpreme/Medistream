package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/pressly/goose"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s&channel_binding=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
		os.Getenv("DB_CHANNEL_BINDING"),
	)
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	// Open sql.DB for Goose
	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to open sql.DB: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}
	log.Println("Database connection established (sql.DB)")

	// Open sql.DB for Goose
	sqlDB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to open sql.DB: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}
	log.Println("Database connection established (sql.DB)")

	// Run Goose migrations
	migrationsDir := os.Getenv("MIGRATIONS_DIR")
	if err := goose.Up(sqlDB, migrationsDir); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}
	log.Println("Database migrations completed successfully")

	// Initialize GORM
	gormLogger := logger.Default.LogMode(logger.Silent)
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		log.Fatalf("failed to initialize GORM: %v", err)
	}

	DB = gormDB // <- Now factories can safely use this
	log.Println("GORM DB initialized successfully")
}
