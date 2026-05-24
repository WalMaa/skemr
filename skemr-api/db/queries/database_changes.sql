-- name: CreateDatabaseChange :one
INSERT INTO database_changes
    (database_id, entity_id, action)
VALUES (@database_id, @entity_id, @action)
RETURNING *;

-- name: GetDatabaseChangeByDatabaseIdAndId :one
SELECT *
FROM database_changes c
WHERE c.id = @id
  AND c.database_id = @database_id
LIMIT 1;

-- name: GetDatabaseChangesByDatabaseIdAndId :many
SELECT *
FROM database_changes c
WHERE c.database_id = @database_id
ORDER BY c.created_at DESC
LIMIT sqlc.narg('limit')::int OFFSET sqlc.narg('offset')::int;
