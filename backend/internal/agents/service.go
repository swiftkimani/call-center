package agents

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	goredis "github.com/redis/go-redis/v9"
	"github.com/yourorg/callcenter/internal/db/queries"
	rkeys "github.com/yourorg/callcenter/internal/redis"
)

var (
	ErrInvalidTransition = errors.New("invalid status transition")
	ErrAgentNotFound     = errors.New("agent not found")
)

// validTransitions maps current status → allowed next statuses.
var validTransitions = map[string][]string{
	"offline":   {"available"},
	"available": {"busy", "break", "offline"},
	"busy":      {"wrap_up", "available"},
	"wrap_up":   {"available", "break", "offline"},
	"break":     {"available", "offline"},
}

const heartbeatTTL = 30 * time.Second

type Service struct {
	q   *queries.Queries
	rdb *goredis.Client
}

func NewService(pool *pgxpool.Pool, rdb *goredis.Client) *Service {
	return &Service{q: queries.New(pool), rdb: rdb}
}

func (s *Service) SetStatus(ctx context.Context, agentID uuid.UUID, newStatus string) error {
	current, err := s.GetPresence(ctx, agentID)
	if err != nil {
		// Agent may not have a Redis key yet (first login)
		current = "offline"
	}

	if !isValidTransition(current, newStatus) {
		return fmt.Errorf("%w: %s → %s", ErrInvalidTransition, current, newStatus)
	}

	pipe := s.rdb.Pipeline()
	agentKey := rkeys.AgentKey(agentID.String())

	pipe.HSet(ctx, agentKey, "status", newStatus, "last_seen", time.Now().UTC().Format(time.RFC3339))

	// Maintain the available/busy sets
	for _, st := range []string{"available", "busy", "wrap_up", "break"} {
		pipe.SRem(ctx, rkeys.AgentSetKey(st), agentID.String())
	}
	if newStatus != "offline" {
		pipe.SAdd(ctx, rkeys.AgentSetKey(newStatus), agentID.String())
	}

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("redis pipeline: %w", err)
	}

	return s.q.UpdateAgentStatus(ctx, queries.UpdateAgentStatusParams{
		ID:     agentID,
		Status: newStatus,
	})
}

func (s *Service) Heartbeat(ctx context.Context, agentID uuid.UUID) error {
	key := rkeys.AgentHeartbeatKey(agentID.String())
	return s.rdb.Set(ctx, key, "1", heartbeatTTL).Err()
}

func (s *Service) GetPresence(ctx context.Context, agentID uuid.UUID) (string, error) {
	status, err := s.rdb.HGet(ctx, rkeys.AgentKey(agentID.String()), "status").Result()
	if err != nil {
		return "", fmt.Errorf("redis HGet: %w", err)
	}
	return status, nil
}

func (s *Service) ListAvailable(ctx context.Context) ([]queries.Agent, error) {
	members, err := s.rdb.SMembers(ctx, rkeys.AgentSetKey("available")).Result()
	if err != nil || len(members) == 0 {
		// Fall back to Postgres
		return s.q.ListAvailableAgents(ctx)
	}
	agents := make([]queries.Agent, 0, len(members))
	for _, id := range members {
		uid, err := uuid.Parse(id)
		if err != nil {
			continue
		}
		row, err := s.q.GetAgentByID(ctx, uid)
		if err != nil {
			continue
		}
		agents = append(agents, queries.Agent{
			ID:          row.ID,
			UserID:      row.UserID,
			Extension:   row.Extension,
			Skills:      row.Skills,
			Status:      row.Status,
			MaxConcurrent: row.MaxConcurrent,
			TeamID:      row.TeamID,
			LastSeenAt:  row.LastSeenAt,
		})
	}
	return agents, nil
}

func isValidTransition(from, to string) bool {
	allowed, ok := validTransitions[from]
	if !ok {
		return false
	}
	for _, s := range allowed {
		if s == to {
			return true
		}
	}
	return false
}
