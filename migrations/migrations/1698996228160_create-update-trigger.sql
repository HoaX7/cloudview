-- Up Migration
CREATE OR REPLACE FUNCTION on_update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Down Migration
DROP FUNCTION IF EXISTS on_update_timestamp;
DROP EXTENSION IF EXISTS "uuid-ossp";