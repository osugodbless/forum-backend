CREATE SCHEMA IF NOT EXISTS profile;

CREATE TABLE IF NOT EXISTS profile.profiles (
    id UUID PRIMARY KEY,
    user_id UUID UNIQUE NOT NULL, --Logical reference to auth.users(id)
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    country VARCHAR(100),
    photo VARCHAR(255),
    bio TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);