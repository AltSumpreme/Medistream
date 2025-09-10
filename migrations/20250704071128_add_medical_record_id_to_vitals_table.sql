-- +goose Up
-- +goose StatementBegin
ALTER TABLE vitals
ADD COLUMN medical_record_id UUID  NULL;

ALTER TABLE vitals
ADD CONSTRAINT fk_medical_record FOREIGN KEY (medical_record_id)
REFERENCES medical_records(id)
ON UPDATE CASCADE
ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE vitals
DROP CONSTRAINT fk_medical_record;

ALTER TABLE vitals
DROP COLUMN medical_record_id;
-- +goose StatementEnd
