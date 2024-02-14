-- Up Migration
ALTER TABLE users ADD COLUMN permissions text;

-- Down Migration
ALTER TABLE users DROP COLUMN permissions;