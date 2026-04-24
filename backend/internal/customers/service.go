package customers

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourorg/callcenter/internal/db/queries"
)

var ErrNotFound = errors.New("customer not found")

type Service struct {
	q *queries.Queries
}

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{q: queries.New(pool)}
}

func (s *Service) FindByPhone(ctx context.Context, phone string) (*queries.Customer, error) {
	c, err := s.q.GetCustomerByPhone(ctx, phone)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("GetCustomerByPhone: %w", err)
	}
	return &c, nil
}

func (s *Service) FindOrCreate(ctx context.Context, phone, name string) (*queries.Customer, error) {
	c, err := s.FindByPhone(ctx, phone)
	if err == nil {
		return c, nil
	}
	if !errors.Is(err, ErrNotFound) {
		return nil, err
	}
	created, err := s.q.CreateCustomer(ctx, queries.CreateCustomerParams{
		PhoneNumber: phone,
		FullName:    name,
		Tags:        []string{},
	})
	if err != nil {
		return nil, fmt.Errorf("CreateCustomer: %w", err)
	}
	return &created, nil
}

func (s *Service) Get(ctx context.Context, id uuid.UUID) (*queries.Customer, error) {
	c, err := s.q.GetCustomerByID(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("GetCustomerByID: %w", err)
	}
	return &c, nil
}

func (s *Service) Search(ctx context.Context, query string, limit, offset int) ([]queries.Customer, error) {
	return s.q.SearchCustomers(ctx, queries.SearchCustomersParams{
		Column1: query,
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, name, email string, tags []string) (*queries.Customer, error) {
	var emailPtr *string
	if email != "" {
		emailPtr = &email
	}
	c, err := s.q.UpdateCustomer(ctx, queries.UpdateCustomerParams{
		ID:       id,
		FullName: name,
		Email:    emailPtr,
		Tags:     tags,
	})
	if err != nil {
		return nil, fmt.Errorf("UpdateCustomer: %w", err)
	}
	return &c, nil
}
