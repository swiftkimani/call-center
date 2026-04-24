package queries

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DBTX interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
}

type Queries struct {
	db DBTX
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

type User struct {
	ID           uuid.UUID  `json:"id"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"password_hash"`
	FullName     string     `json:"full_name"`
	Role         string     `json:"role"`
	CreatedAt    time.Time  `json:"created_at"`
	DeletedAt    *time.Time `json:"deleted_at"`
}

type Agent struct {
	ID            uuid.UUID  `json:"id"`
	UserID        uuid.UUID  `json:"user_id"`
	Extension     string     `json:"extension"`
	Skills        []string   `json:"skills"`
	Status        string     `json:"status"`
	MaxConcurrent int16      `json:"max_concurrent"`
	TeamID        *uuid.UUID `json:"team_id"`
	LastSeenAt    *time.Time `json:"last_seen_at"`
}

type Customer struct {
	ID          uuid.UUID  `json:"id"`
	PhoneNumber string     `json:"phone_number"`
	FullName    string     `json:"full_name"`
	Email       *string    `json:"email"`
	Tags        []string   `json:"tags"`
	DncListed   bool       `json:"dnc_listed"`
	Timezone    string     `json:"timezone"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type Call struct {
	ID           uuid.UUID  `json:"id"`
	ProviderSid  string     `json:"provider_sid"`
	CustomerID   *uuid.UUID `json:"customer_id"`
	AgentID      *uuid.UUID `json:"agent_id"`
	QueueID      *uuid.UUID `json:"queue_id"`
	Direction    string     `json:"direction"`
	Status       string     `json:"status"`
	StartedAt    time.Time  `json:"started_at"`
	AnsweredAt   *time.Time `json:"answered_at"`
	EndedAt      *time.Time `json:"ended_at"`
	WaitSeconds  *int32     `json:"wait_seconds"`
	TalkSeconds  *int32     `json:"talk_seconds"`
	RecordingUrl *string    `json:"recording_url"`
	CostCents    *int32     `json:"cost_cents"`
	FromNumber   string     `json:"from_number"`
	ToNumber     string     `json:"to_number"`
}

type Campaign struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Status      string     `json:"status"`
	ScheduledAt *time.Time `json:"scheduled_at"`
	CreatedBy   uuid.UUID  `json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
}

type Queue struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	SkillsRequired []string  `json:"skills_required"`
	MaxWaitSeconds int32     `json:"max_wait_seconds"`
	SlaSeconds     int32     `json:"sla_seconds"`
	CreatedAt      time.Time `json:"created_at"`
}

type CampaignContact struct {
	ID          uuid.UUID  `json:"id"`
	CampaignID  uuid.UUID  `json:"campaign_id"`
	CustomerID  uuid.UUID  `json:"customer_id"`
	Status      string     `json:"status"`
	AttemptedAt *time.Time `json:"attempted_at"`
	CompletedAt *time.Time `json:"completed_at"`
}

type RefreshToken struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"user_id"`
	TokenHash string     `json:"token_hash"`
	ExpiresAt time.Time  `json:"expires_at"`
	RevokedAt *time.Time `json:"revoked_at"`
	CreatedAt time.Time  `json:"created_at"`
}

type ListCallsRow struct {
	Call
	CustomerName  *string `json:"customer_name"`
	CustomerPhone *string `json:"customer_phone"`
}

type DailySummaryRow struct {
	AgentID         *uuid.UUID `json:"agent_id"`
	TotalCalls      int64      `json:"total_calls"`
	CompletedCalls  int64      `json:"completed_calls"`
	AbandonedCalls  int64      `json:"abandoned_calls"`
	AvgTalkSeconds  *float64   `json:"avg_talk_seconds"`
	AvgWaitSeconds  *float64   `json:"avg_wait_seconds"`
	TotalCostCents  *int64     `json:"total_cost_cents"`
}

func scanUser(row pgx.Row) (User, error) {
	var u User
	err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.FullName, &u.Role, &u.CreatedAt, &u.DeletedAt)
	return u, err
}

func scanAgent(row pgx.Row) (Agent, error) {
	var a Agent
	err := row.Scan(&a.ID, &a.UserID, &a.Extension, &a.Skills, &a.Status, &a.MaxConcurrent, &a.TeamID, &a.LastSeenAt)
	return a, err
}

func scanCustomer(row pgx.Row) (Customer, error) {
	var c Customer
	err := row.Scan(&c.ID, &c.PhoneNumber, &c.FullName, &c.Email, &c.Tags, &c.DncListed, &c.Timezone, &c.CreatedAt, &c.UpdatedAt)
	return c, err
}

func scanCall(row pgx.Row) (Call, error) {
	var c Call
	err := row.Scan(
		&c.ID,
		&c.ProviderSid,
		&c.CustomerID,
		&c.AgentID,
		&c.QueueID,
		&c.Direction,
		&c.Status,
		&c.StartedAt,
		&c.AnsweredAt,
		&c.EndedAt,
		&c.WaitSeconds,
		&c.TalkSeconds,
		&c.RecordingUrl,
		&c.CostCents,
		&c.FromNumber,
		&c.ToNumber,
	)
	return c, err
}

func toJSON(v any) []byte {
	b, _ := json.Marshal(v)
	return b
}
