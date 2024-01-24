-- Up Migration
CREATE TYPE access_type as ENUM('ACCESS KEYS', 'CROSS ACCOUNT ROLE');
ALTER TABLE services ALTER COLUMN access_key_id DROP NOT NULL;
ALTER TABLE services ALTER COLUMN access_key_secret DROP NOT NULL;
ALTER TABLE services ALTER COLUMN rotation_secret_key DROP NOT NULL;
ALTER TABLE services ADD type access_type default 'ACCESS KEYS';
ALTER TABLE services ADD account_id text not null default 'Missing account id';
ALTER TABLE services ADD access_role text;
ALTER TABLE services ADD feature_access_permission text NOT NULL default '7';

DROP TRIGGER IF EXISTS services_updated_at ON services;

ALTER TABLE services RENAME TO provider_accounts;

CREATE OR REPLACE TRIGGER provider_accounts_updated_at
  BEFORE UPDATE ON provider_accounts
  FOR EACH ROW
  EXECUTE PROCEDURE on_update_timestamp();

-- Down Migration
DROP TYPE IF EXISTS access_type cascade;
ALTER TABLE provider_accounts DROP COLUMN account_id;
ALTER TABLE provider_accounts DROP COLUMN access_role;
ALTER TABLE provider_accounts DROP COLUMN feature_access_permission;
ALTER TABLE provider_accounts ALTER COLUMN access_key_id SET NOT NULL;
ALTER TABLE provider_accounts ALTER COLUMN access_key_secret SET NOT NULL;
ALTER TABLE provider_accounts ALTER COLUMN rotation_secret_key SET NOT NULL;

DROP TRIGGER IF EXISTS provider_accounts_updated_at ON provider_accounts;

ALTER TABLE provider_accounts RENAME TO services;

CREATE OR REPLACE TRIGGER services_updated_at
  BEFORE UPDATE ON services 
  FOR EACH ROW
  EXECUTE PROCEDURE on_update_timestamp();