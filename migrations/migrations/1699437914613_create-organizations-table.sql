-- Up Migration
CREATE TYPE visiblity as ENUM('PUBLIC', 'PRIVATE');
CREATE TABLE IF NOT EXISTS projects(
    id uuid default uuid_generate_v4() primary key,
    name varchar(255) not null,
    description varchar(255),
    email varchar(100),
    owner_id uuid not null,
    members jsonb,
    member_limit int default 1,
    type visiblity default 'PRIVATE',
    metadata jsonb,
    is_deleted boolean default false,
    created_at timestamptz default NOW(),
    updated_at timestamptz default NOW(),
    constraint fk_owner foreign key(owner_id) references users(id)
);
CREATE TRIGGER projects_updated_at
  BEFORE UPDATE ON projects
  FOR EACH ROW
  EXECUTE PROCEDURE on_update_timestamp();

-- Down Migration
DROP TYPE IF EXISTS visiblity cascade;
DROP TRIGGER IF EXISTS projects_updated_at ON projects;
DROP TABLE IF EXISTS projects;
