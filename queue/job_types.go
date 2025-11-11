package queue

type JobType string

const (
	JobTypeCreateAppointment JobType = "appointment:create"
	JobTypeWelcomeEmail      JobType = "email:welcome"
	JobOTPEmail              JobType = "email:otp"
	JobTypeResetPassword     JobType = "email:reset_password"
)

type JobPayload struct {
	Type JobType `json:"type"`
	Data any     `json:"data"`
}

type EmailPayload struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}
