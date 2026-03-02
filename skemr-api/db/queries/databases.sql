-- name: GetDatabase :one
SELECT *
FROM databases
WHERE id = $1
LIMIT 1;

-- name: GetDatabaseByIDAndProjectID :one
SELECT *
FROM databases
WHERE id = @id
  AND project_id = @project_id
LIMIT 1;

-- name: GetDatabaseByNameAndProject :one
SELECT *
FROM databases
WHERE display_name = $1
  AND project_id = $2
LIMIT 1;

-- name: GetDatabaseByIdAndProject :one
SELECT *
FROM databases
WHERE id = $1
  AND project_id = $2
LIMIT 1;

-- name: CreateDatabase :one
INSERT INTO databases
(project_id, display_name, db_name, username, password, host, port, database_type)
VALUES (@project_id, @display_name, @db_name, @username, @password, @host, @port, @database_type)
RETURNING *;

-- name: UpdateDatabase :one
UPDATE databases
SET display_name  = COALESCE(sqlc.narg(display_name), display_name),
    db_name       = COALESCE(sqlc.narg(db_name), db_name),
    username      = COALESCE(sqlc.narg(username), username),
    password      = COALESCE(sqlc.narg(password), password),
    host          = COALESCE(sqlc.narg(host), host),
    port          = COALESCE(sqlc.narg(port), port),
    database_type = COALESCE(sqlc.narg(database_type), database_type)
WHERE id = @database_id
RETURNING *;

-- name: UpdateDatabaseSyncedAt :one
UPDATE databases
SET last_synced_at             = @synced_at,
    last_sync_error            = NULL,
    failed_connection_attempts = 0
WHERE id = @database_id
RETURNING *;

-- name: UpdateDatabaseSyncFail :one
UPDATE databases
SET last_sync_error            = @sync_error,
    failed_connection_attempts = failed_connection_attempts + 1,
    last_synced_at             = @synced_at
WHERE id = @database_id
RETURNING *;

-- name: DeleteDatabase :exec
DELETE
FROM databases
WHERE id = $1;

-- name: ListDatabasesByProject :many
SELECT *
FROM databases
WHERE project_id = $1;
