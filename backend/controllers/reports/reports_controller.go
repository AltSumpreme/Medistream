package reports

import (
	"net/http"
	"strconv"

	"github.com/AltSumpreme/Medistream.git/config"
	"github.com/AltSumpreme/Medistream.git/models"
	"github.com/AltSumpreme/Medistream.git/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ReportInput struct {
	Title           string    `json:"title" binding:"required"`
	Description     string    `json:"description" binding:"required"`
	FileURL         string    `json:"file_url" binding:"required"`
	PatientID       uuid.UUID `json:"patient_id" binding:"required"`
	MedicalRecordID uuid.UUID `json:"medical_record_id" binding:"required"`
}

func CreateReport(c *gin.Context) {
	var input ReportInput

	user, err := utils.GetCurrentUser(c)
	if err != nil {
		utils.Log.Warnf("Failed to get current user: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Log.Warnf("Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	report := models.Report{
		Title:           input.Title,
		Description:     input.Description,
		FileURL:         input.FileURL,
		DoctorID:        user.ID,
		PatientID:       input.PatientID,
		MedicalRecordID: &input.MedicalRecordID,
	}
	if err := config.DB.Create(&report).Error; err != nil {
		utils.Log.Errorf("Failed to create report: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create report"})
		return
	}

	var record models.MedicalRecord
	if err := config.DB.Where("id = ? AND doctor_id = ?", input.MedicalRecordID).First(&record).Error; err != nil {
		utils.Log.Warnf("Medical record not found or unauthorized access: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "You are not authorized to create a report for this medical record"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Report created successfully"})

}
func GetReportByPatientID(c *gin.Context) {
	var reports []models.Report
	patientID := c.Param("patient_id")
	if patientID == "" {
		utils.Log.Warnf("Patient ID is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Patient ID is required"})
		return
	}

	page := 1
	limit := 10
	MaxLimit := 20
	if l := c.Query("limit"); l != "" {
		if n, err := strconv.Atoi(l); err == nil && n > 0 && n <= MaxLimit {
			limit = n
		} else {
			limit = MaxLimit
		}
	}
	if p := c.Query("page"); p != "" {
		if n, err := strconv.Atoi(p); err == nil {
			page = n
		}
	}
	offset := (page - 1) * limit

	if err := config.DB.Where("patient_id = ?", patientID).Offset(offset).Limit(limit).Find(&reports).Error; err != nil {
		utils.Log.Warnf("Failed to get reports: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get reports"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reports": reports})
}

func GetReportByID(c *gin.Context) {
	reportID := c.Param("id")
	if reportID == "" {
		utils.Log.Warnf("Report ID is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Report ID is required"})
		return
	}

	var report models.Report
	if err := config.DB.Preload("Doctor").Preload("Patient").Preload("MedicalRecord").First(&report, "id = ?", reportID).Error; err != nil {
		utils.Log.Errorf("Failed to retrieve report: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Report not found"})
		return
	}

	c.JSON(http.StatusOK, report)
}

func UpdateReportByID(c *gin.Context) {
	reportID := c.Param("id")
	if reportID == "" {
		utils.Log.Warnf("Report ID is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Report ID is required"})
		return
	}
	var input ReportInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Log.Warnf("Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	var report models.Report
	if err := config.DB.First(&report, "id = ?", reportID).Error; err != nil {
		utils.Log.Errorf("Failed to retrieve report: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Report not found"})
		return
	}
	report.Title = input.Title
	report.Description = input.Description
	report.FileURL = input.FileURL
	report.PatientID = input.PatientID
	report.MedicalRecordID = &input.MedicalRecordID
	if err := config.DB.Save(&report).Error; err != nil {
		utils.Log.Errorf("Failed to update report: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update report"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Report updated successfully"})
}

func DeleteReportByID(c *gin.Context) {
	reportID := c.Param("id")
	if reportID == "" {
		utils.Log.Warnf("Report ID is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Report ID is required"})
		return
	}

	var report models.Report
	if err := config.DB.First(&report, "id = ?", reportID).Error; err != nil {
		utils.Log.Errorf("Failed to retrieve report: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Report not found"})
		return
	}

	if err := config.DB.Delete(&report).Error; err != nil {
		utils.Log.Errorf("Failed to delete report: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete report"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Report deleted successfully"})
}
