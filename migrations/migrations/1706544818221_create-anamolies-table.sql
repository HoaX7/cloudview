-- Up Migration
CREATE TABLE IF NOT EXISTS anamolies(
    id uuid default uuid_generate_v4() primary key,
    name varchar(255) not null,
    description varchar(255),
    instance_id varchar(255),
    provider_account_id uuid not null,
    metadata jsonb,
    is_deleted boolean default false,
    created_at timestamptz default NOW(),
    updated_at timestamptz default NOW(),
    deleted_at timestamptz default null,
    constraint fk_provider_account foreign key(provider_account_id) references provider_accounts(id)
);
CREATE TRIGGER anamolies_updated_at
  BEFORE UPDATE ON anamolies
  FOR EACH ROW
  EXECUTE PROCEDURE on_update_timestamp();

-- Down Migration
DROP TRIGGER IF EXISTS anamolies_updated_at ON anamolies;
DROP TABLE IF EXISTS anamolies;