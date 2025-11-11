package queue

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

type WelcomePayload struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type OTPPayload struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

func NewTask(jobType JobType, payload any) (*asynq.Task, error) {
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(string(jobType), b), nil
}

func NewWelcomeEmailTask(email, name string) (*asynq.Task, error) {
	p, err := json.Marshal(WelcomePayload{Email: email, Name: name})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(string(JobTypeWelcomeEmail), p), nil
}

func NewOTPEmailTask(email, otp string) (*asynq.Task, error) {
	p, err := json.Marshal(OTPPayload{Email: email, OTP: otp})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(string(JobOTPEmail), p), nil
}

func ResetEmailTask(email, subject, body string) (*asynq.Task, error) {
	payload := map[string]string{
		"email":   email,
		"subject": subject,
		"body":    body,
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(string(JobTypeResetPassword), b), nil
}
