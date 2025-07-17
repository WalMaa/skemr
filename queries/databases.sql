-- name: GetDatabase :one
SELECT *
FROM databases
WHERE id = $1 LIMIT 1;

-- name: CreateDatabase :one
INSERT INTO databases (name, username, password)
VALUES ($1, $2, $3) RETURNING *;

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
