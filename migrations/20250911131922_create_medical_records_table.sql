-- +goose Up
CREATE TABLE IF NOT EXISTS medical_records (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    patient_id UUID NOT NULL,
    doctor_id UUID NOT NULL,
    diagnosis TEXT NOT NULL,
    notes TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP,

    CONSTRAINT fk_med_patient FOREIGN KEY(patient_id) REFERENCES patients(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,
    CONSTRAINT fk_med_doctor FOREIGN KEY(doctor_id) REFERENCES doctors(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS medical_records CASCADE;
