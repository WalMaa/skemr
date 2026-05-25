-- name: GetPipelineRunByDatabaseIdAndId :one
SELECT *
FROM pipeline_runs
WHERE database_id = @database_id
  AND id = @id
LIMIT 1;

-- name: GetPipelineRunsByDatabaseId :many
SELECT *
FROM pipeline_runs
WHERE database_id = @database_id
ORDER BY created_at DESC;

-- name: CreatePipelineRun :one
INSERT INTO pipeline_runs
    (database_id, status, environment, completed_at)
VALUES (@database_id, @status, @environment, @completed_at)
RETURNING *;