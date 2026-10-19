CREATE SCHEMA IF NOT EXISTS profile;

CREATE TABLE IF NOT EXISTS profile.profiles (
    id uuid PRIMARY KEY,
    user_id uuid UNIQUE NOT NULL, --Logical reference to auth.users(id)
    first_name varchar(50),
    last_name varchar(50),
    country varchar(100),
    photo varchar(255),
    bio text,
    created_at timestamp(0) with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP
);