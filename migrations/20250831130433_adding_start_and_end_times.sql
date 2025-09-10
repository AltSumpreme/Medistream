-- +goose Up
-- +goose StatementBegin
ALTER TABLE appointments
    ADD COLUMN start_time TIME  NULL,
    ADD COLUMN end_time TIME  NULL,
    DROP COLUMN appointment_time,
    DROP COLUMN duration;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE appointments
    DROP COLUMN start_time,
    DROP COLUMN end_time,
    ADD COLUMN appointment_time TIME NOT NULL,
    ADD COLUMN duration INT NOT NULL;
-- +goose StatementEnd
