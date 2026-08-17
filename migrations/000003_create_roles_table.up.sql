CREATE TYPE role_name AS ENUM (
    'user',
    'institution_staff',
    'admin'
);


CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    name role_name NOT NULL UNIQUE,

    description TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);



INSERT INTO roles (name, description)
VALUES
    (
        'user',
        'Default role for every registered account'
    ),
    (
        'institution_staff',
        'Manages an institution''s found-item desk'
    ),
    (
        'admin',
        'Full platform access'
    );