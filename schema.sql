CREATE SCHEMA IF NOT EXISTS "public";

CREATE TYPE rule_scope AS ENUM (
    'database',
    'schema',
    'table',
    'column'
);

CREATE TYPE rule_type AS ENUM (
    'lock',
    'warn'
);

CREATE TABLE projects
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL
);

CREATE TABLE databases
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name     text NOT NULL,
    username TEXT,
    password TEXT,
    project_id uuid NOT NULL REFERENCES projects(id) ON DELETE CASCADE
    ADD CONSTRAINT unique_database_name_per_project UNIQUE (name, project_id)
);

CREATE TABLE schemas
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    database_id uuid NOT NULL REFERENCES databases(id) ON DELETE CASCADE
);

CREATE TABLE tables
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    schema_id uuid NOT NULL REFERENCES databases(id) ON DELETE CASCADE
);

CREATE TABLE rules
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    type rule_type    NOT NULL,
    scope rule_scope   NOT NULL,
    target text NOT NULL
);