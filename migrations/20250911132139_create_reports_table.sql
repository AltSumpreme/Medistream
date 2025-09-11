-- +goose Up
CREATE TABLE IF NOT EXISTS reports (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    file_url TEXT NOT NULL,
    doctor_id UUID NOT NULL,
    patient_id UUID NOT NULL,
    medical_record_id UUID,
    deleted_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),

    CONSTRAINT fk_report_patient FOREIGN KEY(patient_id) REFERENCES patients(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,
    CONSTRAINT fk_report_doctor FOREIGN KEY(doctor_id) REFERENCES doctors(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,
    CONSTRAINT fk_report_medical FOREIGN KEY(medical_record_id) REFERENCES medical_records(id)
        ON UPDATE CASCADE
        ON DELETE SET NULL
);

-- +goose Down
DROP TABLE IF EXISTS reports CASCADE;
