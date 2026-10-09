CREATE EXTENSION IF NOT EXISTS pgcrypto;

DO $$ BEGIN
    CREATE TYPE slot_state AS ENUM ('invalid', 'available', 'hold', 'booked', 'cancelled');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

DO $$ BEGIN
    CREATE TYPE payment_state AS ENUM ('invalid', 'pending', 'success', 'failed');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

DO $$ BEGIN
    CREATE TYPE reservation_state AS ENUM ('invalid', 'hold', 'booked', 'cancelled');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

CREATE TABLE IF NOT EXISTS clinics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    address TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS slots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    clinic_id UUID NOT NULL REFERENCES clinics(id),
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE NOT NULL,
    state slot_state NOT NULL,
    amount NUMERIC(10,2) NOT NULL DEFAULT 500.00,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_slots_clinic_id ON slots(clinic_id);

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    address TEXT NOT NULL,
    phone TEXT NOT NULL,
    email TEXT NOT NULL,
    medical_notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS payment_status (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    status payment_state NOT NULL,
    amount NUMERIC(10,2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS reservations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    clinic_id UUID NOT NULL REFERENCES clinics(id),
    slot_id UUID NOT NULL REFERENCES slots(id),
    payment_status_id UUID NOT NULL REFERENCES payment_status(id),
    status reservation_state NOT NULL DEFAULT 'hold',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Partial unique index ensures only one active (hold/booked) reservation per slot, while preserving audit history for cancelled/expired reservations
CREATE UNIQUE INDEX IF NOT EXISTS uq_active_slot_reservation ON reservations(slot_id) WHERE status IN ('hold', 'booked');
CREATE INDEX IF NOT EXISTS idx_reservations_user_id ON reservations(user_id);

-- Seed data
INSERT INTO clinics (id, name, address)
SELECT gen_random_uuid(), 'Healthy Clinic', '123 Health St'
WHERE NOT EXISTS (SELECT 1 FROM clinics);

INSERT INTO clinics (id, name, address)
SELECT gen_random_uuid(), 'Wellness Center', '456 Wellness Blvd'
WHERE (SELECT COUNT(*) FROM clinics) < 2;

DO $$
DECLARE
    c_id1 UUID;
    c_id2 UUID;
BEGIN
    SELECT id INTO c_id1 FROM clinics ORDER BY id LIMIT 1;
    SELECT id INTO c_id2 FROM clinics ORDER BY id OFFSET 1 LIMIT 1;

    IF NOT EXISTS (SELECT 1 FROM slots WHERE clinic_id = c_id1) THEN
        FOR i IN 1..10 LOOP
            INSERT INTO slots (clinic_id, start_time, end_time, state, amount)
            VALUES (
                c_id1,
                NOW() + (i || ' days')::INTERVAL,
                NOW() + (i || ' days 1 hours')::INTERVAL,
                'available',
                500.00
            );
        END LOOP;
    END IF;

    IF NOT EXISTS (SELECT 1 FROM slots WHERE clinic_id = c_id2) THEN
        FOR i IN 1..10 LOOP
            INSERT INTO slots (clinic_id, start_time, end_time, state, amount)
            VALUES (
                c_id2,
                NOW() + (i || ' days')::INTERVAL,
                NOW() + (i || ' days 1 hours')::INTERVAL,
                'available',
                750.00
            );
        END LOOP;
    END IF;
END $$;

-- Seed with fixed UUIDs so the frontend hardcoded USER_ID works immediately
INSERT INTO users (id, name, address, phone, email)
SELECT 'a0000000-0000-0000-0000-000000000001', 'John Doe', '789 User Ave', '555-1234', 'john@example.com'
WHERE NOT EXISTS (SELECT 1 FROM users WHERE id = 'a0000000-0000-0000-0000-000000000001');

INSERT INTO users (id, name, address, phone, email)
SELECT 'a0000000-0000-0000-0000-000000000002', 'Jane Smith', '101 User Blvd', '555-5678', 'jane@example.com'
WHERE NOT EXISTS (SELECT 1 FROM users WHERE id = 'a0000000-0000-0000-0000-000000000002');
