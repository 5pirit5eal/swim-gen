-- Add application metadata table for schema versioning and configuration
-- This table tracks schema versions and can store feature flags

create table if not exists app_metadata (
  key text primary key,
  value jsonb not null,
  description text,
  updated_at timestamptz default now()
);

-- Insert initial schema version
insert into app_metadata (key, value, description)
values (
  'schema_version',
  '{"version": 2, "applied_at": "2025-10-10", "description": "Added swimmer permissions and app_metadata table"}',
  'Current database schema version'
) on conflict (key) do update set
  value = excluded.value,
  updated_at = now();

-- Insert database info
insert into app_metadata (key, value, description)
values (
  'database_info',
  '{"provider": "supabase", "migration_from": "cloud_sql", "migration_date": "2025-10-10"}',
  'Database provider and migration information'
) on conflict (key) do nothing;

-- Create index for faster lookups
create index if not exists idx_app_metadata_updated_at on app_metadata(updated_at desc);

-- Comment for documentation
COMMENT ON TABLE app_metadata IS 'Application metadata including schema version, feature flags, and configuration';
COMMENT ON COLUMN app_metadata.key IS 'Unique identifier for the metadata entry';
COMMENT ON COLUMN app_metadata.value IS 'JSON value containing the metadata';
COMMENT ON COLUMN app_metadata.description IS 'Human-readable description of the metadata entry';

-- Verify the table was created successfully
DO $$
DECLARE
  schema_version_data jsonb;
BEGIN
  SELECT value INTO schema_version_data FROM app_metadata WHERE key = 'schema_version';
  RAISE NOTICE 'Schema version: %', schema_version_data->>'version';
  RAISE NOTICE 'app_metadata table created successfully with swimmer permissions';
END $$;
