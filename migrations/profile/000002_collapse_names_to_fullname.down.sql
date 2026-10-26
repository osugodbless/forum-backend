ALTER TABLE profile.profiles
  DROP COLUMN display_name,
  DROP COLUMN full_name,
  ADD COLUMN first_name text NOT NULL,
  ADD COLUMN text NOT NULL;