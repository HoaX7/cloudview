-- Up Migration
ALTER TABLE metric_panels ADD COLUMN auth_key text;

-- Down Migration
ALTER TABLE metric_panels DROP COLUMN auth_key;