package apitests

import (
	"log"
	"os"
	"testing"

	"github.com/AltSumpreme/Medistream.git/config"
	"github.com/AltSumpreme/Medistream.git/metrics"
	"github.com/AltSumpreme/Medistream.git/tests/helpers"
	"github.com/AltSumpreme/Medistream.git/utils"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func TestMain(m *testing.M) {

	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalf("failed to load .env: %v", err)
	}
	os.Setenv("MIGRATIONS_DIR", "../../migrations")

	config.ConnectDB()
	config.InitRedis()
	helpers.SetupTestDatabase()
	helpers.PatchDatabase()
	utils.InitLogger()
	gin.SetMode(gin.TestMode)
	// Initialize metrics
	metrics.MetricsInit()

	// Start a separate goroutine to serve metrics
	go func() {
		r := gin.Default()
		r.GET("/metrics", gin.WrapH(promhttp.Handler()))
		if err := r.Run("0.0.0.0:2112"); err != nil {
			log.Fatalf("failed to run metrics server: %v", err)
		}
	}()

	// Run Tests
	code := m.Run()
	helpers.UnpatchDatabase()
	helpers.TearDownTestDatabase()

	// Exit with status code
	os.Exit(code)
}
