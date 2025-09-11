-- +goose Up

DROP TYPE IF EXISTS vital_type;
DROP TYPE IF EXISTS appt_type;
DROP TYPE IF EXISTS appointment_status;

-- Create enums
CREATE TYPE appointment_status AS ENUM ('PENDING', 'CONFIRMED', 'CANCELLED', 'COMPLETED');
CREATE TYPE appt_type AS ENUM ('CONSULTATION', 'FOLLOW_UP', 'SURGERY', 'CHECKUP');
CREATE TYPE vital_type AS ENUM ('HEART_RATE', 'BLOOD_PRESSURE', 'TEMPERATURE', 'OXYGEN_SATURATION');

-- +goose Down

-- Drop enums in reverse order
DROP TYPE IF EXISTS vital_type;
DROP TYPE IF EXISTS appt_type;
DROP TYPE IF EXISTS appointment_status;
