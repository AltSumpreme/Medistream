-- +goose Up
CREATE TABLE IF NOT EXISTS appointments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    patient_id UUID NOT NULL,
    doctor_id UUID NOT NULL,
    appointment_date TIMESTAMP NOT NULL,
    start_time VARCHAR(10) NOT NULL,
    end_time VARCHAR(10) NOT NULL,
    status appointment_status NOT NULL DEFAULT 'PENDING',
    location VARCHAR(255),
    mode TEXT NOT NULL DEFAULT 'Online',
    appointment_type appt_type NOT NULL DEFAULT 'CONSULTATION',
    notes TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),

    CONSTRAINT fk_patient FOREIGN KEY(patient_id) REFERENCES patients(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,
    CONSTRAINT fk_doctor FOREIGN KEY(doctor_id) REFERENCES doctors(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS appointments CASCADE;
