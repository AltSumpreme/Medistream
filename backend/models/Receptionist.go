package models

import "github.com/google/uuid"

type Receptionist struct {
	ID     uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	UserID uuid.UUID `gorm:"uniqueIndex"`
	User   User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
