CREATE SCHEMA IF NOT EXISTS profile;

CREATE TABLE IF NOT EXISTS profile.profiles (
    user_id uuid PRIMARY KEY, --Logical reference to auth.users(id)
    username text NOT NULL,
    full_name text NOT NULL,
    display_name text,
    country char(2),
    bio text,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT uq_profile_username UNIQUE (username),
    CONSTRAINT chk_profile_username_length CHECK (char_length(trim(username)) BETWEEN 3 AND 30),
    CONSTRAINT chk_profile_username_lowercase CHECK (username = LOWER(username)),
    CONSTRAINT chk_profiles_bio_length CHECK (bio IS NULL OR char_length(trim(bio)) BETWEEN 3 AND 300)
);