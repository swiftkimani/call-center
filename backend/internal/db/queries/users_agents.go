package queries

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type UpdateAgentStatusParams struct {
	ID     uuid.UUID
	Status string
}

type CreateRefreshTokenParams struct {
	UserID    uuid.UUID
	TokenHash string
	ExpiresAt time.Time
}

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	return scanUser(q.db.QueryRow(ctx, `SELECT id, email, password_hash, full_name, role, created_at, deleted_at FROM users WHERE email = $1 AND deleted_at IS NULL`, email))
}

func (q *Queries) GetUserByID(ctx context.Context, id uuid.UUID) (User, error) {
	return scanUser(q.db.QueryRow(ctx, `SELECT id, email, password_hash, full_name, role, created_at, deleted_at FROM users WHERE id = $1 AND deleted_at IS NULL`, id))
}

func (q *Queries) GetAgentByID(ctx context.Context, id uuid.UUID) (Agent, error) {
	return scanAgent(q.db.QueryRow(ctx, `SELECT id, user_id, extension, skills, status, max_concurrent, team_id, last_seen_at FROM agents WHERE id = $1`, id))
}

func (q *Queries) GetAgentByUserID(ctx context.Context, userID uuid.UUID) (Agent, error) {
	return scanAgent(q.db.QueryRow(ctx, `SELECT id, user_id, extension, skills, status, max_concurrent, team_id, last_seen_at FROM agents WHERE user_id = $1`, userID))
}

func (q *Queries) UpdateAgentStatus(ctx context.Context, arg UpdateAgentStatusParams) error {
	_, err := q.db.Exec(ctx, `UPDATE agents SET status = $2, last_seen_at = NOW() WHERE id = $1`, arg.ID, arg.Status)
	return err
}

func (q *Queries) ListAvailableAgents(ctx context.Context) ([]Agent, error) {
	rows, err := q.db.Query(ctx, `SELECT id, user_id, extension, skills, status, max_concurrent, team_id, last_seen_at FROM agents WHERE status = 'available' ORDER BY last_seen_at ASC NULLS LAST`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Agent
	for rows.Next() {
		item, err := scanAgent(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (q *Queries) CreateRefreshToken(ctx context.Context, arg CreateRefreshTokenParams) error {
	_, err := q.db.Exec(ctx, `INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3)`, arg.UserID, arg.TokenHash, arg.ExpiresAt)
	return err
}

func (q *Queries) GetRefreshToken(ctx context.Context, tokenHash string) (RefreshToken, error) {
	var r RefreshToken
	err := q.db.QueryRow(ctx, `SELECT id, user_id, token_hash, expires_at, revoked_at, created_at FROM refresh_tokens WHERE token_hash = $1 AND revoked_at IS NULL AND expires_at > NOW()`, tokenHash).
		Scan(&r.ID, &r.UserID, &r.TokenHash, &r.ExpiresAt, &r.RevokedAt, &r.CreatedAt)
	return r, err
}

func (q *Queries) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	_, err := q.db.Exec(ctx, `UPDATE refresh_tokens SET revoked_at = NOW() WHERE token_hash = $1`, tokenHash)
	return err
}

func (q *Queries) RevokeAllUserRefreshTokens(ctx context.Context, userID uuid.UUID) error {
	_, err := q.db.Exec(ctx, `UPDATE refresh_tokens SET revoked_at = NOW() WHERE user_id = $1 AND revoked_at IS NULL`, userID)
	return err
}
