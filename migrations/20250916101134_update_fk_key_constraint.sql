-- +goose Up
-- +goose StatementBegin
ALTER TABLE doctor_working_hours
DROP CONSTRAINT IF EXISTS fk_doctor_working_hours_doctor_id;

ALTER TABLE doctor_working_hours
ADD CONSTRAINT fk_doctor_working_hours_doctor_id
FOREIGN KEY (doctor_id) REFERENCES doctors(id)
ON DELETE CASCADE
ON UPDATE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE doctor_working_hours
DROP CONSTRAINT IF EXISTS fk_doctor_working_hours_doctor_id;
ALTER TABLE doctor_working_hours
ADD CONSTRAINT fk_doctor_working_hours_doctor_id
FOREIGN KEY (doctor_id) REFERENCES users(id)
ON DELETE CASCADE
ON UPDATE CASCADE;
-- +goose StatementEnd
