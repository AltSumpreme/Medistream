package appointments

import (
	"net/http"
	"time"

	"github.com/AltSumpreme/Medistream.git/queue"
	"github.com/AltSumpreme/Medistream.git/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AppointmentInput struct {
	UserID          uuid.UUID `json:"userId" binding:"required"`
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

func HandleUserCreateAppointment(c *gin.Context, q *queue.RedisQueueConfig) {
	_, err := utils.GetCurrentUser(c)
	if err != nil {
		utils.Log.Warnf("CreateAppointment: Failed to get current user - %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input AppointmentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Log.Warnf("CreateAppointment: Invalid input - %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}
	job := queue.JobPayload{
		Type: queue.JobTypeCreateAppointment,
		Data: input,
	}
	if err := q.Enqueue(c.Request.Context(), "appointment_queue", job); err != nil {
		utils.Log.Errorf("SignUp: Failed to enqueue job - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process signup"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "user created appointment successfully"})
}
