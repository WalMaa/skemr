-- name: GetRule :one
SELECT *
FROM rules
WHERE id = @id
LIMIT 1;

-- name: CreateRule :one
INSERT INTO rules
    (name, type, database_entity_id, project_id)
VALUES (@name, @type, @database_entity_id, @project_id)
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

-- name: ListRulesByDatabaseId :many
SELECT r.*
FROM rules r
LEFT JOIN database_entities de
    ON r.database_entity_id = de.id
    AND de.database_id = @database_id;


-- name: ListRulesByCriteria :many
SELECT *
FROM rules
WHERE project_id = @project_id
  AND (database_entity_id = @database_entity_id OR @database_entity_id IS NULL);