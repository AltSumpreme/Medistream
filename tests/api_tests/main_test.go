package apitests

import (
	"log"
	"os"
	"testing"

	"github.com/AltSumpreme/Medistream.git/config"
	"github.com/AltSumpreme/Medistream.git/tests/helpers"
	"github.com/AltSumpreme/Medistream.git/utils"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {

	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalf("failed to load .env: %v", err)
	}
	config.ConnectDB()
	helpers.SetupTestDatabase()
	helpers.PatchDatabase()
	utils.InitLogger()

	gin.SetMode(gin.TestMode)
	// Run Tests
	code := m.Run()

	helpers.UnpatchDatabase()
	helpers.TearDownTestDatabase()

	// Exit with status code
	os.Exit(code)
}
