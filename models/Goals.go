package models

import (
	"time"

	"github.com/google/uuid"
)

type HealthGoal struct {
	ID           uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	PatientID    uuid.UUID `gorm:"not null"`
	Type         GoalType  `gorm:"not null"`
	TargetValue  int       `gorm:"not null"`
	CurrentValue int       `gorm:"not null"`
	UpdatedAt    time.Time `gorm:"not null"`

	Patient Patient `gorm:"foreignKey:PatientID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
