-- name: GetRule :one
SELECT *
FROM rules r
WHERE r.database_id = @database_id
  AND r.id = @rule_id
LIMIT 1;

-- name: GetRuleByDatabaseAndName :one
SELECT *
FROM rules
WHERE database_id = @database_id
  AND name = @name
LIMIT 1;

-- name: GetRuleWithEntity :one
SELECT sqlc.embed(r),
       sqlc.embed(de)
FROM rules r
         JOIN databases d ON r.database_id = d.id
         JOIN database_entities de ON r.database_entity_id = de.id
WHERE d.project_id = @project_id
  AND r.database_id = @database_id
  AND r.id = @rule_id
LIMIT 1;

-- name: CreateRule :one
INSERT INTO rules
    (name, type, database_entity_id, database_id, attributes)
VALUES (@name, @type, @database_entity_id, @database_id, COALESCE(@attributes, '{}'::jsonb))
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
WHERE database_id = @database_id
  AND id = @rule_id;


-- name: ListRulesByDatabaseId :many
SELECT *
FROM rules
WHERE database_id = @database_id;

-- name: GetRulesWithEntities :many
SELECT sqlc.embed(rules),
       sqlc.embed(database_entities)
FROM rules
         JOIN databases ON rules.database_id = databases.id
         JOIN database_entities ON rules.database_entity_id = database_entities.id
WHERE rules.database_id = @database_id
  AND databases.project_id = @project_id;


-- name: ListRulesByCriteria :many
SELECT *
FROM rules
WHERE database_id = @database_id
  AND (database_entity_id = @database_entity_id OR @database_entity_id IS NULL);