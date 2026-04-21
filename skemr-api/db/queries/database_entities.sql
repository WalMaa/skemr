-- name: GetDatabaseEntity :one
SELECT *
FROM database_entities
WHERE id = @id
LIMIT 1;

-- name: GetDatabaseEntityByProjectIdAndId :one
SELECT *
FROM database_entities
WHERE id = @id
  AND project_id = @project_id
LIMIT 1;

-- name: GetDatabaseEntitiesByProjectId :many
SELECT *
FROM database_entities
WHERE project_id = @project_id;

-- name: UpdateDatabaseEntityAsDeleted :exec
UPDATE database_entities
SET status     = 'deleted',
    deleted_at = NOW()
WHERE id = @id;

-- name: GetDatabaseEntities :many
SELECT *
FROM database_entities
WHERE database_id = @database_id
  AND (entity_type = sqlc.narg('entity_type') OR sqlc.narg('entity_type') IS NULL)
  AND (parent_id = sqlc.narg('parent_id') OR sqlc.narg('parent_id') IS NULL);


-- name: GetDatabaseEntitiesByDatabaseId :many
SELECT *
FROM database_entities
WHERE database_id = @database_id;

-- name: GetDatabaseEntitiesByDatabaseIdAndParentId :many
SELECT *
FROM database_entities
WHERE database_id = @database_id
  AND parent_id = @parent_id;

-- name: GetDatabaseEntityByFingerprint :one
SELECT *
FROM database_entities
WHERE database_id = @database_id
  AND fingerprint = @fingerprint
LIMIT 1;

-- name: GetDatabaseEntityByDatabaseIdAndTypeAndParentAndName :one
SELECT *
FROM database_entities
WHERE database_id = @database_id
  AND entity_type = @entity_type
  AND parent_id IS NOT DISTINCT FROM @parent_id -- this can be null in case of schema so we use is not distinct to compare
  AND name = @name
LIMIT 1;

-- name: CreateDatabaseEntity :one
INSERT INTO database_entities
(project_id, database_id, entity_type, parent_id, name, attributes, fingerprint)
VALUES (@project_id, @database_id, @entity_type, @parent_id, @name, COALESCE(@attributes, '{}'::jsonb), @fingerprint)
RETURNING *;

-- name: UpdateDatabaseEntityName :one
UPDATE database_entities
SET name = @name
WHERE id = @id
RETURNING *;

-- name: UpdateDatabaseEntity :one
UPDATE database_entities
SET name = COALESCE(sqlc.narg(name), name),
    attributes = COALESCE(sqlc.narg(attributes), attributes),
    fingerprint = COALESCE(sqlc.narg(fingerprint), fingerprint),
    parent_id = COALESCE(sqlc.narg(parent_id), parent_id)
WHERE id = @id
RETURNING *;

