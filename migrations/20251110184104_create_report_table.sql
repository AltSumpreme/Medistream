-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    file_url TEXT NOT NULL, -- stores S3 object key
    doctor_id UUID NOT NULL REFERENCES doctors(id) ON DELETE CASCADE,
    patient_id UUID NOT NULL REFERENCES patients(id) ON DELETE CASCADE,
    medical_record_id UUID REFERENCES medical_records(id) ON DELETE SET NULL
);
-- +goose StatementEnd

ALTER TABLE IF EXISTS reports
    ALTER COLUMN id SET DEFAULT gen_random_uuid();

-- +goose Down
DROP TABLE IF EXISTS reports;
