

CREATE TYPE user_status AS ENUM (
    'active',
    'banned',
    'suspended'
);


CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    first_name TEXT NOT NULL,
    middle_name TEXT NOT NULL,
    last_name TEXT NOT NULL,

    phone TEXT NOT NULL,
    email TEXT,

    fayda TEXT,

    language_code VARCHAR(3) NOT NULL DEFAULT 'en',

    password_hash TEXT NOT NULL,

    is_phone_verified BOOLEAN NOT NULL DEFAULT FALSE,
    is_fayda_verified BOOLEAN NOT NULL DEFAULT FALSE,

    status user_status NOT NULL DEFAULT 'active',

    deleted_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


CREATE UNIQUE INDEX idx_users_phone_active
    ON users(phone)
    WHERE deleted_at IS NULL;


CREATE UNIQUE INDEX idx_users_email_active
    ON users(email)
    WHERE deleted_at IS NULL
      AND email IS NOT NULL;


CREATE UNIQUE INDEX idx_users_fayda_active
    ON users(fayda)
    WHERE deleted_at IS NULL
      AND fayda IS NOT NULL;

CREATE INDEX idx_users_status
    ON users(status);



CREATE TRIGGER trg_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();