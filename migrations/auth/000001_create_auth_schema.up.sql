CREATE SCHEMA IF NOT EXISTS auth;

CREATE TABLE IF NOT EXISTS auth.users (
    id uuid PRIMARY KEY,
    email text NOT NULL,
    password_hash text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Uniqueness constraints
    CONSTRAINT uq_users_email UNIQUE (email),

    -- Format and lowercase guardrail
    CONSTRAINT chk_users_email_length CHECK (char_length(email) <= 254),
    CONSTRAINT chk_users_email_lowercase CHECK (email = LOWER(email))
);