-- name: GetProjectSecretKeys :many
SELECT id, project_id, name, expires_at, created_at, updated_at
FROM projects_secret_keys
WHERE project_id = $1
ORDER BY created_at DESC;

-- name: GetProjectSecretKeyByID :one
SELECT id, project_id, name, expires_at, created_at, updated_at
FROM projects_secret_keys
WHERE id = $1
  AND project_id = $2
LIMIT 1;

-- name: GetProjectBySecretKey :one
SELECT p.id, p.name, p.created_at, p.updated_at
FROM projects_secret_keys psk
JOIN projects p ON p.id = psk.project_id
WHERE psk.secret_key = $1
LIMIT 1;

-- name: CreateProjectSecretKey :one
INSERT INTO projects_secret_keys (project_id, name, secret_key, expires_at)
VALUES ($1, $2, $3, $4)
RETURNING id, project_id, name, secret_key, expires_at, created_at, updated_at;
