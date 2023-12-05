-- Up Migration
CREATE TABLE IF NOT EXISTS services(
    id uuid default uuid_generate_v4() primary key,
    name varchar(255) not null,
    description varchar(255),
    access_key_id text not null,
    access_key_secret text not null,
    rotation_secret_key text not null,
    provider varchar(255),
    project_id uuid not null,
    metadata jsonb,
    is_deleted boolean default false,
    created_at timestamptz default NOW(),
    updated_at timestamptz default NOW(),
    constraint fk_project foreign key(project_id) references projects(id)
);
CREATE TRIGGER services_updated_at
  BEFORE UPDATE ON services 
  FOR EACH ROW
  EXECUTE PROCEDURE on_update_timestamp();

-- Down Migration
DROP TRIGGER IF EXISTS services_updated_at ON services;
DROP TABLE IF EXISTS services;
