package queries

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type CreateCustomerParams struct {
	PhoneNumber string
	FullName    string
	Email       *string
	Tags        []string
}

type SearchCustomersParams struct {
	Column1 string
	Limit   int32
	Offset  int32
}

type UpdateCustomerParams struct {
	ID       uuid.UUID
	FullName string
	Email    *string
	Tags     []string
}

type CreateCallParams struct {
	ProviderSid string
	CustomerID  *uuid.UUID
	AgentID     *uuid.UUID
	QueueID     *uuid.UUID
	Direction   string
	Status      string
	FromNumber  string
	ToNumber    string
}

type AnswerCallParams struct {
	ID      uuid.UUID
	AgentID *uuid.UUID
}

type EndCallParams struct {
	ID        uuid.UUID
	Status    string
	CostCents *int32
}

type UpdateCallRecordingParams struct {
	ID           uuid.UUID
	RecordingUrl *string
}

type ListCallsParams struct {
	Limit  int32
	Offset int32
}

type SaveDispositionParams struct {
	CallID   uuid.UUID
	AgentID  uuid.UUID
	Category string
	Notes    string
}

type InsertCallEventParams struct {
	CallID    uuid.UUID
	EventType string
	Payload   []byte
}

func (q *Queries) GetCustomerByID(ctx context.Context, id uuid.UUID) (Customer, error) {
	return scanCustomer(q.db.QueryRow(ctx, `SELECT id, phone_number, full_name, email, tags, dnc_listed, timezone, created_at, updated_at FROM customers WHERE id = $1`, id))
}

func (q *Queries) GetCustomerByPhone(ctx context.Context, phone string) (Customer, error) {
	return scanCustomer(q.db.QueryRow(ctx, `SELECT id, phone_number, full_name, email, tags, dnc_listed, timezone, created_at, updated_at FROM customers WHERE phone_number = $1`, phone))
}

func (q *Queries) CreateCustomer(ctx context.Context, arg CreateCustomerParams) (Customer, error) {
	return scanCustomer(q.db.QueryRow(ctx, `INSERT INTO customers (phone_number, full_name, email, tags) VALUES ($1, $2, $3, $4) RETURNING id, phone_number, full_name, email, tags, dnc_listed, timezone, created_at, updated_at`, arg.PhoneNumber, arg.FullName, arg.Email, arg.Tags))
}

func (q *Queries) UpdateCustomer(ctx context.Context, arg UpdateCustomerParams) (Customer, error) {
	return scanCustomer(q.db.QueryRow(ctx, `UPDATE customers SET full_name = $2, email = $3, tags = $4, updated_at = NOW() WHERE id = $1 RETURNING id, phone_number, full_name, email, tags, dnc_listed, timezone, created_at, updated_at`, arg.ID, arg.FullName, arg.Email, arg.Tags))
}

func (q *Queries) SearchCustomers(ctx context.Context, arg SearchCustomersParams) ([]Customer, error) {
	rows, err := q.db.Query(ctx, `SELECT id, phone_number, full_name, email, tags, dnc_listed, timezone, created_at, updated_at FROM customers WHERE full_name ILIKE '%' || $1 || '%' OR phone_number ILIKE '%' || $1 || '%' OR email ILIKE '%' || $1 || '%' ORDER BY full_name LIMIT $2 OFFSET $3`, arg.Column1, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Customer
	for rows.Next() {
		item, err := scanCustomer(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (q *Queries) CreateCall(ctx context.Context, arg CreateCallParams) (Call, error) {
	return scanCall(q.db.QueryRow(ctx, `INSERT INTO calls (provider_sid, customer_id, agent_id, queue_id, direction, status, from_number, to_number) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, provider_sid, customer_id, agent_id, queue_id, direction, status, started_at, answered_at, ended_at, wait_seconds, talk_seconds, recording_url, cost_cents, from_number, to_number`, arg.ProviderSid, arg.CustomerID, arg.AgentID, arg.QueueID, arg.Direction, arg.Status, arg.FromNumber, arg.ToNumber))
}

func (q *Queries) GetCallByID(ctx context.Context, id uuid.UUID) (Call, error) {
	return scanCall(q.db.QueryRow(ctx, `SELECT id, provider_sid, customer_id, agent_id, queue_id, direction, status, started_at, answered_at, ended_at, wait_seconds, talk_seconds, recording_url, cost_cents, from_number, to_number FROM calls WHERE id = $1`, id))
}

func (q *Queries) GetCallByProviderSID(ctx context.Context, sid string) (Call, error) {
	return scanCall(q.db.QueryRow(ctx, `SELECT id, provider_sid, customer_id, agent_id, queue_id, direction, status, started_at, answered_at, ended_at, wait_seconds, talk_seconds, recording_url, cost_cents, from_number, to_number FROM calls WHERE provider_sid = $1`, sid))
}

func (q *Queries) AnswerCall(ctx context.Context, arg AnswerCallParams) error {
	_, err := q.db.Exec(ctx, `UPDATE calls SET status = 'in_progress', agent_id = $2, answered_at = NOW(), wait_seconds = EXTRACT(EPOCH FROM (NOW() - started_at))::INT WHERE id = $1`, arg.ID, arg.AgentID)
	return err
}

func (q *Queries) EndCall(ctx context.Context, arg EndCallParams) error {
	_, err := q.db.Exec(ctx, `UPDATE calls SET status = $2, ended_at = NOW(), talk_seconds = CASE WHEN answered_at IS NOT NULL THEN EXTRACT(EPOCH FROM (NOW() - answered_at))::INT ELSE NULL END, cost_cents = $3 WHERE id = $1`, arg.ID, arg.Status, arg.CostCents)
	return err
}

func (q *Queries) UpdateCallRecording(ctx context.Context, arg UpdateCallRecordingParams) error {
	_, err := q.db.Exec(ctx, `UPDATE calls SET recording_url = $2 WHERE id = $1`, arg.ID, arg.RecordingUrl)
	return err
}

func (q *Queries) ListCalls(ctx context.Context, arg ListCallsParams) ([]ListCallsRow, error) {
	rows, err := q.db.Query(ctx, `SELECT c.id, c.provider_sid, c.customer_id, c.agent_id, c.queue_id, c.direction, c.status, c.started_at, c.answered_at, c.ended_at, c.wait_seconds, c.talk_seconds, c.recording_url, c.cost_cents, c.from_number, c.to_number, cu.full_name AS customer_name, cu.phone_number AS customer_phone FROM calls c LEFT JOIN customers cu ON cu.id = c.customer_id ORDER BY c.started_at DESC LIMIT $1 OFFSET $2`, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []ListCallsRow
	for rows.Next() {
		var item ListCallsRow
		err := rows.Scan(
			&item.ID,
			&item.ProviderSid,
			&item.CustomerID,
			&item.AgentID,
			&item.QueueID,
			&item.Direction,
			&item.Status,
			&item.StartedAt,
			&item.AnsweredAt,
			&item.EndedAt,
			&item.WaitSeconds,
			&item.TalkSeconds,
			&item.RecordingUrl,
			&item.CostCents,
			&item.FromNumber,
			&item.ToNumber,
			&item.CustomerName,
			&item.CustomerPhone,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (q *Queries) SaveDisposition(ctx context.Context, arg SaveDispositionParams) (uuid.UUID, error) {
	var id uuid.UUID
	err := q.db.QueryRow(ctx, `INSERT INTO dispositions (call_id, agent_id, category, notes) VALUES ($1, $2, $3, $4) RETURNING id`, arg.CallID, arg.AgentID, arg.Category, arg.Notes).Scan(&id)
	return id, err
}

func (q *Queries) InsertCallEvent(ctx context.Context, arg InsertCallEventParams) error {
	payload := json.RawMessage(arg.Payload)
	_, err := q.db.Exec(ctx, `INSERT INTO call_events (call_id, event_type, payload) VALUES ($1, $2, $3)`, arg.CallID, arg.EventType, payload)
	return err
}

type DailySummaryParams struct {
	StartedAt   time.Time
	StartedAt_2 time.Time
}

func (q *Queries) DailySummary(ctx context.Context, arg DailySummaryParams) ([]DailySummaryRow, error) {
	rows, err := q.db.Query(ctx, `SELECT agent_id, COUNT(*) AS total_calls, COUNT(*) FILTER (WHERE status = 'completed') AS completed_calls, COUNT(*) FILTER (WHERE status = 'abandoned') AS abandoned_calls, AVG(talk_seconds) FILTER (WHERE talk_seconds IS NOT NULL) AS avg_talk_seconds, AVG(wait_seconds) FILTER (WHERE wait_seconds IS NOT NULL) AS avg_wait_seconds, SUM(cost_cents) AS total_cost_cents FROM calls WHERE started_at >= $1 AND started_at < $2 GROUP BY agent_id`, arg.StartedAt, arg.StartedAt_2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []DailySummaryRow
	for rows.Next() {
		var item DailySummaryRow
		if err := rows.Scan(&item.AgentID, &item.TotalCalls, &item.CompletedCalls, &item.AbandonedCalls, &item.AvgTalkSeconds, &item.AvgWaitSeconds, &item.TotalCostCents); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
