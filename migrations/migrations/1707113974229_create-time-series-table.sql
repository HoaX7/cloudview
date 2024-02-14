-- Up Migration
CREATE TYPE metric_type as ENUM('CPU_USAGE', 'RAM_USAGE', 'DISK_OPERATIONS');
CREATE TABLE IF NOT EXISTS time_series(
    id SERIAL,
    truncated_timestamp bigint not null,
    series double precision[][] not null,
    type metric_type not null,
    metric_panel_id uuid not null,
    created_at timestamptz default NOW(),
    updated_at timestamptz default NOW(),
    constraint fk_metric_panel foreign key(metric_panel_id) references metric_panels(id),
    primary key (id, truncated_timestamp)
);
CREATE TRIGGER time_series_updated_at
  BEFORE UPDATE ON time_series
  FOR EACH ROW
  EXECUTE PROCEDURE on_update_timestamp();

-- Down Migration
DROP TYPE IF EXISTS metric_type cascade;
DROP TRIGGER IF EXISTS time_series_updated_at ON time_series;
DROP TABLE IF EXISTS time_series;