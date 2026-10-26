ALTER TABLE profile.profiles
  DROP COLUMN first_name,
  DROP COLUMN last_name,
  ADD COLUMN full_name text NOT NULL,
  ADD COLUMN display_name text;