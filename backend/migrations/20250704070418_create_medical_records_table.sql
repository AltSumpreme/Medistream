-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS medical_records (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    patient_id UUID NOT NULL,
    doctor_id UUID NOT NULL,
    diagnosis TEXT,
    treatment TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_patient
      FOREIGN KEY(patient_id)
      REFERENCES patients(id)
      ON UPDATE CASCADE
      ON DELETE CASCADE,
    CONSTRAINT fk_doctor
      FOREIGN KEY(doctor_id)
      REFERENCES doctors(id)
      ON UPDATE CASCADE
      ON DELETE CASCADE
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS medical_records;
-- +goose StatementEnd
