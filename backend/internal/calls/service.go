package calls

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	goredis "github.com/redis/go-redis/v9"
	"github.com/yourorg/callcenter/internal/db/queries"
	rkeys "github.com/yourorg/callcenter/internal/redis"
	"github.com/yourorg/callcenter/internal/telephony"
)

var ErrNotFound = errors.New("call not found")
var ErrDNCListed = errors.New("customer is on do-not-call list")
var ErrOutsideDialingHours = errors.New("outside permitted dialing hours")

type Service struct {
	q               *queries.Queries
	rdb             *goredis.Client
	telephony       telephony.Adapter
	dialingHourStart int
	dialingHourEnd   int
}

func NewService(pool *pgxpool.Pool, rdb *goredis.Client, adapter telephony.Adapter, hourStart, hourEnd int) *Service {
	return &Service{
		q:               queries.New(pool),
		rdb:             rdb,
		telephony:       adapter,
		dialingHourStart: hourStart,
		dialingHourEnd:   hourEnd,
	}
}

type CreateInboundParams struct {
	ProviderSID string
	FromNumber  string
	ToNumber    string
	QueueID     *uuid.UUID
	CustomerID  *uuid.UUID
}

func (s *Service) CreateInbound(ctx context.Context, p CreateInboundParams) (*queries.Call, error) {
	call, err := s.q.CreateCall(ctx, queries.CreateCallParams{
		ProviderSid: p.ProviderSID,
		CustomerID:  p.CustomerID,
		AgentID:     nil,
		QueueID:     p.QueueID,
		Direction:   "inbound",
		Status:      "queued",
		FromNumber:  p.FromNumber,
		ToNumber:    p.ToNumber,
	})
	if err != nil {
		return nil, fmt.Errorf("CreateCall: %w", err)
	}

	// Cache call state in Redis
	s.cacheCall(ctx, &call)

	s.appendEvent(ctx, call.ID, "created", map[string]any{
		"from": p.FromNumber, "to": p.ToNumber,
	})
	return &call, nil
}

func (s *Service) Answer(ctx context.Context, callID, agentID uuid.UUID) error {
	if err := s.q.AnswerCall(ctx, queries.AnswerCallParams{
		ID:      callID,
		AgentID: &agentID,
	}); err != nil {
		return fmt.Errorf("AnswerCall: %w", err)
	}
	s.rdb.HSet(ctx, rkeys.CallKey(callID.String()), "status", "in_progress", "agent_id", agentID.String())
	s.appendEvent(ctx, callID, "answered", map[string]any{"agent_id": agentID})
	return nil
}

func (s *Service) End(ctx context.Context, callID uuid.UUID, status string, costCents *int32) error {
	if err := s.q.EndCall(ctx, queries.EndCallParams{
		ID:        callID,
		Status:    status,
		CostCents: costCents,
	}); err != nil {
		return fmt.Errorf("EndCall: %w", err)
	}
	s.rdb.Del(ctx, rkeys.CallKey(callID.String()))
	s.appendEvent(ctx, callID, "ended", map[string]any{"status": status})
	return nil
}

func (s *Service) SaveDisposition(ctx context.Context, callID, agentID uuid.UUID, category, notes string) error {
	_, err := s.q.SaveDisposition(ctx, queries.SaveDispositionParams{
		CallID:   callID,
		AgentID:  agentID,
		Category: category,
		Notes:    notes,
	})
	return err
}

func (s *Service) GetByProviderSID(ctx context.Context, sid string) (*queries.Call, error) {
	call, err := s.q.GetCallByProviderSID(ctx, sid)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("GetCallByProviderSID: %w", err)
	}
	return &call, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*queries.Call, error) {
	call, err := s.q.GetCallByID(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("GetCallByID: %w", err)
	}
	return &call, nil
}

func (s *Service) UpdateRecordingURL(ctx context.Context, callID uuid.UUID, url string) error {
	return s.q.UpdateCallRecording(ctx, queries.UpdateCallRecordingParams{
		ID:           callID,
		RecordingUrl: &url,
	})
}

func (s *Service) List(ctx context.Context, limit, offset int) ([]queries.ListCallsRow, error) {
	return s.q.ListCalls(ctx, queries.ListCallsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
}

func (s *Service) InitiateOutbound(ctx context.Context, agentID, customerID uuid.UUID, customerPhone string) (*queries.Call, error) {
	// Compliance: check DNC
	customer, err := s.q.GetCustomerByID(ctx, customerID)
	if err == nil && customer.DncListed {
		return nil, ErrDNCListed
	}

	// Compliance: check dialing hours (use server local time; adjust per customer timezone in v2)
	hour := time.Now().Hour()
	if hour < s.dialingHourStart || hour >= s.dialingHourEnd {
		return nil, ErrOutsideDialingHours
	}

	call, err := s.q.CreateCall(ctx, queries.CreateCallParams{
		ProviderSid: fmt.Sprintf("pending-%s-%d", customerID, time.Now().UnixNano()),
		CustomerID:  &customerID,
		AgentID:     &agentID,
		Direction:   "outbound",
		Status:      "ringing",
		FromNumber:  "system",
		ToNumber:    customerPhone,
	})
	if err != nil {
		return nil, fmt.Errorf("CreateCall: %w", err)
	}

	s.appendEvent(ctx, call.ID, "outbound_initiated", map[string]any{
		"agent_id": agentID, "customer_phone": customerPhone,
	})
	return &call, nil
}

func (s *Service) cacheCall(ctx context.Context, call *queries.Call) {
	key := rkeys.CallKey(call.ID.String())
	s.rdb.HSet(ctx, key,
		"status", call.Status,
		"from", call.FromNumber,
		"to", call.ToNumber,
	)
	s.rdb.Expire(ctx, key, 4*time.Hour)
}

func (s *Service) appendEvent(ctx context.Context, callID uuid.UUID, eventType string, payload map[string]any) {
	b, _ := json.Marshal(payload)
	s.q.InsertCallEvent(ctx, queries.InsertCallEventParams{
		CallID:    callID,
		EventType: eventType,
		Payload:   b,
	})
}
