-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
  DROP COLUMN email,
  DROP COLUMN password,
    ADD CONSTRAINT fk_auth
        FOREIGN KEY(auth_id)
        REFERENCES auth(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
  ADD COLUMN email VARCHAR(255) NOT NULL UNIQUE,
  ADD COLUMN password VARCHAR(255) NOT NULL,
  DROP CONSTRAINT fk_auth;
-- +goose StatementEnd
