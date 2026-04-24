package campaigns

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourorg/callcenter/internal/db/queries"
)

type ContactRow struct {
	CustomerID uuid.UUID
}

type Service struct {
	q *queries.Queries
}

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{q: queries.New(pool)}
}

func (s *Service) Create(ctx context.Context, name string, scheduledAt *time.Time, createdBy uuid.UUID) (*queries.Campaign, error) {
	c, err := s.q.CreateCampaign(ctx, queries.CreateCampaignParams{
		Name:        name,
		ScheduledAt: scheduledAt,
		CreatedBy:   createdBy,
	})
	if err != nil {
		return nil, fmt.Errorf("CreateCampaign: %w", err)
	}
	return &c, nil
}

func (s *Service) ImportContacts(ctx context.Context, campaignID uuid.UUID, contacts []ContactRow) error {
	for _, c := range contacts {
		if _, err := s.q.InsertCampaignContact(ctx, queries.InsertCampaignContactParams{
			CampaignID: campaignID,
			CustomerID: c.CustomerID,
		}); err != nil {
			return fmt.Errorf("InsertCampaignContact: %w", err)
		}
	}
	return nil
}

func (s *Service) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	return s.q.UpdateCampaignStatus(ctx, queries.UpdateCampaignStatusParams{
		ID:     id,
		Status: status,
	})
}

func (s *Service) Get(ctx context.Context, id uuid.UUID) (*queries.Campaign, error) {
	c, err := s.q.GetCampaignByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("GetCampaignByID: %w", err)
	}
	return &c, nil
}

func (s *Service) List(ctx context.Context, limit, offset int) ([]queries.Campaign, error) {
	return s.q.ListCampaigns(ctx, queries.ListCampaignsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
}
