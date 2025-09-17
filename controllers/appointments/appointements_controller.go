package appointments

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/AltSumpreme/Medistream.git/config"
	"github.com/AltSumpreme/Medistream.git/metrics"
	"github.com/AltSumpreme/Medistream.git/models"
	"github.com/AltSumpreme/Medistream.git/services/cache"
	"github.com/AltSumpreme/Medistream.git/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AppointmentInput struct {
	AppointmentDate time.Time `json:"appointmentDate" binding:"required"`
	AppointmentType string    `json:"appointmentType" binding:"required,oneof=CONSULTATION FOLLOWUP CHECKUP EMERGENCY"`
	StartTime       string    `json:"startTime" binding:"required"`
	EndTime         string    `json:"endTime" binding:"required"`
	Mode            string    `json:"mode" binding:"required,oneof=Online In-Person"`
	Notes           string    `json:"notes"`
	DoctorID        uuid.UUID `json:"doctorId" binding:"required"`
}
type AppointmentStatusInput struct {
	Status string `json:"status" binding:"required,oneof=SCHEDULED CONFIRMED CANCELLED COMPLETED"`
}
type AppointmentUpdateInput struct {
	StartTime string `json:"appointment_time" binding:"omitempty"`
	EndTime   string `json:"end_time" binding:"omitempty"`
	Mode      string `json:"mode" binding:"omitempty,oneof=Online In-Person"`
	Location  string `json:"location" binding:"omitempty"`
	Notes     string `json:"notes" binding:"omitempty"`
}

func CreateAppointment(c *gin.Context) {
	user, err := utils.GetCurrentUser(c)
	if err != nil {
		utils.Log.Warnf("CreateAppointment: Failed to get current user - %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input AppointmentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Log.Warnf("CreateAppointment: Invalid input - %v", err)
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	var patient models.Patient
	err = metrics.DbMetrics(config.DB, "select_patient", func(db *gorm.DB) error {
		return db.WithContext(c.Request.Context()).
			Select("id").
			Where("user_id = ?", user.UserID).
			First(&patient).Error
	})
	if err != nil {
		utils.Log.Warnf("CreateAppointment: Failed to get patient ID - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get patient ID"})
		return
	}
	patientID := patient.ID

	appointment := models.Appointment{
		PatientID:       patientID,
		DoctorID:        input.DoctorID,
		AppointmentDate: input.AppointmentDate,
		Status:          models.AppointmentStatusPending,
		StartTime:       input.StartTime,
		EndTime:         input.EndTime,
		Mode:            input.Mode,
		AppointmentType: models.ApptType(input.AppointmentType),
		Notes:           input.Notes,
	}
	scheduleErr := utils.ScheduleAppointment(config.DB, input.DoctorID, patientID, input.AppointmentDate, input.StartTime, input.EndTime, nil)

	if scheduleErr != nil {
		utils.Log.Warnf("CreateAppointment: Scheduling error - %v", scheduleErr)
		c.JSON(400, gin.H{"error": scheduleErr.Error()})
		return
	}

	err = metrics.DbMetrics(config.DB, "insert_appointment", func(db *gorm.DB) error { return db.Create(&appointment).Error })
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
	err := metrics.DbMetrics(config.DB, "get_all_appointments", func(db *gorm.DB) error {
		return db.WithContext(c.Request.Context()).
			Preload("Patient").
			Preload("Doctor").
			Limit(limit).
			Offset((page - 1) * limit).
			Find(&appointments).Error
	})
	if err != nil {
		utils.Log.Errorf("GetAppointments: Database error - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve appointments - " + err.Error()})
		return
	}

	utils.Log.Infof("GetAppointments: Retrieved %d appointments for page %d with limit %d", len(appointments), page, limit)
	c.JSON(200, gin.H{"appointments": appointments, "page": page, "limit": limit, "total": len(appointments)})
}

func GetAppointmentByID(c *gin.Context) {

	appointmentID := c.Param("id")
	user, _ := utils.GetCurrentUser(c)
	var appointment models.Appointment

	// cache
	cachekey := fmt.Sprintf("cache:appointment:%s", appointmentID)
	val, err := config.Rdb.Get(config.Ctx, cachekey).Result()
	switch err {
	case nil:
		var appointment models.Appointment
		metrics.CacheHits.WithLabelValues("appointment_by_id").Inc()
		if jsonErr := json.Unmarshal([]byte(val), &appointment); jsonErr == nil {
			c.JSON(http.StatusOK, gin.H{"appointment": appointment})
			return
		}
	case redis.Nil:
		// key does not exist
		metrics.CacheMisses.WithLabelValues("appointment_by_id").Inc()
	default:
		// Redis error (network, timeout, etc.)
		metrics.CacheMisses.WithLabelValues("appointment_by_id").Inc()
		utils.Log.Warnf("GetAppointmentByID: Redis error - %v", err)
	}

	err = metrics.DbMetrics(config.DB, "get_appointment_by_appt_id", func(db *gorm.DB) error {
		return db.WithContext(c.Request.Context()).Preload("Patient").Preload("Doctor").Where("id = ?", appointmentID).First(&appointment).Error
	})
	if err != nil {
		utils.Log.Errorf("GetAppointmentByID: Appointment not found - %v", err)
		c.JSON(404, gin.H{"error": "Appointment not found - " + err.Error()})
		return
	}
	if models.Role(user.Role) != models.RoleAdmin && (appointment.PatientID != user.UserID || appointment.DoctorID != user.UserID) {
		utils.Log.Warnf("GetAppointmentByID: You are not authorised to access this appointment")
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}
	//store in cache
	data, _ := json.Marshal(appointment)
	config.Rdb.Set(config.Ctx, cachekey, data, 5*time.Minute)

	c.JSON(200, gin.H{"appointment": appointment})
}
func GetAppointmentByDoctorID(c *gin.Context) {
	doctorID := c.Param("id")
	user, err := utils.GetCurrentUser(c)
	if err != nil {
		utils.Log.Warnf("GetAppointmentByDoctorID: Unauthorized access - %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if models.Role(user.Role) == models.RoleDoctor && user.UserID.String() != doctorID {
		utils.Log.Warnf("GetAppointmentByDoctorID: Doctor %s attempted to access data of doctor %s", user.UserID, doctorID)
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

	cachekey := fmt.Sprintf("cache:doctor_appointments:doctor:%s:limit:%d:offset:%d", doctorID, limit, offset)
	val, err := config.Rdb.Get(config.Ctx, cachekey).Result()
	switch err {
	case nil:
		if jsonErr := json.Unmarshal([]byte(val), &appointments); jsonErr == nil {
			c.JSON(http.StatusOK, gin.H{"appointments": appointments})
			return
		}
	case redis.Nil:
		metrics.CacheMisses.WithLabelValues("get_appointments_by_doctors").Inc()

	default:
		metrics.CacheMisses.WithLabelValues("get_appointments_by_doctors").Inc()
		utils.Log.Warnf("GetAppointmentByID: Redis error - %v", err)
	}

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
	data, _ := json.Marshal(appointments)
	config.Rdb.Set(config.Ctx, cachekey, data, 5*time.Minute).Result()

	c.JSON(http.StatusOK, gin.H{"appointments": appointments})
}

func GetAppointmentByPatientID(c *gin.Context) {

	user, err := utils.GetCurrentUser(c)
	if err != nil {
		utils.Log.Warnf("GetAppointmentByPatientID: Unauthorized access attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	patientID := c.Param("id")

	if models.Role(user.Role) == models.RolePatient && user.UserID.String() != patientID {
		utils.Log.Warnf("GetAppointmentByPatientID: Access denied for patient %s to data of %s", user.UserID, patientID)
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

	cachekey := fmt.Sprintf("cache:patient_appointment:patientID %s:limit: %d:offset:%d", patientID, limit, offset)
	val, err := config.Rdb.Get(config.Ctx, cachekey).Result()
	switch err {
	case nil:
		metrics.CacheHits.WithLabelValues("get_appointments_by_patient").Inc()
		if jsonErr := json.Unmarshal([]byte(val), &appointments); jsonErr == nil {
			c.JSON(http.StatusOK, gin.H{"appointments:": appointments})
			return
		}
	case redis.Nil:
		metrics.CacheMisses.WithLabelValues("get_appointments_by_patient").Inc()
	default:
		metrics.CacheMisses.WithLabelValues("get_appointments_by_patient").Inc()
		utils.Log.Warnf("GetAppointmentByID: Redis error - %v", err)
	}
	db := config.DB.WithContext(c.Request.Context()).Preload("Patient").
		Where("patient_id = ?", patientID).
		Order("appointment_date desc").
		Limit(limit).
		Offset(offset)

	if models.Role(user.Role) == models.RoleDoctor {
		db = db.Preload("Doctor").Where("doctor_id = ?", user.UserID)
	}

	err = metrics.DbMetrics(db, "get_appointments_by_patient", func(db *gorm.DB) error { return db.Find(&appointments).Error })
	if err != nil {
		utils.Log.Errorf("GetAppointmentByPatientID: Failed to fetch appointments - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch appointments"})
		return
	}
	data, _ := json.Marshal(appointments)
	config.Rdb.Set(config.Ctx, cachekey, data, 5*time.Minute)

	c.JSON(http.StatusOK, gin.H{"appointments": appointments})
}

func GetAvailableSlots(c *gin.Context) {
	doctorIDParam := c.Query("doctorId")
	dateParam := c.Query("date")

	if doctorIDParam == "" || dateParam == "" {
		utils.Log.Warnf("GetAvailableSlots: Missing required query parameters")
		c.JSON(http.StatusBadRequest, gin.H{"error": "doctorId and date query parameters are required"})
		return
	}

	doctorID, err := uuid.Parse(doctorIDParam)
	if err != nil {
		utils.Log.Warnf("GetAvailableSlots: Invalid doctorId format - %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid doctorId format"})
		return
	}

	appointmentDate, err := time.Parse("2006-01-02", dateParam)

	if err != nil {
		utils.Log.Warnf("GetAvailableSlots: Invalid date format - %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}
	cacheKey := fmt.Sprintf("cache:doctorSchedule:%s:%s", doctorID.String(), appointmentDate.Format("2006-01-02"))
	val, err := config.Rdb.Get(config.Ctx, cacheKey).Result()

	switch err {
	case nil:
		var slots []string
		if jsonErr := json.Unmarshal([]byte(val), &slots); jsonErr == nil {
			metrics.CacheHits.WithLabelValues("get_available_slots").Inc()
			c.JSON(http.StatusOK, gin.H{"availableSlots": slots})
			return
		}
	case redis.Nil:
		metrics.CacheMisses.WithLabelValues("get_available_slots").Inc()
	default:
		metrics.CacheMisses.WithLabelValues("get_available_slots").Inc()
		utils.Log.Warnf("GetAvailableSlots: Redis error - %v", err)
	}

	slots, err := utils.GetAvailableSlots(config.DB, doctorID, appointmentDate)
	if err != nil {
		utils.Log.Errorf("GetAvailableSlots: Failed to retrieve slots - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve available slots"})
		return
	}

	data, _ := json.Marshal(slots)
	config.Rdb.Set(config.Ctx, cacheKey, data, 5*time.Minute)

	c.JSON(http.StatusOK, gin.H{"availableSlots": slots})
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

	// Role-based checks
	switch models.Role(user.Role) {
	case models.RolePatient:
		if user.UserID != appt.PatientID {
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to update this appointment"})
			return
		}
		if appt.Status == "ACCEPTED" {
			c.JSON(http.StatusForbidden, gin.H{"error": "You cannot update an accepted appointment"})
			return
		}

	case models.RoleDoctor:
		if user.UserID != appt.DoctorID {
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to update this appointment"})
			return
		}

	case models.RoleAdmin:
		// Admin can always update
	default:
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid role"})
		return
	}

	var input AppointmentUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	if input.Location != "" && user.Role == string(models.RolePatient) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only doctor or admin can update location"})
		return
	}

	if input.StartTime != "" && input.EndTime != "" {
		if err := utils.ScheduleAppointment(
			config.DB,
			appt.DoctorID,
			appt.PatientID,
			appt.AppointmentDate,
			input.StartTime,
			input.EndTime,
			&appt.ID); err != nil {
			utils.Log.Warnf("UpdateAppointment: Conflict in scheduling the appointment")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	if input.StartTime != "" {
		appt.StartTime = input.StartTime
	}
	if input.EndTime != "" {
		appt.EndTime = input.EndTime
	}
	if input.Mode != "" {
		appt.Mode = input.Mode
	}
	if input.Notes != "" {
		appt.Notes = input.Notes
	}
	if input.Location != "" {
		appt.Location = input.Location
	}

	err = metrics.DbMetrics(config.DB, "update_appointment", func(db *gorm.DB) error { return db.Save(&appt).Error })
	if err != nil {
		utils.Log.Errorf("UpdateAppointment: Failed to update appointment - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update appointment"})
		return
	}

	appointmentCache := cache.NewAppointmentCache(config.Rdb, config.Ctx)

	appointmentCache.Invalidate(appointmentID, appt.DoctorID.String(), appt.PatientID.String(), appt.AppointmentDate.Format("2006-01-02"))

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
	err := metrics.DbMetrics(config.DB, "delete_appointment", func(db *gorm.DB) error { return db.WithContext(c.Request.Context()).Delete(&appointment).Error })
	if err != nil {
		utils.Log.Errorf("DeleteAppointment: Failed to delete appointment - %v", err)
		c.JSON(500, gin.H{"error": "Failed to delete appointment - " + err.Error()})
		return
	}
	appointmentCache := cache.NewAppointmentCache(config.Rdb, config.Ctx)
	appointmentCache.Invalidate(appointmentId, appointment.DoctorID.String(), appointment.PatientID.String(), appointment.AppointmentDate.Format("2006-01-02"))
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

	if models.Role(user.Role) != models.RoleAdmin && (user.UserID != appointment.PatientID || user.UserID != appointment.DoctorID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only reschedule your own appointment"})
		return
	}

	if appointment.Status != "PENDING" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only PENDING appointments can be rescheduled"})
		return
	}

	var input struct {
		Date      time.Time `json:"date" binding:"required"`
		StartTime string    `json:"start_time" binding:"required"`
		EndTime   string    `json:"end_time" binding:"required"`
		Mode      string    `json:"mode" binding:"required,oneof=Online In-Person"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	appointment.AppointmentDate = input.Date
	appointment.StartTime = input.StartTime
	appointment.EndTime = input.EndTime
	appointment.Mode = input.Mode

	if err := utils.ScheduleAppointment(config.DB, appointment.DoctorID, appointment.PatientID, appointment.AppointmentDate, appointment.StartTime, appointment.EndTime, &appointment.ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := metrics.DbMetrics(config.DB, "Reschedule_appointment", func(db *gorm.DB) error { return db.Save(&appointment).Error })
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reschedule"})
		return
	}
	appointmentCache := cache.NewAppointmentCache(config.Rdb, config.Ctx)
	appointmentCache.Invalidate(appointmentID, appointment.DoctorID.String(), appointment.PatientID.String(), appointment.AppointmentDate.Format("2006-01-02"))

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

	allowed := models.Role(user.Role) == models.RoleAdmin || user.UserID == appointment.PatientID
	if !allowed {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		return
	}

	appointment.Status = "CANCELLED"
	err := metrics.DbMetrics(config.DB, "cancel_appointment", func(db *gorm.DB) error { return db.Save(&appointment).Error })
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Appointment cancelled"})
}
