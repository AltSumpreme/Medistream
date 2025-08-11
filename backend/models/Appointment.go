package models

import (
	"time"

	"github.com/google/uuid"
)

type Appointment struct {
	ID              uuid.UUID         `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	PatientID       uuid.UUID         `gorm:"type:uuid;not null"`
	DoctorID        uuid.UUID         `gorm:"type:uuid;not null"`
	AppointmentDate time.Time         `gorm:"column:appointment_date;not null"`
	Status          AppointmentStatus `gorm:"type:appointment_status;not null" default:"PENDING"`
	Duration        int               // in minutes
	Location        string
	Mode            string    `gorm:"type:mode;not null" default:"Online"` // Online or In-Person
	AppointmentType ApptType  `gorm:"type:appt_type;not null" default:"CONSULTATION"`
	Notes           string    `gorm:"type:text"`
	CreatedAt       time.Time `gorm:"autoCreateTime"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`
	Patient         Patient
	Doctor          Doctor
}
