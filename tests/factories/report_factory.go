package factories

import (
	"log"

	"github.com/AltSumpreme/Medistream.git/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SeedReport(db *gorm.DB, doctorID, patientID uuid.UUID, medicalRecordId *uuid.UUID) models.Report {

	if db == nil {
		log.Fatal("db instance is nil")
	}

	report := models.Report{
		ID:              uuid.New(),
		PatientID:       patientID,
		DoctorID:        doctorID,
		FileURL:         "http://example.com/report.pdf",
		Title:           "Test Report",
		Description:     "This is a test report",
		MedicalRecordID: medicalRecordId,
	}

	if err := db.Create(&report).Error; err != nil {
		log.Fatalf("failed to create report: %v", err)
	}

	return report
}
