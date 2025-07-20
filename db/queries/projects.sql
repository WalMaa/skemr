-- name: GetProject :one
SELECT *
FROM projects
WHERE id = $1 LIMIT 1;

-- name: CreateProject :one
INSERT INTO projects (name)
VALUES ($1) RETURNING *;

-- name: UpdateProject :exec
UPDATE projects
set name = $2
WHERE id = $1
    RETURNING *;

-- name: DeleteProject :exec
DELETE
FROM projects
WHERE id = $1;