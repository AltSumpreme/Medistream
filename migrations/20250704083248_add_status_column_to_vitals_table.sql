-- +goose Up
-- +goose StatementBegin
ALTER TABLE vitals
ADD COLUMN status TEXT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE vitals
DROP COLUMN status;
-- +goose StatementEnd
