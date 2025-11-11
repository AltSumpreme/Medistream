package workers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/AltSumpreme/Medistream.git/controllers/appointments"
	"github.com/AltSumpreme/Medistream.git/handlers"
	"github.com/hibiken/asynq"
)

func ProcessCreateAppointmentTask(ctx context.Context, t *asynq.Task) error {
	var input handlers.AppointmentInput

	if err := json.Unmarshal(t.Payload(), &input); err != nil {
		return fmt.Errorf("failed to decode appointment task: %v", err)
	}

	appointments.CreateAppointment(input)

	return nil
}
