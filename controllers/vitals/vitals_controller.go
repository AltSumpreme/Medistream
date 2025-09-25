package vitals

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

type VitalInput struct {
	PatientID  uuid.UUID        `json:"patient_id" binding:"required"`
	Type       models.VitalType `json:"type" binding:"required,oneof=BLOOD_PRESSURE HEART_RATE WEIGHT BMI TEMPERATURE RESPIRATORY_RATE OXYGEN_SATURATION"`
	Value      string           `json:"value" binding:"required"`
	Status     string           `json:"status" binding:"required"`
	RecordedAt time.Time        `json:"recorded_at" binding:"required"`
}
type VitalUpdateInput struct {
	Type       models.VitalType `json:"type" binding:"required,oneof=BLOOD_PRESSURE HEART_RATE WEIGHT BMI TEMPERATURE RESPIRATORY_RATE OXYGEN_SATURATION"`
	Value      string           `json:"value" binding:"required"`
	Status     string           `json:"status" binding:"required"`
	RecordedAt time.Time        `json:"recorded_at" binding:"required"`
}

func CreateVital(c *gin.Context) {
	var input VitalInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Log.Warnf("Invalid input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	vital := models.Vital{
		ID:         uuid.New(),
		PatientID:  input.PatientID,
		Type:       input.Type,
		Value:      input.Value,
		Status:     input.Status,
		RecordedAt: input.RecordedAt,
	}

	err := metrics.DbMetrics(config.DB, "create_vitals", func(db *gorm.DB) error {
		return db.WithContext(c).Create(&vital).Error
	})
	if err != nil {
		utils.Log.Errorf("Failed to create vital: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create vital"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Vital created successfully", "vital_id": vital.ID})
}

func GetVitalsByPatientID(c *gin.Context) {
	patientID := c.Param("id")
	user, _ := utils.GetCurrentUser(c)
	if patientID == "" {
		utils.Log.Warnf("GetVitalsByPatientID: PatientID required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Patient ID is required"})
		return
	}
	switch models.Role(user.Role) {
	case models.RolePatient:
		if user.ID != patientID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized access to another patient's data"})
			return
		}
	case models.RoleDoctor:
		hasAccess, err := utils.DoctorHasAccessToPatient(user.UserID, uuid.MustParse(patientID), c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Access check failed"})
			return
		}
		if !hasAccess {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not assigned to this patient"})
			return
		}

	}

	limit := 10
	page := 1
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil {
			page = parsed
		}
	}
	offset := (page - 1) * limit

	cachekey := fmt.Sprintf("cache:vitals:patient:%s:limit:%d:offset:%d", patientID, limit, offset)

	val, err := config.Rdb.Get(c, cachekey).Result()
	switch err {
	case nil:
		var vitals []models.Vital
		if err := json.Unmarshal([]byte(val), &vitals); err == nil {
			metrics.CacheHits.WithLabelValues("vitals").Inc()
			c.JSON(http.StatusOK, vitals)
			return
		}
	case redis.Nil:
		metrics.CacheMisses.WithLabelValues("vitals").Inc()
	default:
		metrics.CacheMisses.WithLabelValues("vitals").Inc()
		utils.Log.Warnf("Redis GET error: %v", err)
	}

	var vitals []models.Vital
	err = metrics.DbMetrics(config.DB, "get_vitals", func(db *gorm.DB) error {
		return db.WithContext(c).Where("patient_id = ?", patientID).
			Limit(limit).
			Offset(offset).
			Order("recorded_at DESC").
			Find(&vitals).Error
	})
	if err != nil {
		utils.Log.Errorf("Failed to fetch vitals for patient %s: %v", patientID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch vitals"})
		return
	}
	data, _ := json.Marshal(vitals)
	err = config.Rdb.Set(c, cachekey, data, 5*time.Minute).Err()
	if err != nil {
		utils.Log.Warnf("Redis SET error: %v", err)
	}

	c.JSON(http.StatusOK, vitals)
}

func GetVitalByID(c *gin.Context) {
	vitalID := c.Param("id")
	user, _ := utils.GetCurrentUser(c)

	if vitalID == "" {
		utils.Log.Warnf("Vital ID is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Vital ID is required"})
		return
	}

	var vital models.Vital
	err := metrics.DbMetrics(config.DB, "get_vital", func(db *gorm.DB) error {
		return db.WithContext(c).
			Preload("Patient").
			Preload("MedicalRecord").
			First(&vital, "id = ?", vitalID).Error
	})
	if err != nil {
		utils.Log.Warnf("Vital not found: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Vital not found"})
		return
	}

	switch models.Role(user.Role) {
	case models.RolePatient:
		if vital.PatientID != user.UserID {
			utils.Log.Warnf("GetVitalByID: Unauthorized access by patient to vital")
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to access this vital record"})
			return
		}

	case models.RoleDoctor:
		if vital.MedicalRecordID == nil {
			utils.Log.Warnf("GetVitalByID: Vital has no medical record linked")
			c.JSON(http.StatusForbidden, gin.H{"error": "No medical record linked to this vital"})
			return
		}

		var record models.MedicalRecord
		err := metrics.DbMetrics(config.DB, "get_medical_record", func(db *gorm.DB) error {
			return db.WithContext(c).
				Select("id").
				First(&record, "id = ? AND doctor_id = ?", *vital.MedicalRecordID, user.ID).Error
		})
		if err != nil {
			utils.Log.Warnf("Doctor %s tried to access unassigned vital %s", user.ID, vital.ID)
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to access this vital record"})
			return
		}
	}

	c.JSON(http.StatusOK, vital)
}

func UpdateVital(c *gin.Context, vitalsCache *cache.Cache) {
	vitalID := c.Param("id")
	if vitalID == "" {
		utils.Log.Warnf("Vital ID is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Vital ID is required"})
		return
	}

	var input struct {
		Value      string    `json:"value"`
		Status     string    `json:"status"`
		RecordedAt time.Time `json:"recorded_at"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	err := metrics.DbMetrics(config.DB, "update_vital", func(db *gorm.DB) error {
		return db.WithContext(c).
			Model(&models.Vital{}).
			Where("id = ?", vitalID).
			Updates(map[string]interface{}{
				"value":       input.Value,
				"status":      input.Status,
				"recorded_at": input.RecordedAt,
			}).Error
	})
	if err != nil {
		utils.Log.Errorf("Failed to update vital %s: %v", vitalID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update vital"})
		return
	}
	var vital models.Vital
	err = metrics.DbMetrics(config.DB, "get_vital_patient", func(db *gorm.DB) error {
		return db.WithContext(c).Select("patient_id").First(&vital, "id = ?", vitalID).Error
	})
	if err != nil {
		utils.Log.Errorf("Failed to fetch vital %s for cache invalidation: %v", vitalID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update vital"})
		return
	}
	vitalsCache.VitalsInvalidate(vital.PatientID.String())

	c.JSON(http.StatusOK, gin.H{"message": "Vital updated successfully"})
}

func DeleteVital(c *gin.Context, vitalsCache *cache.Cache) {
	vitalID := c.Param("id")
	if vitalID == "" {
		utils.Log.Warnf("Vital ID is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Vital ID is required"})
		return
	}

	err := metrics.DbMetrics(config.DB, "delete_vital", func(db *gorm.DB) error {
		return db.WithContext(c).Delete(&models.Vital{}, "id = ?", vitalID).Error
	})
	if err != nil {
		utils.Log.Errorf("Failed to soft delete vital %s: %v", vitalID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete vital"})
		return
	}

	var vital models.Vital
	err = metrics.DbMetrics(config.DB, "get_vital_patient", func(db *gorm.DB) error {
		return db.WithContext(c).Select("patient_id").First(&vital, "id = ?", vitalID).Error
	})
	if err != nil {
		utils.Log.Errorf("Failed to fetch vital %s for cache invalidation: %v", vitalID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete vital"})
		return
	}
	vitalsCache.VitalsInvalidate(vital.PatientID.String())
	c.JSON(http.StatusOK, gin.H{"message": "Vital soft deleted"})
}
