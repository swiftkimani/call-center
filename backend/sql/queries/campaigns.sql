-- name: CreateCampaign :one
INSERT INTO campaigns (name, scheduled_at, created_by)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetCampaignByID :one
SELECT * FROM campaigns WHERE id = $1;

-- name: UpdateCampaignStatus :exec
UPDATE campaigns SET status = $2 WHERE id = $1;

-- name: ListCampaigns :many
SELECT * FROM campaigns ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: InsertCampaignContact :one
INSERT INTO campaign_contacts (campaign_id, customer_id)
VALUES ($1, $2)
RETURNING *;

-- name: GetNextPendingContact :one
SELECT * FROM campaign_contacts
WHERE campaign_id = $1 AND status = 'pending'
ORDER BY id
LIMIT 1
FOR UPDATE SKIP LOCKED;

-- name: UpdateContactStatus :exec
UPDATE campaign_contacts
SET status = $2, attempted_at = NOW()
WHERE id = $1;

-- name: CompleteContact :exec
UPDATE campaign_contacts
SET status = 'completed', attempted_at = NOW(), completed_at = NOW()
WHERE id = $1;

-- name: CountContactsByStatus :many
SELECT status, COUNT(*) AS count
FROM campaign_contacts
WHERE campaign_id = $1
GROUP BY status;
