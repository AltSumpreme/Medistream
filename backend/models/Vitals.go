package models

import (
	"time"

	"github.com/google/uuid"
)

type Vital struct {
	ID              uuid.UUID  `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	PatientID       uuid.UUID  `gorm:"not null"`
	Type            VitalType  `gorm:"not null"`
	Value           string     `gorm:"not null"`
	Status          string     `gorm:"not null"`
	RecordedAt      time.Time  `gorm:"not null"`
	MedicalRecordID *uuid.UUID `gorm:"type:uuid"`
	CreatedAt       time.Time  `gorm:"autoCreateTime"`
	UpdatedAt       time.Time  `gorm:"autoUpdateTime"`

	Patient       Patient       `gorm:"foreignKey:PatientID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	MedicalRecord MedicalRecord `gorm:"foreignKey:MedicalRecordID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
