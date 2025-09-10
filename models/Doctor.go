package models

import (
	"github.com/google/uuid"
)

type Doctor struct {
	ID     uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID uuid.UUID `gorm:"uniqueIndex"`
	User   *User

	Specialization string
	Appointments   []Appointment
	Prescriptions  []Prescription
}
