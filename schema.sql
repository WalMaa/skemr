CREATE TABLE projects
(
    id   BIGSERIAL PRIMARY KEY,
    name text NOT NULL
);

CREATE TABLE databases
(
    id       BIGSERIAL PRIMARY KEY,
    name     text NOT NULL,
    username varchar(255),
    password varchar(255)
);

CREATE TABLE schemas
(
    id   BIGSERIAL PRIMARY KEY,
    name varchar(255) NOT NULL
);

CREATE TABLE tables
(
    id   BIGSERIAL PRIMARY KEY,
    name varchar(255) NOT NULL
);

CREATE TABLE rules
(
    id   BIGSERIAL PRIMARY KEY,
    name varchar(255) NOT NULL,
    type rule_type    NOT NULL,
    scope rule_scope   NOT NULL,
    target text NOT NULL
);

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