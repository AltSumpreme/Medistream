package utils

import (
	"errors"
	"time"

	"github.com/AltSumpreme/Medistream.git/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func ScheduleAppointment(db *gorm.DB, doctorID uuid.UUID, patientID uuid.UUID, appointmentDate time.Time, start string, end string, excludeID *uuid.UUID) error {

	layoutTime := "15:04"
	startTime, err1 := time.Parse(layoutTime, start)
	endTime, err2 := time.Parse(layoutTime, end)
	if err1 != nil || err2 != nil {
		return errors.New("invalid time format, expected HH:MM")
	}

	if !endTime.After(startTime) {
		return errors.New("end time must be after start time")
	}

	// will later have an overlap buffer for after every 5 appointments have a break of 15mins to half and hour
	overlapCondition := "DATE(appointment_date) = ? AND ((start_time < ? AND end_time > ?) OR (start_time >= ? AND start_time < ?))"
	var count int64
	query := db.Model(&models.Appointment{}).
		Where("patient_id = ?", patientID).
		Where(overlapCondition, appointmentDate, endTime, startTime, startTime, endTime)

	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return errors.New("time slot already booked for this doctor")
	}

	query = db.Model(&models.Appointment{}).Where("doctor_id=?", doctorID).Where(overlapCondition, appointmentDate, endTime, startTime, startTime, endTime)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return errors.New("time slot already booked for this patient")
	}

	return nil
}
