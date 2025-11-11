package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/pressly/goose"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)

	var sqlDB *sql.DB
	var err error

	for i := 1; i <= 10; i++ {
		sqlDB, err = sql.Open("postgres", dsn)
		if err == nil && sqlDB.Ping() == nil {
			log.Println("Database is ready (attempt", i, ")")
			break
		}
		log.Printf("Database not ready (attempt %d/10), retrying in 5s...\n", i)
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		log.Fatalf("failed to connect to database after retries: %v", err)
	}

	// Run Goose migrations
	migrationsDir := os.Getenv("MIGRATIONS_DIR")
	if err := goose.Up(sqlDB, migrationsDir); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}
	log.Println("Database migrations completed successfully")

	// Initialize GORM
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("failed to initialize GORM: %v", err)
	}

	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(50)
	sqlDB.SetConnMaxLifetime(time.Hour)
	DB = gormDB
	log.Println("DB initialized successfully")
}
func CloseDB() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("failed to get sql.DB from GORM: %v", err)
	}
	if err := sqlDB.Close(); err != nil {
		log.Fatalf("failed to close database connection: %v", err)
	}
	log.Println("Database connection closed successfully")
}
