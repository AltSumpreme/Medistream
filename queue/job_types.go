package queue

type JobType string

const (
	JobUserSignUp JobType = "user_sign_up"
	JobUserLogin  JobType = "user_login"
)

type JobPayload struct {
	Type JobType `json:"type"`
	Data any     `json:"data"`
}
