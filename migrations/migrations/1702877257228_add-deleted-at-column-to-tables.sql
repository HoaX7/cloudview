-- Up Migration
ALTER TABLE projects ADD COLUMN deleted_at timestamptz default null;
ALTER TABLE project_members ADD COLUMN deleted_at timestamptz default null;
ALTER TABLE services ADD COLUMN deleted_at timestamptz default null;

-- Down Migration
ALTER TABLE projects DROP COLUMN deleted_at;
ALTER TABLE project_members DROP COLUMN deleted_at;
ALTER TABLE services DROP COLUMN deleted_at;
