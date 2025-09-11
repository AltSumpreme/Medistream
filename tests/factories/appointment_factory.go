package factories

import (
	"log"
	"time"

	"github.com/AltSumpreme/Medistream.git/config"
	"github.com/AltSumpreme/Medistream.git/models"
	"github.com/google/uuid"
)

func CreateAppointment(patientID, doctorID uuid.UUID) models.Appointment {
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

	if err := config.DB.Create(&appt).Error; err != nil {
		log.Fatalf("failed to create appointment: %v", err)
	}

	return appt
}
