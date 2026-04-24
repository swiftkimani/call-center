package queries

import (
	"context"

	"github.com/google/uuid"
)

func (q *Queries) GetQueueByID(ctx context.Context, id uuid.UUID) (Queue, error) {
	var item Queue
	err := q.db.QueryRow(ctx, `SELECT id, name, description, skills_required, max_wait_seconds, sla_seconds, created_at FROM queues WHERE id = $1`, id).
		Scan(&item.ID, &item.Name, &item.Description, &item.SkillsRequired, &item.MaxWaitSeconds, &item.SlaSeconds, &item.CreatedAt)
	return item, err
}

func (q *Queries) ListQueues(ctx context.Context) ([]Queue, error) {
	rows, err := q.db.Query(ctx, `SELECT id, name, description, skills_required, max_wait_seconds, sla_seconds, created_at FROM queues ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Queue
	for rows.Next() {
		var item Queue
		if err := rows.Scan(&item.ID, &item.Name, &item.Description, &item.SkillsRequired, &item.MaxWaitSeconds, &item.SlaSeconds, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
