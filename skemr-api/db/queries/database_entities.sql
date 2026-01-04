-- name: GetDatabaseEntity :one
SELECT *
FROM database_entities
WHERE id = @id
LIMIT 1;

-- name: GetDatabaseEntityByProjectIdAndId :one
SELECT *
FROM database_entities
WHERE id = @id AND project_id = @project_id
LIMIT 1;

-- name: GetDatabaseEntitiesByProjectId :many
SELECT *
FROM database_entities
WHERE project_id = @project_id;


-- name: GetDatabaseEntitiesByDatabaseIdAndType :many
SELECT  *
FROM database_entities
WHERE database_id = @database_id AND entity_type = @type;

-- name: GetDatabaseEntityByDatabaseIdAndTypeAndName :one
SELECT *
FROM database_entities
WHERE database_id = @database_id AND entity_type = @type AND name = @name
LIMIT 1;

-- name: CreateDatabaseEntity :one
INSERT INTO database_entities
(project_id, database_id, entity_type, parent_id, name)
VALUES
(@project_id, @database_id, @entity_type, @parent_id, @name)
RETURNING *;


