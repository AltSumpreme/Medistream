-- +goose Up
-- +goose StatementBegin
ALTER TABLE appointments
ADD COLUMN notes TEXT NULL;

ALTER TABLE appointments
ADD COLUMN appointment_type VARCHAR(50) NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE appointments
DROP COLUMN notes;

ALTER TABLE appointments
DROP COLUMN appointment_type;
-- +goose StatementEnd
