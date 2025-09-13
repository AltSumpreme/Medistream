package factories

import (
	"log"

	"github.com/AltSumpreme/Medistream.git/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateMedicalRecord(db *gorm.DB, patientID, doctorID uuid.UUID) models.MedicalRecord {
	record := models.MedicalRecord{
		ID:        uuid.New(),
		PatientID: patientID,
		DoctorID:  doctorID,
		Diagnosis: "Test Diagnosis",
		Notes:     "Test Notes",
	}

	if err := db.Create(&record).Error; err != nil {
		log.Fatalf("failed to create medical record: %v", err)
	}

	return record
}
