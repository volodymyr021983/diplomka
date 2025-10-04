-- +migrate Down
ALTER TABLE user_profiles
DROP COLUMN IF EXISTS email;
