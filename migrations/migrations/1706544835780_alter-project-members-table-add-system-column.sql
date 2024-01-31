-- Up Migration
ALTER TABLE project_members ADD COLUMN is_system boolean default false;

-- Down Migration
ALTER TABLE project_members DROP COLUMN is_system;