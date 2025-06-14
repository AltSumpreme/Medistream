package models

import "time"

type Appointment struct {
	ID        string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	PatientID string
	DoctorID  string
	Date      time.Time
	Status    AppointmentStatus
	Duration  int // in minutes
	Location  string
	Mode      string // "Online" or "In-Person"

	Patient Patient
	Doctor  Doctor
}
