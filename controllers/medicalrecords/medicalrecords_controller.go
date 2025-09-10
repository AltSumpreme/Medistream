package medicalrecords

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

	if err := tx.Create(&record).Error; err != nil {
		utils.Log.Errorf("Failed to create medical record: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create medical record"})
		return
	}

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

		if err := tx.Create(&vitalsToInsert).Error; err != nil {
			utils.Log.Errorf("Failed to create vitals: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create vitals"})
			return
		}
	}

	if len(input.VitalIDsToLink) > 0 {
		if err := tx.Model(&models.Vital{}).
			Where("id IN ?", input.VitalIDsToLink).
			Update("medical_record_id", record.ID).Error; err != nil {
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

func GetMedicalRecord(c *gin.Context) {
	recordID := c.Param("id")
	if recordID == "" {
		utils.Log.Warnf("Medical record ID is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Medical record ID is required"})
		return
	}

	var record models.MedicalRecord
	if err := config.DB.WithContext(c).Preload("Vitals").Preload("Doctor").Preload("Patient").First(&record, "id = ?", recordID).Error; err != nil {
		utils.Log.Errorf("Failed to retrieve medical record: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Medical record not found"})
		return
	}

	c.JSON(http.StatusOK, record)
}

func GetRecordsByPatientID(c *gin.Context) {
	patientID := c.Param("id")
	if patientID == "" {
		utils.Log.Warnf("Patient ID is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Patient ID is required"})
		return
	}
	var records []models.MedicalRecord

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

	if err := config.DB.WithContext(c).
		Preload("Vitals").
		Preload("Doctor").
		Where("patient_id = ?", patientID).
		Find(&records).Limit(limit).Offset(offset).Error; err != nil {
		utils.Log.Errorf("Failed to fetch records for patient %s: %v", patientID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch records"})
		return
	}

	c.JSON(http.StatusOK, records)
}

func UpdateMedicalRecord(c *gin.Context) {
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
	if err := tx.Model(&models.MedicalRecord{}).
		Where("id = ?", recordID).
		Updates(map[string]interface{}{
			"diagnosis": input.Diagnosis,
			"notes":     input.Notes,
		}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update record"})
		return
	}

	// Link vitals if provided
	if len(input.VitalIDsToLink) > 0 {
		if err := tx.Model(&models.Vital{}).
			Where("id IN ?", input.VitalIDsToLink).
			Update("medical_record_id", recordID).Error; err != nil {
			tx.Rollback()
			utils.Log.Errorf("Failed to link vitals: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to link vitals"})
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		utils.Log.Errorf("Transaction commit failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Medical record updated successfully"})
}

func SoftDeleteMedicalRecord(c *gin.Context) {
	recordID := c.Param("id")
	if recordID == "" {
		utils.Log.Warnf("Medical record ID is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record ID required"})
		return
	}

	if err := config.DB.WithContext(c).Delete(&models.MedicalRecord{}, "id = ?", recordID).Error; err != nil {
		utils.Log.Errorf("Failed to soft delete medical record %s: %v", recordID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to soft delete"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Medical record soft deleted"})
}

func HardDeleteMedicalRecord(c *gin.Context) {
	recordID := c.Param("id")
	if recordID == "" {
		utils.Log.Warnf("Record ID is required for hard delete")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record ID required"})
		return
	}

	if err := config.DB.WithContext(c).Unscoped().Delete(&models.MedicalRecord{}, "id = ?", recordID).Error; err != nil {
		utils.Log.Errorf("Failed to hard delete medical record %s: %v", recordID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hard delete"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Medical record hard deleted"})
}
