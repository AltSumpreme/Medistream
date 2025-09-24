package reports

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

type ReportInput struct {
	Title           string    `json:"title" binding:"required"`
	Description     string    `json:"description" binding:"required"`
	FileURL         string    `json:"file_url" binding:"required"`
	PatientID       uuid.UUID `json:"patient_id" binding:"required"`
	DoctorID        uuid.UUID `json:"doctor_id" binding:"required"`
	MedicalRecordID uuid.UUID `json:"medical_record_id" binding:"required"`
}

func CreateReport(c *gin.Context) {
	var input ReportInput

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Log.Warnf("Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	var record models.MedicalRecord
	if err := metrics.DbMetrics(config.DB, "get_medical_record_by_id_vitals", func(db *gorm.DB) error {
		return db.Where("id = ? AND doctor_id = ?", input.MedicalRecordID).First(&record).Error
	}); err != nil {
		utils.Log.Warnf("Medical record not found or unauthorized access: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "You are not authorized to create a report for this medical record"})
		return
	}

	report := models.Report{
		Title:           input.Title,
		Description:     input.Description,
		FileURL:         input.FileURL,
		DoctorID:        input.DoctorID,
		PatientID:       input.PatientID,
		MedicalRecordID: &input.MedicalRecordID,
	}
	if err := metrics.DbMetrics(config.DB, "create_report", func(db *gorm.DB) error {
		return db.Create(&report).Error
	}); err != nil {
		utils.Log.Errorf("Failed to create report: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create report"})
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

	cachekey := fmt.Sprintf("cache:reports:patient:%s:limit:%d:offset:%d", patientID, limit, offset)
	val, err := config.Rdb.Get(c, cachekey).Result()
	switch err {
	case nil:
		var reports []models.Report
		metrics.CacheHits.WithLabelValues("GetReportByPatientID").Inc()
		if err := json.Unmarshal([]byte(val), &reports); err == nil {
			c.JSON(http.StatusOK, reports)
			return
		}
	case redis.Nil:
		metrics.CacheMisses.WithLabelValues("GetReportByPatientID").Inc()
	default:
		metrics.CacheMisses.WithLabelValues("GetReportsByPatientID").Inc()
	}

	if err := metrics.DbMetrics(config.DB, "get_reports_by_patient_id", func(db *gorm.DB) error {
		return db.Where("patient_id = ?", patientID).Offset(offset).Limit(limit).Find(&reports).Error
	}); err != nil {
		utils.Log.Warnf("Failed to get reports: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get reports"})
		return
	}
	data, _ := json.Marshal(reports)
	config.Rdb.Set(c, cachekey, data, 10*time.Minute)

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

func UpdateReportByID(c *gin.Context, reportsCache *cache.Cache) {
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
	if err := metrics.DbMetrics(config.DB, "update_reports_by_id", func(db *gorm.DB) error {
		return db.Save(&report).Error
	}); err != nil {
		utils.Log.Errorf("Failed to update report: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update report"})
		return
	}
	if err := metrics.DbMetrics(config.DB, "updated_reports_by_id", func(db *gorm.DB) error {
		return db.Preload("Patient").First("patient_id = ?", report.PatientID).Error
	}); err != nil {
		utils.Log.Errorf("Failed to preload patient: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to preload patient"})
		return
	}
	reportsCache.ReportInvalidate(report.PatientID.String())
	c.JSON(http.StatusOK, gin.H{"message": "Report updated successfully"})
}

func DeleteReportByID(c *gin.Context, reportsCache *cache.Cache) {
	reportID := c.Param("id")
	if reportID == "" {
		utils.Log.Warnf("Report ID is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Report ID is required"})
		return
	}

	var report models.Report
	err := metrics.DbMetrics(config.DB, "get_report_by_id", func(db *gorm.DB) error {
		return db.First(&report, "id = ?", reportID).Error
	})
	if err != nil {
		utils.Log.Errorf("Failed to retrieve report: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Report not found"})
		return
	}

	err = metrics.DbMetrics(config.DB, "delete_report_by_id", func(db *gorm.DB) error {
		return db.Delete(&report).Error
	})
	if err != nil {
		utils.Log.Errorf("Failed to delete report: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete report"})
		return
	}
	reportsCache.ReportInvalidate(report.PatientID.String())

	c.JSON(http.StatusOK, gin.H{"message": "Report deleted successfully"})
}
