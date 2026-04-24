package queues

import (
	"context"
	"errors"
	"sort"
	"time"

	"github.com/google/uuid"
	goredis "github.com/redis/go-redis/v9"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourorg/callcenter/internal/db/queries"
	rkeys "github.com/yourorg/callcenter/internal/redis"
)

var ErrNoAgentAvailable = errors.New("no agent available")

type Router struct {
	q   *queries.Queries
	rdb *goredis.Client
}

func NewRouter(pool *pgxpool.Pool, rdb *goredis.Client) *Router {
	return &Router{q: queries.New(pool), rdb: rdb}
}

// SelectAgent picks the longest-idle available agent that satisfies the required skills.
// If requiredSkills is empty, any available agent qualifies.
func (r *Router) SelectAgent(ctx context.Context, requiredSkills []string) (*queries.Agent, error) {
	members, err := r.rdb.SMembers(ctx, rkeys.AgentSetKey("available")).Result()
	if err != nil || len(members) == 0 {
		return r.selectAgentFromDB(ctx, requiredSkills)
	}

	type candidate struct {
		agent      queries.Agent
		lastSeenAt time.Time
	}

	var candidates []candidate
	for _, id := range members {
		uid, err := uuid.Parse(id)
		if err != nil {
			continue
		}
		row, err := r.q.GetAgentByID(ctx, uid)
		if err != nil {
			continue
		}
		if row.Status != "available" {
			continue
		}
		if !hasRequiredSkills(row.Skills, requiredSkills) {
			continue
		}
		ts := time.Time{}
		if row.LastSeenAt != nil {
			ts = *row.LastSeenAt
		}
		candidates = append(candidates, candidate{
			agent:      queries.Agent{
				ID:            row.ID,
				UserID:        row.UserID,
				Extension:     row.Extension,
				Skills:        row.Skills,
				Status:        row.Status,
				MaxConcurrent: row.MaxConcurrent,
				TeamID:        row.TeamID,
				LastSeenAt:    row.LastSeenAt,
			},
			lastSeenAt: ts,
		})
	}

	if len(candidates) == 0 {
		return nil, ErrNoAgentAvailable
	}

	// Oldest last_seen = longest idle
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].lastSeenAt.Before(candidates[j].lastSeenAt)
	})
	selected := candidates[0].agent
	return &selected, nil
}

func (r *Router) selectAgentFromDB(ctx context.Context, requiredSkills []string) (*queries.Agent, error) {
	agents, err := r.q.ListAvailableAgents(ctx)
	if err != nil {
		return nil, err
	}
	for _, a := range agents {
		if hasRequiredSkills(a.Skills, requiredSkills) {
			ag := a
			return &ag, nil
		}
	}
	return nil, ErrNoAgentAvailable
}

func hasRequiredSkills(agentSkills, required []string) bool {
	if len(required) == 0 {
		return true
	}
	skillSet := make(map[string]struct{}, len(agentSkills))
	for _, s := range agentSkills {
		skillSet[s] = struct{}{}
	}
	for _, r := range required {
		if _, ok := skillSet[r]; !ok {
			return false
		}
	}
	return true
}
