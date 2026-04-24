package reports

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourorg/callcenter/internal/db/queries"
)

type AgentSummary struct {
	AgentID        string  `json:"agent_id"`
	TotalCalls     int64   `json:"total_calls"`
	CompletedCalls int64   `json:"completed_calls"`
	AbandonedCalls int64   `json:"abandoned_calls"`
	AvgTalkSeconds float64 `json:"avg_talk_seconds"`
	AvgWaitSeconds float64 `json:"avg_wait_seconds"`
	TotalCostCents int64   `json:"total_cost_cents"`
}

type DailySummary struct {
	Date   string         `json:"date"`
	Agents []AgentSummary `json:"agents"`
}

type Service struct {
	q *queries.Queries
}

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{q: queries.New(pool)}
}

func (s *Service) Daily(ctx context.Context, date time.Time) (*DailySummary, error) {
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	end := start.Add(24 * time.Hour)

	rows, err := s.q.DailySummary(ctx, queries.DailySummaryParams{
		StartedAt:   start,
		StartedAt_2: end,
	})
	if err != nil {
		return nil, fmt.Errorf("DailySummary: %w", err)
	}

	summary := &DailySummary{Date: date.Format("2006-01-02")}
	for _, row := range rows {
		agentID := ""
		if row.AgentID != nil {
			agentID = row.AgentID.String()
		}
		as := AgentSummary{
			AgentID:        agentID,
			TotalCalls:     row.TotalCalls,
			CompletedCalls: row.CompletedCalls,
			AbandonedCalls: row.AbandonedCalls,
		}
		if row.AvgTalkSeconds != nil {
			as.AvgTalkSeconds = *row.AvgTalkSeconds
		}
		if row.AvgWaitSeconds != nil {
			as.AvgWaitSeconds = *row.AvgWaitSeconds
		}
		if row.TotalCostCents != nil {
			as.TotalCostCents = *row.TotalCostCents
		}
		summary.Agents = append(summary.Agents, as)
	}
	return summary, nil
}

func (s *Service) ExportCSV(ctx context.Context, w io.Writer, date time.Time) error {
	summary, err := s.Daily(ctx, date)
	if err != nil {
		return err
	}
	cw := csv.NewWriter(w)
	cw.Write([]string{"agent_id", "total_calls", "completed_calls", "abandoned_calls", "avg_talk_seconds", "avg_wait_seconds", "total_cost_cents"})
	for _, a := range summary.Agents {
		cw.Write([]string{
			a.AgentID,
			fmt.Sprint(a.TotalCalls),
			fmt.Sprint(a.CompletedCalls),
			fmt.Sprint(a.AbandonedCalls),
			fmt.Sprintf("%.1f", a.AvgTalkSeconds),
			fmt.Sprintf("%.1f", a.AvgWaitSeconds),
			fmt.Sprint(a.TotalCostCents),
		})
	}
	cw.Flush()
	return cw.Error()
}
