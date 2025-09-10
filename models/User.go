package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	AuthID    uuid.UUID `gorm:"type:uuid;not null;unique"`
	FirstName string    `gorm:"not null"`
	LastName  string    `gorm:"not null"`
	Role      Role      `gorm:"type:role;default:'PATIENT'"`
	Phone     string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	Auth             Auth `gorm:"foreignKey:AuthID;references:ID;constraint:OnDelete:CASCADE"`
	Patient          *Patient
	Doctor           *Doctor
	Receptionist     *Receptionist
	MessagesSent     []Message `gorm:"foreignKey:SenderID"`
	MessagesReceived []Message `gorm:"foreignKey:ReceiverID"`
}
