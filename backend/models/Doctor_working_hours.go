package models

import (
	"time"

	"github.com/google/uuid"
)

type Doctor_working_hours struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	DoctorID  uuid.UUID `gorm:"type:uuid;not null"`
	Weekday   int       `gorm:"type:int;not null"`
	StartTime string    `gorm:"type:time;not null"`
	EndTime   string    `gorm:"type:time;not null"`
	IsActive  bool      `gorm:"type:boolean;not null" default:"true"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
