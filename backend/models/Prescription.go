package models

import (
	"time"

	"github.com/google/uuid"
)

type Prescription struct {
	ID              uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	PatientID       uuid.UUID
	DoctorID        uuid.UUID
	Medication      string
	Dosage          string
	Instructions    string
	IssuedAt        time.Time
	CreatedAt       time.Time  `gorm:"autoCreateTime"`
	UpdatedAt       time.Time  `gorm:"autoUpdateTime"`
	DeletedAt       *time.Time `gorm:"index"`
	MedicalRecordID *uuid.UUID `gorm:"type:uuid"`

	MedicalRecord MedicalRecord `gorm:"foreignKey:MedicalRecordID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Patient       Patient       `gorm:"foreignKey:PatientID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Doctor        Doctor        `gorm:"foreignKey:DoctorID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
