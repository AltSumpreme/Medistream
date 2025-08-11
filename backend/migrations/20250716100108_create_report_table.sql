-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS report (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    file_url VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    patient_id UUID NOT NULL,
    doctor_id UUID NOT NULL,
    medical_record_id UUID,
    deleted_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_patient
        FOREIGN KEY (patient_id)
        REFERENCES patients(id)
        ON UPDATE CASCADE,

    CONSTRAINT fk_doctor
        FOREIGN KEY (doctor_id)
        REFERENCES doctors(id)
        ON UPDATE CASCADE,

    CONSTRAINT fk_medical_record
        FOREIGN KEY (medical_record_id)
        REFERENCES medical_records(id)
        ON UPDATE CASCADE
        ON DELETE SET NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS report;
-- +goose StatementEnd
