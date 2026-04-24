package queries

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type CreateCampaignParams struct {
	Name        string
	ScheduledAt *time.Time
	CreatedBy   uuid.UUID
}

type UpdateCampaignStatusParams struct {
	ID     uuid.UUID
	Status string
}

type ListCampaignsParams struct {
	Limit  int32
	Offset int32
}

type InsertCampaignContactParams struct {
	CampaignID uuid.UUID
	CustomerID uuid.UUID
}

func (q *Queries) CreateCampaign(ctx context.Context, arg CreateCampaignParams) (Campaign, error) {
	var c Campaign
	err := q.db.QueryRow(ctx, `INSERT INTO campaigns (name, scheduled_at, created_by) VALUES ($1, $2, $3) RETURNING id, name, status, scheduled_at, created_by, created_at`, arg.Name, arg.ScheduledAt, arg.CreatedBy).
		Scan(&c.ID, &c.Name, &c.Status, &c.ScheduledAt, &c.CreatedBy, &c.CreatedAt)
	return c, err
}

func (q *Queries) GetCampaignByID(ctx context.Context, id uuid.UUID) (Campaign, error) {
	var c Campaign
	err := q.db.QueryRow(ctx, `SELECT id, name, status, scheduled_at, created_by, created_at FROM campaigns WHERE id = $1`, id).
		Scan(&c.ID, &c.Name, &c.Status, &c.ScheduledAt, &c.CreatedBy, &c.CreatedAt)
	return c, err
}

func (q *Queries) UpdateCampaignStatus(ctx context.Context, arg UpdateCampaignStatusParams) error {
	_, err := q.db.Exec(ctx, `UPDATE campaigns SET status = $2 WHERE id = $1`, arg.ID, arg.Status)
	return err
}

func (q *Queries) ListCampaigns(ctx context.Context, arg ListCampaignsParams) ([]Campaign, error) {
	rows, err := q.db.Query(ctx, `SELECT id, name, status, scheduled_at, created_by, created_at FROM campaigns ORDER BY created_at DESC LIMIT $1 OFFSET $2`, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Campaign
	for rows.Next() {
		var item Campaign
		if err := rows.Scan(&item.ID, &item.Name, &item.Status, &item.ScheduledAt, &item.CreatedBy, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (q *Queries) InsertCampaignContact(ctx context.Context, arg InsertCampaignContactParams) (CampaignContact, error) {
	var c CampaignContact
	err := q.db.QueryRow(ctx, `INSERT INTO campaign_contacts (campaign_id, customer_id) VALUES ($1, $2) RETURNING id, campaign_id, customer_id, status, attempted_at, completed_at`, arg.CampaignID, arg.CustomerID).
		Scan(&c.ID, &c.CampaignID, &c.CustomerID, &c.Status, &c.AttemptedAt, &c.CompletedAt)
	return c, err
}
