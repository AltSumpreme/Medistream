-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS goals(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    patient_id UUID NOT NULL,
    Type goaltype NOT NULL,
    TargetValue NUMERIC NOT NULL,
    CurrentValue NUMERIC NOT NULL DEFAULT 0,
    UpdatedAt TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_patient
        FOREIGN KEY(patient_id)
        REFERENCES patients(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS goals;
-- +goose StatementEnd
