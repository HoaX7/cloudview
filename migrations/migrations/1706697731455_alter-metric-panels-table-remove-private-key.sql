-- Up Migration
ALTER TABLE metric_panels DROP COLUMN private_key;

-- Down Migration
ALTER TABLE metric_panels ADD COLUMN private_key text;