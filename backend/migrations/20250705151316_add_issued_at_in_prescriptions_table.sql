-- +goose Up
-- +goose StatementBegin
ALTER TABLE prescriptions
ADD COLUMN issued_at TIMESTAMP NULL DEFAULT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE prescriptions
DROP COLUMN issued_at;
-- +goose StatementEnd
