-- name: GetProjectAccessTokens :many
SELECT *
FROM project_access_tokens
WHERE project_id = @project_id
ORDER BY created_at DESC;

-- name: GetProjectSecretKeyByID :one
SELECT id, project_id, name, expires_at, created_at, updated_at
FROM project_access_tokens
WHERE id = @id
  AND project_id = @project_id
LIMIT 1;

-- name: GetProjectBySecretPrefix :one
SELECT p.id, p.name, p.created_at, p.updated_at
FROM project_access_tokens pat
JOIN projects p ON p.id = pat.project_id
WHERE pat.prefix = @prefix
LIMIT 1;

-- name: GetHashByPrefixAndProjectID :one
SELECT hash
FROM project_access_tokens
WHERE prefix = @prefix
  AND project_id = @project_id
LIMIT 1;

-- name: CreateProjectSecretKey :one
INSERT INTO project_access_tokens (project_id, name, prefix, hash, expires_at)
VALUES (@project_id, @name, @prefix, @hash, @expires_at)
RETURNING id, project_id, name, expires_at, created_at, updated_at;


-- name: DeleteProjectAccessToken :exec
DELETE FROM project_access_tokens WHERE project_id = @project_id AND id = @secret_id;