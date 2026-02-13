-- Initialize Keycloak database
CREATE DATABASE IF NOT EXISTS keycloak;

-- Create additional schemas or tables if needed
-- Example:
-- CREATE SCHEMA IF NOT EXISTS keycloak_schema;

-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE keycloak TO keycloak;

-- You can add more initialization scripts here
