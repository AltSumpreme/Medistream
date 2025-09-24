package medicalrecords

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

type MedicalRecordInput struct {
	PatientID      uuid.UUID    `json:"patient_id" binding:"required"`
	DoctorID       uuid.UUID    `json:"doctor_id" binding:"required"`
	Diagnosis      string       `json:"diagnosis" binding:"required"`
	Notes          string       `json:"notes"`
	VitalsToCreate []VitalInput `json:"vitals_to_create"`
	VitalIDsToLink []uuid.UUID  `json:"vital_ids_to_link"`
}
type VitalInput struct {
	Type       string    `json:"type" binding:"required"`
	Value      string    `json:"value" binding:"required,oneof=BLOOD_PRESSURE HEART_RATE WEIGHT BMI TEMPERATURE RESPIRATORY_RATE OXYGEN_SATURATION"`
	Status     string    `json:"status" binding:"required"`
	RecordedAt time.Time `json:"recorded_at" binding:"required"`
}

func CreateMedicalRecord(c *gin.Context) {
	var input MedicalRecordInput

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Log.Warnf("Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tx := config.DB.WithContext(c).Begin()
	if tx.Error != nil {
		utils.Log.Errorf("Failed to start transaction: %v", tx.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	committed := false
	defer func() {
		if !committed {
			utils.Log.Warnf("Rolling back transaction due to error")
			if err := tx.Rollback().Error; err != nil {
				utils.Log.Errorf("Failed to rollback transaction: %v", err)
			}
			return
		}
	}()

	record := models.MedicalRecord{
		PatientID: input.PatientID,
		DoctorID:  input.DoctorID,
		Diagnosis: input.Diagnosis,
		Notes:     input.Notes,
	}

	err := metrics.DbMetrics(config.DB, "create_medical_records", func(d *gorm.DB) error { return tx.Create(&record).Error })
	if err != nil {
		utils.Log.Errorf("Failed to create medical record: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create medical record"})
		return
	}

	// Inserting any vitals (if any)
	if len(input.VitalsToCreate) > 0 {
		vitalsToInsert := make([]models.Vital, 0, len(input.VitalsToCreate))
		for _, v := range input.VitalsToCreate {
			vitalsToInsert = append(vitalsToInsert, models.Vital{
				ID:              uuid.New(),
				PatientID:       input.PatientID,
				Type:            models.VitalType(v.Type),
				Value:           v.Value,
				Status:          v.Status,
				RecordedAt:      v.RecordedAt,
				MedicalRecordID: &record.ID,
			})
		}

		if err := metrics.DbMetrics(config.DB, "create_vitals", func(d *gorm.DB) error {
			return tx.Create(&vitalsToInsert).Error
		}); err != nil {
			utils.Log.Errorf("Failed to create vitals: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create vitals"})
			return
		}
	}

	// Linking existing vitals (if any)
	if len(input.VitalIDsToLink) > 0 {
		if err := metrics.DbMetrics(config.DB, "link_vitals", func(d *gorm.DB) error {
			return tx.Model(&models.Vital{}).
				Where("id IN ?", input.VitalIDsToLink).
				Update("medical_record_id", record.ID).Error
		}); err != nil {
			utils.Log.Errorf("Failed to associate existing vitals: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to associate vitals"})
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		utils.Log.Errorf("Transaction commit failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create medical record"})
		return
	}
	committed = true

	c.JSON(http.StatusCreated, gin.H{
		"message":   "Medical record created successfully",
		"record_id": record.ID,
	})
}

func GetMedicalRecordByID(c *gin.Context) {
	recordID := c.Param("id")
	if recordID == "" {
		utils.Log.Warnf("Medical record ID is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Medical record ID is required"})
		return
	}
	var record models.MedicalRecord
	cachekey := fmt.Sprintf("medical_record_%s", recordID)
	val, err := config.Rdb.Get(c, cachekey).Result()
	switch err {
	case nil:

		metrics.CacheHits.WithLabelValues("medical_by_id").Inc()
		if jsonErr := json.Unmarshal([]byte(val), &record); jsonErr == nil {
			c.JSON(http.StatusOK, gin.H{"record": record})
			return
		}
	case redis.Nil:
		metrics.CacheMisses.WithLabelValues("medical_by_id").Inc()
	default:
		metrics.CacheMisses.WithLabelValues("medical_by_id").Inc()
		utils.Log.Warnf("Redis GET error: %v", err)
	}

	err = metrics.DbMetrics(config.DB, "get_medical_record", func(d *gorm.DB) error {
		return d.WithContext(c).
			Preload("Vitals").
			Preload("Doctor").
			Preload("Patient").
			First(&record, "id = ?", recordID).Error
	})
	if err != nil {
		utils.Log.Errorf("Failed to retrieve medical record: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Medical record not found"})
		return
	}
	data, _ := json.Marshal(record)
	config.Rdb.Set(config.Ctx, cachekey, data, 5*time.Minute).Err()

	c.JSON(http.StatusOK, record)
}

func GetRecordsByPatientID(c *gin.Context) {
	patientID := c.Param("id")
	if patientID == "" {
		utils.Log.Warnf("Patient ID is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Patient ID is required"})
		return
	}

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
	offset := (page - 1) * limit

	cachekey := fmt.Sprintf("medical_records_patient_%s_page_%d_limit_%d", patientID, page, limit)
	val, err := config.Rdb.Get(config.Ctx, cachekey).Result()
	switch err {
	case nil:
		var records []models.MedicalRecord
		metrics.CacheHits.WithLabelValues("medical_by_patient").Inc()
		if jsonErr := json.Unmarshal([]byte(val), &records); jsonErr == nil {
			c.JSON(http.StatusOK, gin.H{"records": records})
			return
		}
	case redis.Nil:
		metrics.CacheMisses.WithLabelValues("medical_by_patient").Inc()
	default:
		metrics.CacheMisses.WithLabelValues("medical_by_patient").Inc()
		utils.Log.Warnf("Redis GET error: %v", err)
	}

	var records []models.MedicalRecord
	err = metrics.DbMetrics(config.DB, "get_records_by_patient", func(d *gorm.DB) error {
		return d.WithContext(c).
			Preload("Vitals").
			Preload("Doctor").
			Where("patient_id = ?", patientID).
			Limit(limit).Offset(offset).
			Find(&records).Error
	})
	if err != nil {
		utils.Log.Errorf("Failed to fetch records for patient %s: %v", patientID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch records"})
		return
	}
	data, _ := json.Marshal(records)
	config.Rdb.Set(config.Ctx, cachekey, data, 5*time.Minute).Err()

	c.JSON(http.StatusOK, records)
}

func UpdateMedicalRecord(c *gin.Context, medicalrecordCache *cache.Cache) {
	recordID := c.Param("id")
	if recordID == "" {
		utils.Log.Warnf("Medical record ID is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Medical record ID is required"})
		return
	}

	var input struct {
		Diagnosis      string      `json:"diagnosis"`
		Notes          string      `json:"notes"`
		VitalIDsToLink []uuid.UUID `json:"vital_ids_to_link"` // Only linking
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Log.Warnf("Invalid input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	tx := config.DB.WithContext(c).Begin()
	if tx.Error != nil {
		utils.Log.Errorf("Failed to start transaction: %v", tx.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	// Update the medical record fields
	if err := metrics.DbMetrics(config.DB, "update_medical_record", func(d *gorm.DB) error {
		return tx.Model(&models.MedicalRecord{}).
			Where("id = ?", recordID).
			Updates(map[string]interface{}{
				"diagnosis": input.Diagnosis,
				"notes":     input.Notes,
			}).Error
	}); err != nil {
		utils.Log.Errorf("Failed to update record %s: %v", recordID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update record"})
		return
	}

	// Link vitals if provided
	if len(input.VitalIDsToLink) > 0 {
		if err := metrics.DbMetrics(config.DB, "link_vitals_update", func(d *gorm.DB) error {
			return tx.Model(&models.Vital{}).
				Where("id IN ?", input.VitalIDsToLink).
				Update("medical_record_id", recordID).Error
		}); err != nil {
			utils.Log.Errorf("Failed to link vitals for record %s: %v", recordID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to link vitals"})
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		utils.Log.Errorf("Transaction commit failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update record"})
		return
	}
	var rec models.MedicalRecord
	if err := config.DB.First(&rec, "id = ?", recordID).Error; err == nil {
		medicalrecordCache.MedicalRecordInvalidate(recordID, rec.PatientID.String())
	}
	c.JSON(http.StatusOK, gin.H{"message": "Medical record updated successfully"})
}

func SoftDeleteMedicalRecord(c *gin.Context, medicalrecordCache *cache.Cache) {
	recordID := c.Param("id")
	if recordID == "" {
		utils.Log.Warnf("Medical record ID is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record ID required"})
		return
	}

	err := metrics.DbMetrics(config.DB, "soft_delete_medical_record", func(db *gorm.DB) error {
		return db.WithContext(c).Delete(&models.MedicalRecord{}, "id = ?", recordID).Error
	})
	if err != nil {
		utils.Log.Errorf("Failed to soft delete medical record %s: %v", recordID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to soft delete"})
		return
	}
	var rec models.MedicalRecord
	if err := config.DB.First(&rec, "id = ?", recordID).Error; err == nil {
		medicalrecordCache.MedicalRecordInvalidate(recordID, rec.PatientID.String())
	}

	c.JSON(http.StatusOK, gin.H{"message": "Medical record  deleted"})
}

func HardDeleteMedicalRecord(c *gin.Context, medicalrecordCache *cache.Cache) {
	recordID := c.Param("id")
	if recordID == "" {
		utils.Log.Warnf("Record ID is required for hard delete")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record ID required"})
		return
	}

	err := metrics.DbMetrics(config.DB, "hard_delete_medical_record", func(db *gorm.DB) error {
		return db.WithContext(c).Unscoped().Delete(&models.MedicalRecord{}, "id = ?", recordID).Error
	})
	if err != nil {
		utils.Log.Errorf("Failed to hard delete medical record %s: %v", recordID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hard delete"})
		return
	}
	var rec models.MedicalRecord
	if err := config.DB.First(&rec, "id = ?", recordID).Error; err == nil {
		medicalrecordCache.MedicalRecordInvalidate(recordID, rec.PatientID.String())
	}

	c.JSON(http.StatusOK, gin.H{"message": "Medical record  deleted"})
}
