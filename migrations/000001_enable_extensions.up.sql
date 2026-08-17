-- Extensions needed platform-wide
CREATE EXTENSION IF NOT EXISTS pgcrypto;   -- gen_random_uuid()
CREATE EXTENSION IF NOT EXISTS postgis;    -- geography(Point) for locations

-- Shared trigger function: auto-update `updated_at` on any row change.
-- Reused across every table that has an updated_at column, instead of
-- relying on application code to remember to set it.
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;