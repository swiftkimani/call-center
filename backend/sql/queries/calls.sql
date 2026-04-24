-- name: CreateCall :one
INSERT INTO calls (provider_sid, customer_id, agent_id, queue_id, direction, status, from_number, to_number)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetCallByID :one
SELECT * FROM calls WHERE id = $1;

-- name: GetCallByProviderSID :one
SELECT * FROM calls WHERE provider_sid = $1;

-- name: UpdateCallStatus :exec
UPDATE calls SET status = $2 WHERE id = $1;

-- name: AnswerCall :exec
UPDATE calls
SET status = 'in_progress', agent_id = $2, answered_at = NOW(),
    wait_seconds = EXTRACT(EPOCH FROM (NOW() - started_at))::INT
WHERE id = $1;

-- name: EndCall :exec
UPDATE calls
SET status = $2, ended_at = NOW(),
    talk_seconds = CASE
      WHEN answered_at IS NOT NULL
      THEN EXTRACT(EPOCH FROM (NOW() - answered_at))::INT
      ELSE NULL
    END,
    cost_cents = $3
WHERE id = $1;

-- name: UpdateCallRecording :exec
UPDATE calls SET recording_url = $2 WHERE id = $1;

-- name: ListCallsByAgent :many
SELECT * FROM calls
WHERE agent_id = $1
ORDER BY started_at DESC
LIMIT $2 OFFSET $3;

-- name: ListCallsByCustomer :many
SELECT * FROM calls
WHERE customer_id = $1
ORDER BY started_at DESC
LIMIT $2 OFFSET $3;

-- name: ListCalls :many
SELECT c.*, cu.full_name AS customer_name, cu.phone_number AS customer_phone
FROM calls c
LEFT JOIN customers cu ON cu.id = c.customer_id
ORDER BY c.started_at DESC
LIMIT $1 OFFSET $2;

-- name: InsertCallEvent :exec
INSERT INTO call_events (call_id, event_type, payload)
VALUES ($1, $2, $3);

-- name: DailySummary :many
SELECT
  agent_id,
  COUNT(*)                                    AS total_calls,
  COUNT(*) FILTER (WHERE status = 'completed') AS completed_calls,
  COUNT(*) FILTER (WHERE status = 'abandoned') AS abandoned_calls,
  AVG(talk_seconds) FILTER (WHERE talk_seconds IS NOT NULL) AS avg_talk_seconds,
  AVG(wait_seconds) FILTER (WHERE wait_seconds IS NOT NULL) AS avg_wait_seconds,
  SUM(cost_cents)                             AS total_cost_cents
FROM calls
WHERE started_at >= $1 AND started_at < $2
GROUP BY agent_id;
