-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_auth_email ON auth(email);
CREATE INDEX idx_user_auth_id ON users(auth_id);
CREATE INDEX idx_refresh_token_user_id ON refresh_tokens(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_auth_email;
DROP INDEX IF EXISTS idx_user_auth_id;
DROP INDEX IF EXISTS idx_refresh_token_user_id;
-- +goose StatementEnd
