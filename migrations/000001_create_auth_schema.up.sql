CREATE SCHEMA IF NOT EXISTS auth;

CREATE TABLE IF NOT EXISTS auth.users (
    id uuid PRIMARY KEY,
    username varchar(50) NOT NULL,
    email varchar(100) NOT NULL,
    password_hash varchar(255) NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_users_email_lower ON auth.users (LOWER(email));
CREATE UNIQUE INDEX idx_users_username_lower ON auth.users (LOWER(username));