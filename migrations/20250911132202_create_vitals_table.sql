-- +goose Up
CREATE TABLE IF NOT EXISTS vitals (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    patient_id UUID NOT NULL,
    medical_record_id UUID,
    type vital_type NOT NULL,
    value VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL,
    recorded_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP,

    CONSTRAINT fk_vital_patient FOREIGN KEY(patient_id) REFERENCES patients(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,
    CONSTRAINT fk_vital_medical FOREIGN KEY(medical_record_id) REFERENCES medical_records(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS vitals CASCADE;
