package factories

import (
	"fmt"
	"log"

	"github.com/AltSumpreme/Medistream.git/models"
	"github.com/AltSumpreme/Medistream.git/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SeedUser(db *gorm.DB, role models.Role) models.User {

	if db == nil {
		log.Fatal("db instance is nil")
	}
	email := fmt.Sprintf("test_%s_%s@example.com", role, uuid.New().String())

	auth := models.Auth{
		ID:       uuid.New(),
		Email:    email,
		Password: "hashedpassword",
	}
	if err := db.Create(&auth).Error; err != nil {
		panic(fmt.Sprintf("Failed to create auth: %v", err))
	}

	user := models.User{
		ID:        uuid.New(),
		FirstName: "Test",
		LastName:  string(role) + "User",
		Phone:     "1234567890",
		Role:      role,
		AuthID:    auth.ID,
	}

	if err := db.Create(&user).Error; err != nil {
		utils.Log.Errorf("failed to create user: %v", err)
	}

	return user
}

func SeedPatient(db *gorm.DB, user models.User) models.Patient {
	if db == nil {
		log.Fatal("db instance is nil")
	}
	patient := models.Patient{
		ID:     uuid.New(),
		UserID: user.ID,
	}
	if err := db.FirstOrCreate(&patient, models.Patient{UserID: user.ID}).Error; err != nil {
		log.Fatalf("failed to seed patient: %v", err)
	}
	return patient
}

func SeedDoctor(db *gorm.DB, user models.User) models.Doctor {
	if db == nil {
		log.Fatal("db instance is nil")
	}
	doctor := models.Doctor{
		ID:             uuid.New(),
		UserID:         user.ID,
		Specialization: "GENERAL",
	}
	if err := db.FirstOrCreate(&doctor, models.Doctor{UserID: user.ID}).Error; err != nil {
		log.Fatalf("failed to seed doctor: %v", err)
	}
	return doctor
}

func CreateEntries(db *gorm.DB) (models.User, models.Patient, models.User, models.Doctor, models.User) {
	// Seed patient
	userPatient := SeedUser(db, models.RolePatient)
	patient := SeedPatient(db, userPatient)

	// Seed doctor
	userDoctor := SeedUser(db, models.RoleDoctor)
	doctor := SeedDoctor(db, userDoctor)

	userAdmin := SeedUser(db, models.RoleAdmin)

	return userPatient, patient, userDoctor, doctor, userAdmin
}
