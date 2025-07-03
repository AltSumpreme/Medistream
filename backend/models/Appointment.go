package models

import (
	"time"

	"github.com/google/uuid"
)

type Appointment struct {
	ID              uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	PatientID       uuid.UUID `gorm:"type:uuid;not null"`
	DoctorID        uuid.UUID `gorm:"type:uuid;not null"`
	AppointmentDate time.Time `gorm:"column:appointment_date;not null"`
	Status          AppointmentStatus
	Duration        int       // in minutes
	Location        string    // e.g., "Room 101", "Online"
	Mode            string    // "Online" or "In-Person"
	CreatedAt       time.Time `gorm:"autoCreateTime"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`
	Patient         Patient
	Doctor          Doctor
}
