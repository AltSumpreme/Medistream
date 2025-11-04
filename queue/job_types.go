package queue

type JobType string

const (
	JobTypeCreateAppointment JobType = "create_appointment"
)

type JobPayload struct {
	Type JobType `json:"type"`
	Data any     `json:"data"`
}
