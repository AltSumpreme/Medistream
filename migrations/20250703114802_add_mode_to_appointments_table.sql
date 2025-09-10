-- +goose Up
-- +goose StatementBegin
ALTER TABLE appointments ADD COLUMN mode VARCHAR(255) NOT NULL DEFAULT 'Online';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE appointments DROP COLUMN mode;
-- +goose StatementEnd
