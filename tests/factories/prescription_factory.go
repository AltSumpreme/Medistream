package factories

import (
	"log"
	"time"

	"github.com/AltSumpreme/Medistream.git/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SeedPrescription(db *gorm.DB, medicalRecordID *uuid.UUID, doctorID, patientID uuid.UUID) models.Prescription {

	if db == nil {
		log.Fatal("db instance is nil")
	}

	prescription := models.Prescription{
		ID:              uuid.New(),
		PatientID:       patientID,
		DoctorID:        doctorID,
		MedicalRecordID: medicalRecordID,
		Medication:      "TestMed",
		Dosage:          "1x daily",
		Instructions:    "After food",
		IssuedAt:        time.Now(),
	}
	if err := db.Save(&prescription).Error; err != nil {
		log.Fatalf("Prescription not created")
	}

	return prescription

}

func SeedPrescriptionWithMedicalRecord(db *gorm.DB, doctorID, patientID uuid.UUID) models.Prescription {

	if db == nil {
		log.Fatal("db instance is nil")
	}
	medicalRecord := CreateMedicalRecord(db, patientID, doctorID)
	return SeedPrescription(db, &medicalRecord.ID, doctorID, patientID)
}

func SeedPrescriptionWithoutMedicalRecord(db *gorm.DB, doctorID, patientID uuid.UUID) models.Prescription {
	if db == nil {
		log.Fatal("db instance is nil")
	}
	return SeedPrescription(db, nil, doctorID, patientID)
}
