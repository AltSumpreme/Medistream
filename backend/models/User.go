package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	FirstName string    `gorm:"not null"`
	LastName  string    `gorm:"not null"`
	Email     string    `gorm:"unique;not null"`
	Password  string    `gorm:"not null"`
	Role      Role      `gorm:"type:role;default:'PATIENT'"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	Patient          *Patient
	Doctor           *Doctor
	Receptionist     *Receptionist
	MessagesSent     []Message `gorm:"foreignKey:SenderID"`
	MessagesReceived []Message `gorm:"foreignKey:ReceiverID"`
}
