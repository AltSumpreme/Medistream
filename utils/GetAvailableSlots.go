package utils

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetAvailableSlots(db *gorm.DB, doctorID uuid.UUID, appointmentDate time.Time) ([]string, error) {
	var workingHours struct {
		StartTime string
		EndTime   string
	}

	weekday := int(appointmentDate.Weekday()) // 0 = Sunday
	err := db.Raw(`
		SELECT start_time, end_time 
		FROM doctor_working_hours 
		WHERE doctor_id = ? AND weekday = ? AND is_active = true
	`, doctorID, weekday).Scan(&workingHours).Error
	if err != nil {
		return nil, err
	}
	if workingHours.StartTime == "" || workingHours.EndTime == "" {
		return []string{}, nil
	}

	// Parse start/end times
	layout := "15:04"
	startTime, err1 := time.Parse(layout, workingHours.StartTime)
	endTime, err2 := time.Parse(layout, workingHours.EndTime)
	if err1 != nil || err2 != nil {
		return nil, err
	}

	// Fetch booked slots
	var bookedSlots []struct {
		StartTime string
		EndTime   string
	}
	err = db.Raw(`
		SELECT start_time, end_time 
		FROM appointments 
		WHERE doctor_id = ? AND DATE(appointment_date) = ?
	`, doctorID, appointmentDate.Format("2006-01-02")).Scan(&bookedSlots).Error
	if err != nil {
		return nil, err
	}

	booked := make(map[string]bool)
	for _, slot := range bookedSlots {
		booked[slot.StartTime] = true
	}
	// Generate 30-min interval slots
	var available []string
	for t := startTime; t.Add(30*time.Minute).Before(endTime) || t.Add(30*time.Minute).Equal(endTime); t = t.Add(30 * time.Minute) {
		slot := t.Format(layout)
		if !booked[slot] {
			available = append(available, slot)
		}
	}

	return available, nil
}
