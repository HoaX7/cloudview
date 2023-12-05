-- Up Migration
CREATE TABLE IF NOT EXISTS users(
    id uuid default uuid_generate_v4() primary key,
    username varchar(255) not null,
    email varchar(255) not null,
    avatar_url varchar(255),
    subscribed_since timestamp,
    subscription_days_left int,
    last_login_at timestamp,
    metadata jsonb,
    is_deleted boolean default false,
    created_at timestamp default NOW(),
    updated_at timestamp default NOW()
);
CREATE TRIGGER users_updated_at
  BEFORE UPDATE ON users
  FOR EACH ROW
  EXECUTE PROCEDURE on_update_timestamp();

-- Down Migration
DROP TRIGGER IF EXISTS users_updated_at ON users;
DROP TABLE IF EXISTS users;