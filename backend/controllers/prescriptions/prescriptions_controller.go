package prescriptions

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

type PrescriptionInput struct {
	PatientID       uuid.UUID `json:"patient_id" binding:"required"`
	DoctorID        uuid.UUID `json:"doctor_id" binding:"required"`
	MedicalRecordID uuid.UUID `json:"medical_record_id"`
	Medication      string    `json:"medication" binding:"required"`
	Dosage          string    `json:"dosage" binding:"required"`
	Instructions    string    `json:"instructions"`
	IssuedAt        time.Time `json:"issued_at" binding:"required"`
}

func CreatePrescription(c *gin.Context) {
	var input PrescriptionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Log.Warnf("Invalid prescription input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	prescription := models.Prescription{
		ID:              uuid.New(),
		PatientID:       input.PatientID,
		DoctorID:        input.DoctorID,
		MedicalRecordID: &input.MedicalRecordID,
		Medication:      input.Medication,
		Dosage:          input.Dosage,
		Instructions:    input.Instructions,
		IssuedAt:        input.IssuedAt,
	}

	if err := config.DB.WithContext(c).Create(&prescription).Error; err != nil {
		utils.Log.Errorf("Failed to create prescription: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create prescription"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Prescription created", "prescription_id": prescription.ID})
}

func GetPrescriptionsByPatientID(c *gin.Context) {
	patientID := c.Param("id")
	user, _ := utils.GetCurrentUser(c)

	if patientID == "" {
		utils.Log.Warnf("GetPrescriptionsByPatientID: Patient ID is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Patient ID is required"})
		return
	}

	switch user.Role {
	case models.RolePatient:
		if user.ID.String() != patientID {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to access this patient's prescriptions"})
			return
		}
	case models.RoleDoctor:
		hasAccess, err := utils.DoctorHasAccessToPatient(user.ID, uuid.MustParse(patientID), c)
		if err != nil {
			utils.Log.Errorf("Access check failed for doctor %s and patient %s: %v", user.ID, patientID, err)
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
		if n, err := strconv.Atoi(l); err == nil {
			limit = n
		}
	}
	if p := c.Query("page"); p != "" {
		if n, err := strconv.Atoi(p); err == nil {
			page = n
		}
	}
	offset := (page - 1) * limit

	var prescriptions []models.Prescription
	if err := config.DB.WithContext(c).
		Where("patient_id = ?", patientID).
		Limit(limit).
		Offset(offset).
		Order("issued_at DESC").
		Find(&prescriptions).Error; err != nil {
		utils.Log.Errorf("Failed to fetch prescriptions: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch prescriptions"})
		return
	}

	c.JSON(http.StatusOK, prescriptions)
}

func GetPrescriptionByID(c *gin.Context) {
	prescriptionID := c.Param("id")
	user, _ := utils.GetCurrentUser(c)

	if prescriptionID == "" {
		utils.Log.Warnf("GetPrescriptionByID: Prescription ID is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Prescription ID is required"})
		return
	}

	var prescription models.Prescription
	if err := config.DB.WithContext(c).
		Preload("Patient").
		Preload("MedicalRecord").
		First(&prescription, "id = ?", prescriptionID).Error; err != nil {
		utils.Log.Warnf("Prescription not found: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Prescription not found"})
		return
	}

	switch user.Role {
	case models.RolePatient:
		if prescription.PatientID != user.ID {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to access this prescription"})
			return
		}

	case models.RoleDoctor:
		if prescription.MedicalRecordID == nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "No medical record linked to this prescription"})
			return
		}

		var record models.MedicalRecord
		if err := config.DB.WithContext(c).
			Select("id").
			Where("id = ? AND doctor_id = ?", prescription.MedicalRecordID, user.ID).
			First(&record).Error; err != nil {
			utils.Log.Errorf("GetPrescriptionByID: You are not authorized to access this prescription %v", err)
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to access this prescription"})
			return
		}
	}

	c.JSON(http.StatusOK, prescription)
}

func UpdatePrescription(c *gin.Context) {
	prescriptionID := c.Param("id")
	if prescriptionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Prescription ID is required"})
		return
	}

	var input PrescriptionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	updateData := map[string]interface{}{
		"medication":   input.Medication,
		"dosage":       input.Dosage,
		"instructions": input.Instructions,
		"issued_at":    input.IssuedAt,
	}
	if input.MedicalRecordID != uuid.Nil {
		updateData["medical_record_id"] = input.MedicalRecordID
	}

	if err := config.DB.WithContext(c).
		Model(&models.Prescription{}).
		Where("id = ?", prescriptionID).
		Updates(updateData).Error; err != nil {
		utils.Log.Errorf("Failed to update prescription %s: %v", prescriptionID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update prescription"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Prescription updated successfully"})
}

func DeletePrescription(c *gin.Context) {
	prescriptionID := c.Param("id")
	if prescriptionID == "" {
		utils.Log.Warnf("DeletePrescription: Prescription ID is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Prescription ID is required"})
		return
	}

	if err := config.DB.WithContext(c).
		Delete(&models.Prescription{}, "id = ?", prescriptionID).Error; err != nil {
		utils.Log.Errorf("Failed to delete prescription: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete prescription"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Prescription deleted"})
}
