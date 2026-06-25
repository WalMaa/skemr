-- +goose Up
CREATE SCHEMA IF NOT EXISTS public;

CREATE TYPE database_entity_type AS ENUM (
    'database',
    'schema',
    'table',
    'column',
    'index',
    'constraint'
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
    CONSTRAINT unique_database_name_per_project UNIQUE (display_name, project_id),
    CONSTRAINT databases_id_project_id_unique UNIQUE (id, project_id)
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
    project_id    uuid                   NOT NULL,
    database_id   uuid                   NOT NULL,
    status        database_entity_status NOT NULL DEFAULT 'active',
    deleted_at    TIMESTAMPTZ            NULL,                            -- Set when status is 'deleted' to track when it was deleted
    first_seen_at TIMESTAMPTZ            NOT NULL DEFAULT NOW(),          -- Track when we first saw this entity
    entity_type   database_entity_type   NOT NULL,
    parent_id     uuid                   NULL,

    -- generic identity at this node
    name          text                   NOT NULL,                        -- e.g. "public", "users", "email", "my_view"
    attributes    jsonb                  NOT NULL DEFAULT '{}'::jsonb,    -- Store any additional metadata about the entity here

    created_at    TIMESTAMPTZ            NOT NULL DEFAULT NOW(),

    CONSTRAINT database_entities_parent_same_database_fkey
        FOREIGN KEY (parent_id, database_id)
        REFERENCES database_entities (id, database_id) ON DELETE CASCADE, -- Ensure parent_id references an entity in the same database

    CONSTRAINT database_entities_database_project_fkey
        FOREIGN KEY (database_id, project_id)
            REFERENCES databases (id, project_id)
            ON DELETE CASCADE, -- Ensure entities are only in the same project

    CONSTRAINT database_entities_id_database_id_unique
        UNIQUE (id, database_id), -- for rule composite fkey

    UNIQUE NULLS NOT DISTINCT (database_id, name, entity_type, parent_id) -- Ensure we do not map the same entity twice, use NULLS NOT DISTINCT so parentless are not duplicated

);

CREATE TABLE database_changes
(
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    database_id   UUID NOT NULL REFERENCES databases (id) ON DELETE CASCADE,
    entity_id     UUID NOT NULL REFERENCES database_entities (id) ON DELETE CASCADE,
    action        migration_statement_action NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE pipeline_runs
(
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    database_id   UUID NOT NULL REFERENCES databases (id) ON DELETE CASCADE,
    status        migration_status NOT NULL DEFAULT 'pending',
    environment   TEXT,
    started_at    TIMESTAMPTZ,
    completed_at  TIMESTAMPTZ,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


-- Rules specify the protection mechanisms for databases, schemas, tables, and columns.
CREATE TABLE rules
(
    id                 UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    name               TEXT        NOT NULL,                     -- Defined by user
    type               rule_type   NOT NULL,
    attributes         jsonb       NOT NULL DEFAULT '{}'::jsonb, -- Metadata about the rule, removal_date for deprecated types for example
    database_entity_id uuid        NOT NULL REFERENCES database_entities (id),
    database_id        uuid        NOT NULL REFERENCES databases (id) ON DELETE CASCADE,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT unique_rule_name_per_database
        UNIQUE (name, database_id),
    CONSTRAINT rules_database_entity_database_fkey
        FOREIGN KEY (database_entity_id, database_id)
        REFERENCES database_entities (id, database_id) ON DELETE CASCADE --- ensure rules are only applied to entities in the same database
);