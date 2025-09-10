package utils

import (
	"github.com/AltSumpreme/Medistream.git/config"
	"github.com/AltSumpreme/Medistream.git/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func DoctorHasAccessToPatient(doctorID, patientID uuid.UUID, c *gin.Context) (bool, error) {
	var count int64
	err := config.DB.WithContext(c).
		Model(&models.Appointment{}).
		Where("doctor_id = ? AND patient_id = ?", doctorID, patientID).
		Where("status != ?", models.AppointmentStatusCancelled).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
