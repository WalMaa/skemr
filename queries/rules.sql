-- name: GetRule :one
SELECT *
FROM rules
WHERE id = $1 LIMIT 1;

-- name: CreateRule :one
INSERT INTO rules
(name, type, scope, target, project_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateRule :exec
UPDATE rules
SET name = $2,
    type = $3,
    scope = $4,
    target = $5
WHERE id = $1
    RETURNING *;

-- name: DeleteRule :exec
DELETE
FROM rules
WHERE id = $1;

-- name: ListRulesByProject :many
SELECT *
FROM rules
WHERE project_id = $1;

-- name: ListRulesByCriteria :many
SELECT *
FROM rules
WHERE project_id = $1
    AND (scope = $2 OR $2 IS NULL)
    AND (type = $3 OR $3 IS NULL)
    AND (target = $4 OR $4 IS NULL);