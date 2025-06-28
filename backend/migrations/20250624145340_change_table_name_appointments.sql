-- +goose Up
-- +goose StatementBegin
ALTER TABLE appoinments RENAME TO appointments;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE appointments RENAME TO appoinments;
-- +goose StatementEnd
