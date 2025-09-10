-- +goose Up
-- +goose StatementBegin
ALTER TABLE appointments
ADD COLUMN appointment_time TIME NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE appointments
DROP COLUMN appointment_time;
-- +goose StatementEnd
