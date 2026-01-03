-- name: GetRule :one
SELECT *
FROM rules
WHERE id = @id
LIMIT 1;

-- name: CreateRule :one
INSERT INTO rules
    (name, type, project_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateRule :exec
UPDATE rules
SET name = $2,
    type = $3
WHERE id = $1
RETURNING *;

-- name: DeleteRule :exec
DELETE
FROM rules
WHERE id = @id;

-- name: ListRulesByProject :many
SELECT *
FROM rules
WHERE project_id = @id;

-- name: ListRulesByCriteria :many
SELECT *
FROM rules
WHERE project_id = @project_id
  AND (type = @type OR @type IS NULL);