package factories

import (
	"log"
	"time"

	"github.com/AltSumpreme/Medistream.git/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateAppointment(db *gorm.DB, patientID, doctorID uuid.UUID) models.Appointment {
	appt := models.Appointment{
		ID:              uuid.New(),
		PatientID:       patientID,
		DoctorID:        doctorID,
		AppointmentDate: time.Now().Add(48 * time.Hour),
		StartTime:       "10:00",
		EndTime:         "10:30",
		Notes:           "Initial Test Appointment",
		Location:        "Test Room",
		Mode:            "Online",
		Status:          models.AppointmentStatus("PENDING"),
		AppointmentType: models.ApptType("CONSULTATION"),
	}

	if err := db.Create(&appt).Error; err != nil {
		log.Fatalf("failed to create appointment: %v", err)
	}

	return appt
}
