-- +goose Up
CREATE SCHEMA IF NOT EXISTS public;

CREATE TYPE database_entity_type AS ENUM (
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

CREATE TYPE database_entity_status AS ENUM (
    'active',
    'deleted'
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
    'locked',
    'warn',
    'advisory',
    'deprecated'
    );

CREATE TYPE database_type AS ENUM (
    'postgres'
    );

CREATE TABLE users
(
    id       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email    TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL
);

CREATE TABLE projects
(
    id         UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    name       TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE projects_sources
(
    id         UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    project_id UUID        NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
    name       TEXT        NOT NULL,
    source     TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE project_access_tokens
(
    id         UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    project_id UUID        NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
    name       TEXT        NOT NULL,
    prefix     TEXT        NOT NULL, -- lookup key of the token
    hash       TEXT        NOT NULL,
    last_used  TIMESTAMPTZ,
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT unique_secret_key_name_per_project UNIQUE (name, project_id)
);

CREATE TABLE databases
(
    id                         UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    display_name               TEXT        NOT NULL,
    db_name                    TEXT,
    username                   TEXT,
    password                   TEXT,
    host                       TEXT,
    port                       INTEGER,
    ssl_mode                   TEXT        NOT NULL DEFAULT 'prefer',
    database_type              database_type        DEFAULT 'postgres',
    project_id                 uuid        NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
    last_synced_at             TIMESTAMPTZ,
    last_sync_error            TEXT,
    failed_connection_attempts INTEGER     NOT NULL DEFAULT 0,
    created_at                 TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at                 TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT unique_database_name_per_project UNIQUE (display_name, project_id)
);

CREATE TABLE migration_statements
(
    id            UUID PRIMARY KEY                    DEFAULT gen_random_uuid(),
    raw_statement TEXT                       NOT NULL,
    action        migration_statement_action NOT NULL,
    status        migration_status           NOT NULL DEFAULT 'pending',
    target        TEXT,
    relation_name TEXT
);


CREATE TABLE tables
(
    id        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name      TEXT NOT NULL,
    schema_id uuid NOT NULL REFERENCES databases (id) ON DELETE CASCADE
);

CREATE TABLE database_entities
(
    id            uuid PRIMARY KEY                DEFAULT gen_random_uuid(),
    fingerprint   text                   NOT NULL,                        -- this is used to track the same entity across syncs even if it is renamed.
    project_id    uuid                   NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
    database_id   uuid                   NOT NULL REFERENCES databases (id) ON DELETE CASCADE,
    status        database_entity_status NOT NULL DEFAULT 'active',
    deleted_at    TIMESTAMPTZ            NULL,                            -- Set when status is 'deleted' to track when it was deleted
    first_seen_at TIMESTAMPTZ            NOT NULL DEFAULT NOW(),          -- Track when we first saw this entity
    entity_type   database_entity_type   NOT NULL,
    parent_id     uuid                   NULL REFERENCES database_entities (id),

    -- generic identity at this node
    name          text                   NOT NULL,                        -- e.g. "public", "users", "email", "my_view"
    attributes    jsonb                  NOT NULL DEFAULT '{}'::jsonb,    -- Store any additional metadata about the entity here

    created_at    TIMESTAMPTZ            NOT NULL DEFAULT NOW(),

    UNIQUE NULLS NOT DISTINCT (database_id, name, entity_type, parent_id) -- Ensure we do not map the same entity twice, use NULLS NOT DISTINCT so parentless are not duplicated
);


-- Rules specify the protection mechanisms for databases, schemas, tables, and columns.
CREATE TABLE rules
(
    id                 UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    name               TEXT        NOT NULL,                     -- Defined by user
    type               rule_type   NOT NULL,
    attributes         jsonb       NOT NULL DEFAULT '{}'::jsonb, -- Metadata about the rule, removal_date for deprecated types for example
    database_entity_id uuid        NOT NULL REFERENCES database_entities (id) ON DELETE CASCADE,
    database_id        uuid        NOT NULL REFERENCES databases (id) ON DELETE CASCADE,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT unique_rule_name_per_database UNIQUE (name, database_id)
);