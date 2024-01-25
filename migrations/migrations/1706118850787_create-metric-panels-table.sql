-- Up Migration
CREATE TABLE IF NOT EXISTS metric_panels(
    id uuid default uuid_generate_v4() primary key,
    name varchar(255) not null,
    description varchar(255),
    panels jsonb,
    provider_account_id uuid not null,
    metadata jsonb,
    is_deleted boolean default false,
    created_at timestamptz default NOW(),
    updated_at timestamptz default NOW(),
    constraint fk_provider_account foreign key(provider_account_id) references provider_accounts(id)
);
CREATE TRIGGER metric_panels_updated_at
  BEFORE UPDATE ON metric_panels
  FOR EACH ROW
  EXECUTE PROCEDURE on_update_timestamp();

-- Down Migration
DROP TRIGGER IF EXISTS metric_panels_updated_at ON metric_panels;
DROP TABLE IF EXISTS metric_panels;