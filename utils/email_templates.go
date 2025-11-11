package utils

import (
	"fmt"
	"strings"
	"time"
)

type EmailTemplate struct {
	Subject string
	Body    string
}

func formatDuration(d time.Duration) string {
	if d >= time.Hour {
		h := int(d.Hours())
		if h == 1 {
			return "1 hour"
		}
		return fmt.Sprintf("%d hours", h)
	}

	m := int(d.Minutes())
	if m == 1 {
		return "1 minute"
	}
	return fmt.Sprintf("%d minutes", m)
}

func GetWelcomeEmailTemplate(name string) EmailTemplate {
	subject := GetEnvWithDefault(
		"WELCOME_EMAIL_SUBJECT",
		"Welcome to Medistream!",
	)

	bodyTemplate := GetEnvWithDefault(
		"WELCOME_EMAIL_BODY",
		"Hi <strong>{{.NAME}}</strong>,<br><br>Welcome to Medistream! We're excited to have you with us.",
	)

	body := strings.ReplaceAll(bodyTemplate, "{{.NAME}}", name)

	return EmailTemplate{
		Subject: subject,
		Body:    body,
	}
}

func GetEmailVerificationTemplate(otp string, duration time.Duration) EmailTemplate {
	subject := GetEnvWithDefault(
		"EMAIL_VERIFICATION_SUBJECT",
		"Verify your email address",
	)

	bodyTemplate := GetEnvWithDefault(
		"EMAIL_VERIFICATION_BODY",
		"Your OTP is: <strong>{{.OTP}}</strong>. It is valid for {{.DURATION}}.",
	)

	body := strings.ReplaceAll(bodyTemplate, "{{.OTP}}", otp)
	body = strings.ReplaceAll(body, "{{.DURATION}}", formatDuration(duration))

	return EmailTemplate{
		Subject: subject,
		Body:    body,
	}
}

func GetResendVerificationTemplate(otp string, duration time.Duration) EmailTemplate {
	subject := GetEnvWithDefault(
		"EMAIL_RESEND_VERIFICATION_SUBJECT",
		"Your new verification code",
	)

	bodyTemplate := GetEnvWithDefault(
		"EMAIL_RESEND_VERIFICATION_BODY",
		"Your new OTP is: <strong>{{.OTP}}</strong>. It is valid for {{.DURATION}}.",
	)

	body := strings.ReplaceAll(bodyTemplate, "{{.OTP}}", otp)
	body = strings.ReplaceAll(body, "{{.DURATION}}", formatDuration(duration))

	return EmailTemplate{
		Subject: subject,
		Body:    body,
	}
}

func GetForgotPasswordOTPTemplate(otp string, duration time.Duration) EmailTemplate {
	subject := GetEnvWithDefault(
		"EMAIL_FORGOT_PASSWORD_SUBJECT",
		"Your password reset code",
	)

	bodyTemplate := GetEnvWithDefault(
		"EMAIL_FORGOT_PASSWORD_BODY",
		"You requested to reset your password.<br>Your OTP is: <strong>{{.OTP}}</strong>.<br>Valid for {{.DURATION}}.",
	)

	body := strings.ReplaceAll(bodyTemplate, "{{.OTP}}", otp)
	body = strings.ReplaceAll(body, "{{.DURATION}}", formatDuration(duration))

	return EmailTemplate{
		Subject: subject,
		Body:    body,
	}
}

func GetPasswordResetSuccessTemplate() EmailTemplate {
	subject := GetEnvWithDefault(
		"EMAIL_PASSWORD_RESET_SUCCESS_SUBJECT",
		"Your password was successfully reset",
	)

	body := GetEnvWithDefault(
		"EMAIL_PASSWORD_RESET_SUCCESS_BODY",
		"Your password has been updated. If this wasn't you, contact support immediately.",
	)

	return EmailTemplate{
		Subject: subject,
		Body:    body,
	}
}
