package workers

import (
	"context"
	"encoding/json"

	"github.com/AltSumpreme/Medistream.git/queue"
	"github.com/AltSumpreme/Medistream.git/services/mail"

	"github.com/hibiken/asynq"
)

func RegisterEmailHandlers(mux *asynq.ServeMux) {
	mux.HandleFunc(string(queue.JobTypeWelcomeEmail), handleWelcomeEmail)
	mux.HandleFunc(string(queue.JobOTPEmail), handleOTPEmail)

}

func handleWelcomeEmail(ctx context.Context, task *asynq.Task) error {
	var p queue.WelcomePayload
	if err := json.Unmarshal(task.Payload(), &p); err != nil {
		return err
	}

	return mail.SendEmail(
		p.Email,
		"Welcome to Medistream",
		"Hi "+p.Name+", welcome aboard!",
	)
}

func handleOTPEmail(ctx context.Context, task *asynq.Task) error {
	var p queue.OTPPayload
	if err := json.Unmarshal(task.Payload(), &p); err != nil {
		return err
	}

	return mail.SendEmail(
		p.Email,
		"Your OTP Code",
		"Your OTP is: "+p.OTP,
	)
}

func handleresetPasswordEmail(body string, ctx context.Context, task *asynq.Task) error {

	var p queue.OTPPayload
	if err := json.Unmarshal(task.Payload(), &p); err != nil {
		return err
	}
	return mail.SendEmail(
		p.Email,
		"Reset Your Password",
		body,
	)

}
