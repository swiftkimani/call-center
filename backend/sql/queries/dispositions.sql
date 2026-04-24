-- name: SaveDisposition :one
INSERT INTO dispositions (call_id, agent_id, category, notes)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetDispositionByCall :one
SELECT * FROM dispositions WHERE call_id = $1;

-- name: ListDispositionCategories :many
SELECT * FROM disposition_categories WHERE active = true ORDER BY name;
