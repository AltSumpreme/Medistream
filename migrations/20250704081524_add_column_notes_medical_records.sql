-- +goose Up
-- +goose StatementBegin
ALTER TABLE medical_records
ADD COLUMN notes TEXT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE medical_records
DROP COLUMN notes;
-- +goose StatementEnd
