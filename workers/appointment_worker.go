package workers

import (
	"github.com/AltSumpreme/Medistream.git/controllers/appointments"
	"github.com/AltSumpreme/Medistream.git/queue"
)

func HandleAppointmentJobs(data *queue.JobPayload) error {

	switch data.Type {
	case queue.JobTypeCreateAppointment:
		appointments.CreateAppointment(data.Data)
	}
	return nil
}
