package models

import (
	"github.com/google/uuid"
)

type Report struct {
	ID              uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Title           string     `json:"title" binding:"required"`
	Description     string     `json:"description" binding:"required"`
	FileURL         string     `json:"file_url" binding:"required"`
	DoctorID        uuid.UUID  `gorm:"type:uuid;not null" json:"doctor_id"`
	PatientID       uuid.UUID  `gorm:"type:uuid;not null" json:"patient_id"`
	MedicalRecordID *uuid.UUID `gorm:"type:uuid;" json:"medical_record_id"`

	MedicalRecord MedicalRecord `gorm:"foreignKey:MedicalRecordID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Patient       Patient       `gorm:"foreignKey:PatientID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Doctor        Doctor        `gorm:"foreignKey:DoctorID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
