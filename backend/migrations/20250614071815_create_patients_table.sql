-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS patients(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,

 CONSTRAINT fk_user
   FOREIGN KEY(user_id)
    REFERENCES users(id)
    ON UPDATE CASCADE
    ON DELETE CASCADE);
-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin
DROP TABLE IF EXISTS patients;
-- +goose StatementEnd
