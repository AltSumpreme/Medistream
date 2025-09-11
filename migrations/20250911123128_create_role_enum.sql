-- +goose Up
-- +goose StatementBegin
CREATE TYPE role AS ENUM ('ADMIN', 'DOCTOR', 'PATIENT');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TYPE IF EXISTS role;
-- +goose StatementEnd
