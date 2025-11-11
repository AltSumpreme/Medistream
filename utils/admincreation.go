package utils

import (
	"github.com/AltSumpreme/Medistream.git/config"
	"github.com/AltSumpreme/Medistream.git/models"
	"github.com/google/uuid"
)

func CreateAdminUserIfNotExists() error {

	adminEmail := GetEnvWithDefault("ADMIN_EMAIL", "")
	adminPassword := GetEnvWithDefault("ADMIN_PASSWORD", "")
	adminFirstName := GetEnvWithDefault("ADMIN_FIRST_NAME", "")
	adminLastName := GetEnvWithDefault("ADMIN_LAST_NAME", "")
	adminPhone := GetEnvWithDefault("ADMIN_PHONE", "")

	if adminEmail == "" || adminPassword == "" {
		Log.Info("Admin credentials not provided, skipping admin user creation")
		return nil // No admin credentials provided

	}

	hashedPassword, err := HashPassword(adminPassword)
	if err != nil {
		return err
	}

	var existingAdmin models.Auth

	err = config.DB.Where("email =?", adminEmail).First(&existingAdmin).Error
	if err == nil {
		Log.Info("Admin user already exists, skipping creation")
		return nil // Admin user already exists
	}

	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	newAdmin := models.Auth{
		ID:       uuid.New(),
		Email:    adminEmail,
		Password: hashedPassword,
	}
	if err := tx.Create(&newAdmin).Error; err != nil {
		tx.Rollback()
		return err
	}

	adminProfile := models.User{
		ID:        uuid.New(),
		AuthID:    newAdmin.ID,
		FirstName: adminFirstName,
		LastName:  adminLastName,
		Phone:     adminPhone,
		Role:      models.RoleAdmin,
	}
	if err := tx.Create(&adminProfile).Error; err != nil {
		tx.Rollback()
		return err
	}

	Log.Info("Admin user created successfully with email")
	return tx.Commit().Error

}
