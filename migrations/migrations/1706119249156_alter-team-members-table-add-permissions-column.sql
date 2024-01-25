-- Up Migration
ALTER TABLE project_members ALTER COLUMN permissions TYPE text;
ALTER TABLE project_members ALTER COLUMN permissions SET DEFAULT '7'; 

-- Down Migration
ALTER TABLE project_members ALTER COLUMN permissions DROP DEFAULT;
ALTER TABLE project_members ALTER COLUMN permissions TYPE jsonb;