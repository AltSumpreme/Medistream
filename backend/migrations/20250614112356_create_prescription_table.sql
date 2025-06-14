-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS prescriptions(
     id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        patient_id UUID NOT NULL,
        doctor_id UUID NOT NULL,
        prescription_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        medication TEXT NOT NULL,
        dosage TEXT NOT NULL,
        instructions TEXT NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
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
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS prescriptions;
-- +goose StatementEnd
