package utils

import (
	"errors"
	"time"

	"github.com/AltSumpreme/Medistream.git/config"
	"github.com/AltSumpreme/Medistream.git/models"
)

func ScheduleAppointment(appointment models.Appointment) error {

	if appointment.Duration <= 0 {
		return errors.New("invalid duration: must be greater than 0")
	}

	if appointment.Mode != "Online" && appointment.Mode != "In-Person" {
		return errors.New("invalid mode: must be 'Online' or 'In-Person'")
	}

	if time.Now().After(appointment.Date) {
		return errors.New("invalid date: appointment date must be in the future")

	}

	endtime := appointment.Date.Add(time.Duration(appointment.Duration) * time.Minute)

	var conflict int64

	if err := config.DB.Model(&models.Appointment{}).
		Where("doctor_id = ? AND date < ? AND (date + interval '1 minute' * duration) > ?", appointment.DoctorID, endtime, appointment.Date).
		Count(&conflict).Error; err != nil {
		return errors.New("database error: " + err.Error())
	}

	if conflict > 0 {
		return errors.New("appointment conflict: another appointment exists during this time")
	}

	return nil
}
