-- +goose Up
-- +goose StatementBegin
ALTER TABLE vitals 
ADD COLUMN deleted_at TIMESTAMP NULL DEFAULT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE vitals 
DROP COLUMN deleted_at;
-- +goose StatementEnd
