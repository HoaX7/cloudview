-- Up Migration
CREATE TYPE health_status as ENUM('DISCONNECTED', 'ACTIVE');
ALTER TABLE metric_panels ADD COLUMN private_key text;
ALTER TABLE metric_panels ADD health_status health_status default 'DISCONNECTED';
ALTER TABLE metric_panels ADD COLUMN instance_id varchar(255) not null default '';
-- Down Migration
DROP TYPE IF EXISTS health_status cascade;
ALTER TABLE metric_panels DROP COLUMN private_key;
ALTER TABLE metric_panels DROP COLUMN instance_id;
