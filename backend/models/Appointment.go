package models

import (
	"time"

	"github.com/google/uuid"
)

type Appointment struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	PatientID uuid.UUID `gorm:"type:uuid;not null"`
	DoctorID  uuid.UUID `gorm:"type:uuid;not null"`
	Date      time.Time
	Status    AppointmentStatus
	Duration  int // in minutes
	Location  string
	Mode      string // "Online" or "In-Person"

	Patient Patient
	Doctor  Doctor
}
