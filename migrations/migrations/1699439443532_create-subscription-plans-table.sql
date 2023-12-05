-- Up Migration
CREATE TABLE IF NOT EXISTS subscription_plans(
    id serial primary key,
    name varchar(255) not null,
    description varchar(255),
    tier int,
    cost float, --base price
    currency varchar(10), --base currency
    metadata jsonb,
    is_deleted boolean default false,
    created_at timestamptz default NOW(),
    updated_at timestamptz default NOW()
);
CREATE TRIGGER subscription_plans_updated_at
  BEFORE UPDATE ON subscription_plans
  FOR EACH ROW
  EXECUTE PROCEDURE on_update_timestamp();

-- Down Migration
DROP TRIGGER IF EXISTS subscription_plans_updated_at ON subscription_plans;
DROP TABLE IF EXISTS subscription_plans;