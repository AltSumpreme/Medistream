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

	if time.Now().After(appointment.AppointmentDate) {
		return errors.New("invalid date: appointment date must be in the future")

	}

	endtime := appointment.AppointmentDate.Add(time.Duration(appointment.Duration) * time.Minute)

	var conflict int64

	if err := config.DB.Model(&models.Appointment{}).
		Where("doctor_id = ? AND appointment_date < ? AND (appointment_date + interval '1 minute' * duration) > ?", appointment.DoctorID, endtime, appointment.AppointmentDate).
		Count(&conflict).Error; err != nil {
		return errors.New("database error: " + err.Error())
	}

	if conflict > 0 {
		return errors.New("appointment conflict: another appointment exists during this time")
	}

	return nil
}
