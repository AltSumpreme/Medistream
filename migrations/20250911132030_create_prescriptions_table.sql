-- +goose Up
CREATE TABLE IF NOT EXISTS prescriptions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    patient_id UUID,
    doctor_id UUID,
    medical_record_id UUID,
    medication TEXT,
    dosage TEXT,
    instructions TEXT,
    issued_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP,

    CONSTRAINT fk_presc_medical FOREIGN KEY(medical_record_id) REFERENCES medical_records(id)
        ON UPDATE CASCADE
        ON DELETE SET NULL,
    CONSTRAINT fk_presc_patient FOREIGN KEY(patient_id) REFERENCES patients(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,
    CONSTRAINT fk_presc_doctor FOREIGN KEY(doctor_id) REFERENCES doctors(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS prescriptions CASCADE;
