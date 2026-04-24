package queues

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

var ErrQueueEmpty = errors.New("queue is empty")
var ErrNoQueuesConfigured = errors.New("no queues configured")

type Service struct {
	rdb *goredis.Client
	q   *queries.Queries
}

type Snapshot struct {
	QueueID           uuid.UUID `json:"queue_id"`
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	SkillsRequired    []string  `json:"skills_required"`
	MaxWaitSeconds    int32     `json:"max_wait_seconds"`
	SLASeconds        int32     `json:"sla_seconds"`
	Waiting           int64     `json:"waiting"`
	OldestWaitSeconds int64     `json:"oldest_wait_seconds"`
	SLABreachCount    int64     `json:"sla_breach_count"`
}

func NewService(pool *pgxpool.Pool, rdb *goredis.Client) *Service {
	return &Service{rdb: rdb, q: queries.New(pool)}
}

func (s *Service) Enqueue(ctx context.Context, queueID, callID uuid.UUID) error {
	score := float64(time.Now().UnixMilli())
	key := rkeys.QueueWaitKey(queueID.String())
	return s.rdb.ZAdd(ctx, key, goredis.Z{
		Score:  score,
		Member: callID.String(),
	}).Err()
}

func (s *Service) DequeueOldest(ctx context.Context, queueID uuid.UUID) (uuid.UUID, int, error) {
	key := rkeys.QueueWaitKey(queueID.String())
	results, err := s.rdb.ZPopMin(ctx, key, 1).Result()
	if err != nil {
		return uuid.Nil, 0, fmt.Errorf("ZPopMin: %w", err)
	}
	if len(results) == 0 {
		return uuid.Nil, 0, ErrQueueEmpty
	}

	callID, err := uuid.Parse(results[0].Member.(string))
	if err != nil {
		return uuid.Nil, 0, fmt.Errorf("invalid call id: %w", err)
	}

	enqueuedAtMs := int64(results[0].Score)
	waitMs := time.Now().UnixMilli() - enqueuedAtMs
	waitSecs := int(waitMs / 1000)

	return callID, waitSecs, nil
}

func (s *Service) Depth(ctx context.Context, queueID uuid.UUID) (int64, error) {
	return s.rdb.ZCard(ctx, rkeys.QueueWaitKey(queueID.String())).Result()
}

func (s *Service) Remove(ctx context.Context, queueID, callID uuid.UUID) error {
	return s.rdb.ZRem(ctx, rkeys.QueueWaitKey(queueID.String()), callID.String()).Err()
}

func (s *Service) Snapshot(ctx context.Context, queueID uuid.UUID) (*Snapshot, error) {
	queue, err := s.q.GetQueueByID(ctx, queueID)
	if err != nil {
		return nil, fmt.Errorf("GetQueueByID: %w", err)
	}

	key := rkeys.QueueWaitKey(queueID.String())
	nowMs := time.Now().UnixMilli()

	pipe := s.rdb.Pipeline()
	depthCmd := pipe.ZCard(ctx, key)
	oldestCmd := pipe.ZRangeWithScores(ctx, key, 0, 0)
	breachCmd := pipe.ZCount(ctx, key, "-inf", fmt.Sprintf("%d", nowMs-int64(queue.SlaSeconds)*1000))
	if _, err := pipe.Exec(ctx); err != nil {
		return nil, fmt.Errorf("redis pipeline: %w", err)
	}

	oldestWaitSeconds := int64(0)
	if oldest, err := oldestCmd.Result(); err == nil && len(oldest) > 0 {
		enqueuedAtMs := int64(oldest[0].Score)
		if nowMs > enqueuedAtMs {
			oldestWaitSeconds = (nowMs - enqueuedAtMs) / 1000
		}
	}

	waiting, err := depthCmd.Result()
	if err != nil {
		return nil, fmt.Errorf("ZCard: %w", err)
	}
	slaBreaches, err := breachCmd.Result()
	if err != nil {
		return nil, fmt.Errorf("ZCount: %w", err)
	}

	return &Snapshot{
		QueueID:           queue.ID,
		Name:              queue.Name,
		Description:       queue.Description,
		SkillsRequired:    queue.SkillsRequired,
		MaxWaitSeconds:    queue.MaxWaitSeconds,
		SLASeconds:        queue.SlaSeconds,
		Waiting:           waiting,
		OldestWaitSeconds: oldestWaitSeconds,
		SLABreachCount:    slaBreaches,
	}, nil
}

func (s *Service) DefaultQueue(ctx context.Context) (*queries.Queue, error) {
	queues, err := s.q.ListQueues(ctx)
	if err != nil {
		return nil, fmt.Errorf("ListQueues: %w", err)
	}
	if len(queues) == 0 {
		return nil, ErrNoQueuesConfigured
	}
	queue := queues[0]
	return &queue, nil
}
