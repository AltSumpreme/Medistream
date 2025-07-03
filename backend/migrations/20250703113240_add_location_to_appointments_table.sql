-- +goose Up
-- +goose StatementBegin
ALTER TABLE appointments ADD COLUMN location VARCHAR(255) NOT NULL DEFAULT 'Unknown';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE appointments DROP COLUMN location;
-- +goose StatementEnd
