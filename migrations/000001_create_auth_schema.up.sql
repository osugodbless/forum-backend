CREATE SCHEMA IF NOT EXISTS auth;

CREATE TABLE IF NOT EXISTS auth.users (
    id uuid PRIMARY KEY,
    username text NOT NULL,
    email text NOT NULL,
    password_hash text NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP

    -- Uniqueness constraints
    CONSTRAINT uq_users_email UNIQUE (email),
    CONSTRAINT uq_users_username UNIQUE (username),

    -- Format constraints
    CONSTRAINT chk_users_username_length CHECK (char_length(username) BETWEEN 3 AND 30),
    CONSTRAINT chk_users_email_length CHECK (char_length(email) <= 254),

    -- Lowercase guardrails
    CONSTRAINT chk_users_email_lowercase CHECK (email = LOWER(email)),
    CONSTRAINT chk_users_username_lowercase CHECK (username = LOWER(username))
);