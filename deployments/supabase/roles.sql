CREATE EXTENSION IF NOT EXISTS vector WITH SCHEMA public;

-- Create backend user role

CREATE ROLE coach WITH LOGIN PASSWORD 'YOUR_PASSWORD_HERE' -- Replace with actual password from secret manager
 NOCREATEDB NOCREATEROLE BYPASSRLS VALID UNTIL 'infinity';

-- Create frontend user role

CREATE ROLE swimmer WITH LOGIN PASSWORD 'YOUR_PASSWORD_HERE' -- Replace with actual password from secret manager
 NOCREATEDB NOCREATEROLE BYPASSRLS VALID UNTIL 'infinity';

-- Grant schema privileges to postgres (owner)
GRANT USAGE,
CREATE ON SCHEMA public TO postgres;

-- Grant schema privileges to coach
GRANT USAGE,
CREATE ON SCHEMA public TO coach;

-- Grant schema privileges to swimmer
GRANT USAGE ON SCHEMA public TO swimmer;

-- Grant coach privileges on tables
GRANT
SELECT,
INSERT,
UPDATE,
DELETE,
TRUNCATE, REFERENCES,
          TRIGGER ON ALL TABLES IN SCHEMA public TO coach;


ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT
SELECT,
INSERT,
UPDATE,
DELETE,
TRUNCATE, REFERENCES,
          TRIGGER ON TABLES TO coach;

-- Grant coach privileges on sequences
GRANT USAGE,
SELECT ON ALL SEQUENCES IN SCHEMA public TO coach;


ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT USAGE,
SELECT ON SEQUENCES TO coach;

-- Grant coach database connection
GRANT CONNECT ON DATABASE postgres TO coach;

-- Grant swimmer privileges on tables
GRANT
SELECT, REFERENCES,
        TRIGGER ON ALL TABLES IN SCHEMA public TO swimmer;


ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT
SELECT, REFERENCES,
        TRIGGER ON TABLES TO swimmer;

-- Grant swimmer privileges on sequences
GRANT USAGE,
SELECT ON ALL SEQUENCES IN SCHEMA public TO swimmer;


ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT USAGE,
SELECT ON SEQUENCES TO swimmer;

-- Grant swimmer database connection
GRANT CONNECT ON DATABASE postgres TO swimmer;

-- Allow postgres to SET ROLE to coach and swimmer
GRANT coach TO postgres;

GRANT swimmer TO postgres;

-- Revoke CREATE privilege from PUBLIC role on public schema
REVOKE
CREATE ON SCHEMA public
FROM PUBLIC;
