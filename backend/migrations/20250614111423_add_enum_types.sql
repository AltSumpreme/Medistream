-- +goose Up
-- +goose StatementBegin

-- Create Role enum
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'role') THEN
        CREATE TYPE role AS ENUM ('ADMIN', 'PATIENT', 'DOCTOR', 'RECEPTIONIST');
    END IF;
END$$;
-- Create AppointmentStatus Enum
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'appointment_status') THEN
        CREATE TYPE appointment_status AS ENUM ('PENDING', 'CONFIRMED', 'CANCELLED', 'COMPLETED');
    END IF;
END $$;

-- Create VitalType Enum
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'vital_type') THEN
        CREATE TYPE vital_type AS ENUM ('BLOOD_PRESSURE', 'HEART_RATE', 'WEIGHT', 'BMI');

    END IF;
END $$;

-- Create GoalType enum
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'goaltype') THEN
        CREATE TYPE goaltype AS ENUM ('STEPS', 'WATER', 'SLEEP');
    END IF;
END$$;


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TYPE IF EXISTS goaltype;
DROP TYPE IF EXISTS vital_type;
DROP TYPE IF EXISTS appointment_status;
DROP TYPE IF EXISTS role;
-- +goose StatementEnd
