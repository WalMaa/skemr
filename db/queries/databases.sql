-- name: GetDatabase :one
SELECT *
FROM databases
WHERE id = $1 LIMIT 1;

-- name: GetDatabaseByName :one
SELECT *
FROM databases
WHERE name = $1 AND project_id = $2 LIMIT 1;

-- name: CreateDatabase :one
INSERT INTO databases (project_id, name)
VALUES ($1, $2) RETURNING *;

-- name: UpdateDatabase :exec
UPDATE databases
SET name = $2,
    username = $3,
    password = $4
WHERE id = $1
    RETURNING *;

-- name: DeleteDatabase :exec
DELETE
FROM databases
WHERE id = $1;

-- name: ListDatabasesByProject :many
SELECT *
FROM databases
WHERE project_id = $1;
