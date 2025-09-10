-- +goose Up
-- +goose StatementBegin
ALTER TABLE prescriptions
ADD COLUMN medical_record_id UUID NULL;
ALTER TABLE prescriptions
ADD CONSTRAINT fk_prescriptions_medical_record
FOREIGN KEY (medical_record_id) REFERENCES medical_records(id)
ON UPDATE CASCADE ON DELETE SET NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE prescriptions
DROP CONSTRAINT fk_prescriptions_medical_record;
ALTER TABLE prescriptions
DROP COLUMN medical_record_id;
-- +goose StatementEnd
