package appointments

import (
	"net/http"
	"strconv"
	"time"

	"github.com/AltSumpreme/Medistream.git/config"
	"github.com/AltSumpreme/Medistream.git/models"
	"github.com/AltSumpreme/Medistream.git/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AppointmentInput struct {
	PatientID       uuid.UUID `json:"patient_id" binding:"required"`
	DoctorID        uuid.UUID `json:"doctor_id" binding:"required"`
	AppointmentDate time.Time `json:"appointment_date" binding:"required"`
	Duration        int       `json:"duration" binding:"required"`
	AppointmentType string    `json:"appointment_type" binding:"required,oneof=CONSULTATION FOLLOWUP CHECKUP EMERGENCY"`
	Location        string    `json:"location" binding:"required"`
	Mode            string    `json:"mode" binding:"required,oneof=Online In-Person"`
	Notes           string    `json:"notes"`
}
type AppointmentStatusInput struct {
	Status string `json:"status" binding:"required,oneof=SCHEDULED CONFIRMED CANCELLED COMPLETED"`
}

func CreateAppointment(c *gin.Context) {

	var input AppointmentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Log.Warnf("CreateAppointment: Invalid input - %v", err)
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}
	appointment := models.Appointment{
		PatientID:       input.PatientID,
		DoctorID:        input.DoctorID,
		AppointmentDate: input.AppointmentDate,
		Status:          models.AppointmentStatusPending,
		Duration:        input.Duration,
		AppointmentType: models.ApptType(input.AppointmentType),
		Notes:           input.Notes,
	}

	scheduleErr := utils.ScheduleAppointment(appointment)
	if scheduleErr != nil {
		utils.Log.Warnf("CreateAppointment: Scheduling error - %v", scheduleErr)
		c.JSON(400, gin.H{"error": scheduleErr.Error()})
		return
	}

	err := config.DB.Create(&appointment).Error
	if err != nil {
		utils.Log.Errorf("CreateAppointment: Database error - %v", err)
		c.JSON(500, gin.H{"error": "Failed to create appointment - " + err.Error()})
		return
	}
	utils.Log.Infof("CreateAppointment: Appointment created successfully with ID %s", appointment.ID)
	c.JSON(http.StatusCreated, gin.H{"message": "Appointment created successfully"})

}

func GetAllAppointments(c *gin.Context) {
	limit := 10
	page := 1
	Maxlimit := 100
	if l := c.Query("limit"); l != "" {
		if n, err := strconv.Atoi(l); err == nil && n > 0 && n <= Maxlimit {
			limit = n
		} else {
			limit = Maxlimit
		}

	}
	if p := c.Query("page"); p != "" {
		if n, err := strconv.Atoi(p); err == nil {
			page = n
		}
	}
	var appointments []models.Appointment
	if err := config.DB.WithContext(c.Request.Context()).Preload("Patient").Preload("Doctor").Limit(limit).Offset((page - 1) * limit).Find(&appointments).Error; err != nil {
		utils.Log.Errorf("GetAppointments: Database error - %v", err)
		c.JSON(500, gin.H{"error": "Failed to retrieve appointments - " + err.Error()})
		return
	}
	utils.Log.Infof("GetAppointments: Retrieved %d appointments for page %d with limit %d", len(appointments), page, limit)
	c.JSON(200, gin.H{"appointments": appointments, "page": page, "limit": limit, "total": len(appointments)})
}

func GetAppointmentByID(c *gin.Context) {

	appointmentID := c.Param("id")
	user, _ := utils.GetCurrentUser(c)
	var appointment models.Appointment

	if err := config.DB.WithContext(c.Request.Context()).Preload("Patient").Preload("Doctor").Where("id = ?", appointmentID).First(&appointment).Error; err != nil {
		utils.Log.Errorf("GetAppointmentByID: Appointment not found - %v", err)
		c.JSON(404, gin.H{"error": "Appointment not found - " + err.Error()})
		return
	}

	if user.Role != models.RoleAdmin && (appointment.PatientID != user.ID || appointment.DoctorID != user.ID) {
		utils.Log.Warnf("GetAppointmentByID: You are not authorised to access this appointment")
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.JSON(200, gin.H{"appointment": appointment})
}

func UpdateAppointment(c *gin.Context) {
	appointmentID := c.Param("id")

	user, err := utils.GetCurrentUser(c)
	if err != nil {
		utils.Log.Warnf("UpdateAppointment: Unauthorized access - %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: " + err.Error()})
		return
	}

	// Fetch appointment
	var appt models.Appointment
	if err := config.DB.WithContext(c).First(&appt, "id = ?", appointmentID).Error; err != nil {
		utils.Log.Errorf("UpdateAppointment: Appointment not found - %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found"})
		return
	}
	switch {
	case user.Role == models.RolePatient:
		if user.ID != appt.PatientID {
			utils.Log.Warnf("UpdateAppointment:Unauthorized access")
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to update this appointment"})
			return
		}
		if appt.Status == "ACCEPTED" {
			utils.Log.Warnf("UpdateAppointment:Cannot update an accepted appointment")
			c.JSON(http.StatusForbidden, gin.H{"error": "You cannot update an accepted appointment"})
		}
	case user.Role != models.RoleAdmin && !(user.Role == models.RoleDoctor && user.ID == appt.DoctorID):
		utils.Log.Warnf("UpdateAppointment:Unauthorized access")
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to update this appointment"})
		return
	}

	var input AppointmentInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	appt.PatientID = input.PatientID
	appt.DoctorID = input.DoctorID
	appt.AppointmentDate = input.AppointmentDate
	appt.Duration = input.Duration
	appt.AppointmentType = models.ApptType(input.AppointmentType)
	appt.Mode = input.Mode
	appt.Notes = input.Notes

	if input.Location != "" {
		if user.Role == "PATIENT" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Only doctor or admin can update location"})
			return
		}
	}
	appt.Location = input.Location

	// Schedule conflict check
	if err := utils.ScheduleAppointment(appt); err != nil {
		utils.Log.Warnf("UpdateAppointment:Conflict in scheduling the appointment")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Save(&appt).Error; err != nil {
		utils.Log.Errorf("UpdateAppointment:Failed to update appointment")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update appointment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Appointment updated successfully", "appointment": appt})
}

func DeleteAppointment(c *gin.Context) {
	appointmentId := c.Param("id")
	var appointment models.Appointment
	if err := config.DB.WithContext(c.Request.Context()).Where("id = ?", appointmentId).First(&appointment).Error; err != nil {
		utils.Log.Errorf("DeleteAppointment: Appointment not found - %v", err)
		c.JSON(404, gin.H{"error": "Appointment not found - " + err.Error()})
		return
	}
	if err := config.DB.WithContext(c.Request.Context()).Delete(&appointment).Error; err != nil {
		utils.Log.Errorf("DeleteAppointment: Failed to delete appointment - %v", err)
		c.JSON(500, gin.H{"error": "Failed to delete appointment - " + err.Error()})
		return
	}
	utils.Log.Infof("DeleteAppointment: Appointment with ID %s deleted successfully", appointment.ID)
	c.JSON(200, gin.H{"message": "Appointment deleted successfully"})
}

func ChangeAppointmentStatus(c *gin.Context) {
	appointmentID := c.Param("id")

	_, err := utils.GetCurrentUser(c)

	if err != nil {
		c.JSON(403, gin.H{"error": "Unauthorized: " + err.Error()})
	}

	var appointment models.Appointment

	if err := config.DB.WithContext(c.Request.Context()).First(&appointment, "id = ?", appointmentID).Error; err != nil {
		utils.Log.Errorf("ChangeAppointmentStatus: Appointment not found - %v", err)
		c.JSON(404, gin.H{"error": "Appointment not found - " + err.Error()})
		return
	}

	var input AppointmentStatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Log.Warnf("ChangeAppointmentStatus: Invalid input - %v", err)
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}
	appointment.Status = models.AppointmentStatus(input.Status)
	if err := config.DB.WithContext(c.Request.Context()).Save(&appointment).Error; err != nil {
		utils.Log.Errorf("ChangeAppointmentStatus: Failed to update appointment status - %v", err)
		c.JSON(500, gin.H{"error": "Failed to update appointment status - " + err.Error()})
		return
	}
}

func RescheduleAppointment(c *gin.Context) {
	appointmentID := c.Param("id")
	user, _ := utils.GetCurrentUser(c)

	var appointment models.Appointment
	if err := config.DB.First(&appointment, "id = ?", appointmentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found"})
		return
	}

	if user.Role != models.RoleAdmin && (user.ID != appointment.PatientID || user.ID != appointment.DoctorID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only reschedule your own appointment"})
		return
	}

	if appointment.Status != "PENDING" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only PENDING appointments can be rescheduled"})
		return
	}

	var input struct {
		Date     time.Time `json:"date" binding:"required"`
		Duration int       `json:"duration" binding:"required"`
		Mode     string    `json:"mode" binding:"required,oneof=Online In-Person"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	appointment.AppointmentDate = input.Date
	appointment.Duration = input.Duration
	appointment.Mode = input.Mode

	if err := utils.ScheduleAppointment(appointment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Save(&appointment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reschedule"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Appointment rescheduled", "appointment": appointment})
}

func CancelAppointment(c *gin.Context) {
	appointmentID := c.Param("id")
	user, exists := utils.GetCurrentUser(c)
	if exists != nil {
		utils.Log.Warnf("CancelAppointment:Unauthorised access to the route")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You have not been authenticated"})
	}

	var appointment models.Appointment
	if err := config.DB.First(&appointment, "id = ?", appointmentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found"})
		return
	}

	allowed := user.Role == models.RoleAdmin || user.ID == appointment.PatientID
	if !allowed {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		return
	}

	appointment.Status = "CANCELLED"
	if err := config.DB.Save(&appointment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Appointment cancelled"})
}

func GetAppointmentByDoctorID(c *gin.Context) {
	doctorID := c.Param("id")
	user, err := utils.GetCurrentUser(c)
	if err != nil {
		utils.Log.Warnf("GetAppointmentByDoctorID: Unauthorized access - %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if user.Role == models.RoleDoctor && user.ID.String() != doctorID {
		utils.Log.Warnf("GetAppointmentByDoctorID: Doctor %s attempted to access data of doctor %s", user.ID, doctorID)
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	limit := 10
	offset := 0
	const MaxLimit = 100

	if l := c.Query("limit"); l != "" {
		if n, err := strconv.Atoi(l); err == nil && n > 0 && n <= MaxLimit {
			limit = n
		} else {
			limit = MaxLimit
		}
	}
	if o := c.Query("offset"); o != "" {
		if n, err := strconv.Atoi(o); err == nil && n >= 0 {
			offset = n
		}
	}

	var appointments []models.Appointment
	if err := config.DB.WithContext(c.Request.Context()).
		Where("doctor_id = ?", doctorID).
		Order("appointment_date desc").
		Limit(limit).
		Offset(offset).
		Find(&appointments).Error; err != nil {
		utils.Log.Errorf("GetAppointmentByDoctorID: Failed to fetch appointments - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch appointments"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"appointments": appointments})
}

func GetAppointmentByPatientID(c *gin.Context) {
	patientID := c.Param("id")
	user, err := utils.GetCurrentUser(c)
	if err != nil {
		utils.Log.Warnf("GetAppointmentByPatientID: Unauthorized access attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if user.Role == models.RolePatient && user.ID.String() != patientID {
		utils.Log.Warnf("GetAppointmentByPatientID: Access denied for patient %s to data of %s", user.ID, patientID)
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	limit := 10
	offset := 0
	const MaxLimit = 100

	if l := c.Query("limit"); l != "" {
		if n, err := strconv.Atoi(l); err == nil && n > 0 && n <= MaxLimit {
			limit = n
		} else {
			limit = MaxLimit
		}
	}

	if o := c.Query("offset"); o != "" {
		if n, err := strconv.Atoi(o); err == nil && n >= 0 {
			offset = n
		}
	}

	var appointments []models.Appointment
	db := config.DB.WithContext(c.Request.Context()).Preload("Patient").
		Where("patient_id = ?", patientID).
		Order("appointment_date desc").
		Limit(limit).
		Offset(offset)

	if user.Role == models.RoleDoctor {
		db = db.Preload("Doctor").Where("doctor_id = ?", user.ID)
	}

	result := db.Find(&appointments)
	if err := result.Error; err != nil {
		utils.Log.Errorf("GetAppointmentByPatientID: Failed to fetch appointments - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch appointments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"appointments": appointments})
}
