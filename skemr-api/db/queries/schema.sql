

-- name: GetSchema :one
SELECT *
FROM schemas
WHERE id = $1 LIMIT 1;

-- name: GetSchemaByIdAndDatabase :one
SELECT *
FROM schemas
WHERE id = $1 AND database_id = $2 LIMIT 1;

-- name: GetSchemaByIdAndProject :one
SELECT s.*
FROM schemas s
JOIN databases d ON s.database_id = d.id
WHERE s.id = $1 AND d.project_id = $2 LIMIT 1;

-- name: GetSchemaByNameAndDatabase :one
SELECT *
FROM schemas
WHERE name = $1 AND database_id = $2 LIMIT 1;

-- name: CreateSchema :one
INSERT INTO schemas (database_id, name)
VALUES ($1, $2) RETURNING *;

-- name: UpdateSchema :exec
UPDATE schemas
SET name = $2
WHERE id = $1
	RETURNING *;

-- name: DeleteSchema :exec
DELETE
FROM schemas
WHERE id = $1;

-- name: ListSchemasByDatabase :many
SELECT *
FROM schemas
WHERE database_id = $1;

-- name: ListSchemasByProject :many
SELECT s.*
FROM schemas s
JOIN databases d ON s.database_id = d.id
WHERE d.project_id = $1;

