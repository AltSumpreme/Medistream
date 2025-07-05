-- +goose Up
-- +goose StatementBegin
ALTER TABLE prescriptions
ADD COLUMN deleted_at TIMESTAMP NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE prescriptions
DROP COLUMN deleted_at;
-- +goose StatementEnd
