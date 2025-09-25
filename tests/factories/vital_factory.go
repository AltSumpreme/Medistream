package factories

import (
	"log"
	"time"

	"github.com/AltSumpreme/Medistream.git/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SeedVital(db *gorm.DB, patientID uuid.UUID) models.Vital {

	if db == nil {
		log.Fatal("db instance is nil")
	}
	vital := models.Vital{
		ID:        uuid.New(),
		PatientID: patientID,

		Type:       "HEART_RATE",
		Value:      "85",
		Status:     "normal",
		RecordedAt: time.Now(),
	}

	if err := db.Create(&vital).Error; err != nil {
		log.Fatalf("failed to create vital: %v", err)
	}
	return vital
}
