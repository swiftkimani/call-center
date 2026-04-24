-- name: GetQueueByID :one
SELECT * FROM queues WHERE id = $1;

-- name: ListQueues :many
SELECT * FROM queues ORDER BY name;

-- name: CreateQueue :one
INSERT INTO queues (name, description, skills_required, max_wait_seconds, sla_seconds)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateQueue :one
UPDATE queues
SET name = $2, description = $3, skills_required = $4,
    max_wait_seconds = $5, sla_seconds = $6
WHERE id = $1
RETURNING *;
