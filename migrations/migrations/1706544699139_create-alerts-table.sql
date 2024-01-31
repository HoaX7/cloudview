-- Up Migration
CREATE TABLE IF NOT EXISTS alerts(
    id uuid default uuid_generate_v4() primary key,
    name varchar(255) not null,
    description varchar(255),
    configurations jsonb,
    provider_account_id uuid not null,
    metadata jsonb,
    is_deleted boolean default false,
    created_at timestamptz default NOW(),
    updated_at timestamptz default NOW(),
    deleted_at timestamptz default null,
    constraint fk_provider_account foreign key(provider_account_id) references provider_accounts(id)
);
CREATE TRIGGER alerts_updated_at
  BEFORE UPDATE ON alerts
  FOR EACH ROW
  EXECUTE PROCEDURE on_update_timestamp();

-- Down Migration
DROP TRIGGER IF EXISTS alerts_updated_at ON alerts;
DROP TABLE IF EXISTS alerts;