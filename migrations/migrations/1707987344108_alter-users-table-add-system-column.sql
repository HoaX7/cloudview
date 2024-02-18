-- Up Migration
ALTER TABLE users ADD COLUMN is_system boolean default false;
ALTER TABLE project_members DROP COLUMN is_system;

-- Down Migration
ALTER TABLE users DROP COLUMN is_system;
ALTER TABLE project_members ADD COLUMN is_system boolean default false;