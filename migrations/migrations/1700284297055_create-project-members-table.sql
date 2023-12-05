-- Up Migration
CREATE TABLE IF NOT EXISTS project_members(
    id uuid default uuid_generate_v4() primary key,
    project_id uuid not null,
    user_id uuid not null,
    is_owner boolean default false,
    is_active boolean default true, -- If the project plan expires some users will lose access. Use this flag for indication.
    permissions jsonb, -- for future use for restricted access to projects.
    metadata jsonb,
    is_deleted boolean default false,
    created_at timestamptz default NOW(),
    updated_at timestamptz default NOW(),
    constraint fkm_project foreign key(project_id) references projects(id),
    constraint fkm_user foreign key(user_id) references users(id)
);
CREATE TRIGGER project_members_updated_at
  BEFORE UPDATE ON project_members 
  FOR EACH ROW
  EXECUTE PROCEDURE on_update_timestamp();
CREATE INDEX user_project_idx on project_members(user_id, project_id);

-- Down Migration
DROP INDEX user_project_idx;
DROP TRIGGER IF EXISTS project_members_updated_at ON project_members;
DROP TABLE IF EXISTS project_members;
