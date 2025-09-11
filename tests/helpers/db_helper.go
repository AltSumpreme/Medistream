package helpers

import (
	"database/sql"
	"log"
	"os"

	"github.com/AltSumpreme/Medistream.git/config"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	TestDB     *gorm.DB
	originalDB *gorm.DB
)

// SetupTestDatabase connects to the test DB and migrates all models

func SetupTestDatabase() *gorm.DB {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		log.Fatal("TEST_DATABASE_URL not set")
	}

	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}

	// Create extension for UUIDs
	_, err = sqlDB.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	if err != nil {
		log.Fatalf("Failed to create uuid-ossp extension: %v", err)
	}

	// Run Goose migrations
	migrationsDir := "../../migrations"
	if err := goose.Up(sqlDB, migrationsDir); err != nil {
		log.Fatalf("Failed to migrate test database: %v", err)
	}
	gormLogger := logger.Default.LogMode(logger.Silent)
	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{Logger: gormLogger})
	if err != nil {
		log.Fatalf("Failed to initialize GORM for test database: %v", err)
	}
	TestDB = gormDB
	return gormDB
}

// PatchDatabase temporarily replaces the global DB with the test DB
func PatchDatabase() {
	if config.DB == nil {
		log.Fatal("global config.DB is nil; cannot patch")
	}

	originalDB = config.DB
	config.DB = TestDB
}

// UnpatchDatabase restores the original global DB
func UnpatchDatabase() {
	if originalDB != nil {
		config.DB = originalDB
		log.Println("Restored original DB")
	}
}

// TearDownTestDatabase closes the test DB connection
func TearDownTestDatabase() {
	if TestDB == nil {
		return
	}
	sqlDB, err := TestDB.DB()
	if err != nil {
		log.Printf("failed to get sql.DB from test DB: %v", err)
		return
	}
	sqlDB.Close()
	log.Println("Test database connection closed")
}
