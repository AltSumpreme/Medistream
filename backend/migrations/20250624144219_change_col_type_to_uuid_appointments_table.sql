-- +goose Up
-- +goose StatementBegin
ALTER TABLE appoinments
  ALTER COLUMN patient_id TYPE UUID USING patient_id::uuid,
  ALTER COLUMN doctor_id TYPE UUID USING doctor_id::uuid;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE appoinments
  ALTER COLUMN patient_id TYPE TEXT
  USING patient_id::TEXT,
  ALTER COLUMN doctor_id TYPE TEXT
  USING doctor_id::TEXT;
-- +goose StatementEnd
