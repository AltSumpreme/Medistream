-- +goose Up
-- +goose StatementBegin
ALTER TABLE medical_records
ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE medical_records
DROP COLUMN deleted_at;
-- +goose StatementEnd
