-- Up Migration
CREATE TABLE IF NOT EXISTS audit_logs(
    id uuid default uuid_generate_v4() primary key,
    name varchar(255) not null,
    action varchar(255) not null,
    description varchar(255),
    email varchar(100),
    user_id uuid not null,
    project_id uuid,
    metadata jsonb,
    is_deleted boolean default false,
    created_at timestamptz default NOW(),
    updated_at timestamptz default NOW(),
    constraint fk_user foreign key(user_id) references users(id),
    constraint fk_project foreign key(project_id) references projects(id)
);
CREATE TRIGGER audit_logs_updated_at
  BEFORE UPDATE ON audit_logs
  FOR EACH ROW
  EXECUTE PROCEDURE on_update_timestamp();

-- Down Migration
DROP TRIGGER IF EXISTS audit_logs_updated_at ON audit_logs;
DROP TABLE IF EXISTS audit_logs;