-- name: GetDatabase :one
SELECT *
FROM databases
WHERE id = $1 LIMIT 1;

-- name: GetDatabaseByIDAndProjectID :one
SELECT *
FROM databases
WHERE id = @id
  AND project_id = @project_id LIMIT 1;

-- name: GetDatabaseByNameAndProject :one
SELECT *
FROM databases
WHERE display_name = $1 AND project_id = $2 LIMIT 1;

-- name: GetDatabaseByIdAndProject :one
SELECT *
FROM databases
WHERE id = $1 AND project_id = $2 LIMIT 1;

-- name: CreateDatabase :one
INSERT INTO databases (project_id, display_name)
VALUES (@project_id, @display_name) RETURNING *;

-- name: UpdateDatabase :exec
UPDATE databases
SET display_name = $2,
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
