package models

import (
	"time"

	"github.com/google/uuid"
)

type MedicalRecord struct {
	ID        uuid.UUID  `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	PatientID uuid.UUID  `gorm:"not null"`
	DoctorID  uuid.UUID  `gorm:"not null"`
	Diagnosis string     `gorm:"type:text;not null"`
	Notes     string     `gorm:"type:text"`
	CreatedAt time.Time  `gorm:"autoCreateTime"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime"`
	DeletedAt *time.Time `gorm:"index"`

	Patient Patient `gorm:"foreignKey:PatientID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Doctor  Doctor  `gorm:"foreignKey:DoctorID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Vitals  []Vital `gorm:"foreignKey:MedicalRecordID"`
}
