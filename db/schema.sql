CREATE SCHEMA IF NOT EXISTS public;

CREATE TYPE rule_scope AS ENUM (
    'database',
    'schema',
    'table',
    'column'
);

CREATE TYPE migration_status AS ENUM (
    'pending',
    'in_progress',
    'completed',
    'failed'
);

CREATE TYPE migration_statement_action AS ENUM (
    'create',
    'alter',
    'drop',
    'insert',
    'update',
    'delete'
);

CREATE TYPE rule_type AS ENUM (
    'lock',
    'warn'
);

CREATE TYPE database_type AS ENUM (
    'postgres'
);

CREATE TABLE  users
(
    id       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email    TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL
);

CREATE TABLE projects
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE projects_secret_keys
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    secret_key TEXT NOT NULL,
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT unique_secret_key_name_per_project UNIQUE (name, project_id)
);

CREATE TABLE databases
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    display_name TEXT NOT NULL,
    db_name     text NOT NULL,
    username TEXT,
    password TEXT,
    host     TEXT,
    port     INTEGER NOT NULL DEFAULT 5432,
    type     database_type NOT NULL DEFAULT 'postgres',
    project_id uuid NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    CONSTRAINT unique_database_name_per_project UNIQUE (display_name, project_id)
);

CREATE TABLE schemas
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    database_id uuid NOT NULL REFERENCES databases(id) ON DELETE CASCADE
);

CREATE TABLE migration_statements
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    schema_id uuid NOT NULL REFERENCES schemas(id) ON DELETE CASCADE,
    raw_statement TEXT NOT NULL,
    action migration_statement_action NOT NULL,
    status migration_status NOT NULL DEFAULT 'pending',
    target TEXT,
    relation_name TEXT
);


CREATE TABLE tables
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    schema_id uuid NOT NULL REFERENCES databases(id) ON DELETE CASCADE
);


-- Rules specify the protection mechanisms for databases, schemas, tables, and columns.
CREATE TABLE rules
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    type rule_type    NOT NULL,
    scope rule_scope   NOT NULL,
    relation_name TEXT,
    target text NOT NULL,
    database_id uuid NOT NULL REFERENCES databases(id) ON DELETE CASCADE
);