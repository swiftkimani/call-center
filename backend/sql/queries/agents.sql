-- name: GetAgentByID :one
SELECT a.*, u.email, u.full_name, u.role
FROM agents a
JOIN users u ON u.id = a.user_id
WHERE a.id = $1;

-- name: GetAgentByUserID :one
SELECT * FROM agents WHERE user_id = $1;

-- name: CreateAgent :one
INSERT INTO agents (user_id, extension, skills, team_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateAgentStatus :exec
UPDATE agents
SET status = $2, last_seen_at = NOW()
WHERE id = $1;

-- name: UpdateAgentLastSeen :exec
UPDATE agents SET last_seen_at = NOW() WHERE id = $1;

-- name: ListAvailableAgents :many
SELECT * FROM agents
WHERE status = 'available'
ORDER BY last_seen_at ASC NULLS LAST;

-- name: ListAgentsByStatus :many
SELECT a.*, u.full_name, u.email
FROM agents a
JOIN users u ON u.id = a.user_id
WHERE a.status = $1
ORDER BY a.last_seen_at ASC;

-- name: ListAllAgents :many
SELECT a.*, u.full_name, u.email
FROM agents a
JOIN users u ON u.id = a.user_id
WHERE u.deleted_at IS NULL
ORDER BY u.full_name;
