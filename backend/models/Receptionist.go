package models

type Receptionist struct {
	ID     string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	UserID string `gorm:"uniqueIndex"`
	User   User   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
