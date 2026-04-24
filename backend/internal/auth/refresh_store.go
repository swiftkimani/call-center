package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourorg/callcenter/internal/db/queries"
)

const refreshTokenTTL = 7 * 24 * time.Hour

var ErrRefreshTokenInvalid = errors.New("refresh token invalid or expired")

type RefreshStore struct {
	q *queries.Queries
}

func NewRefreshStore(pool *pgxpool.Pool) *RefreshStore {
	return &RefreshStore{q: queries.New(pool)}
}

func (s *RefreshStore) Issue(ctx context.Context, userID uuid.UUID) (string, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", fmt.Errorf("rand.Read: %w", err)
	}
	token := base64.URLEncoding.EncodeToString(raw)
	hash := hashToken(token)

	if err := s.q.CreateRefreshToken(ctx, queries.CreateRefreshTokenParams{
		UserID:    userID,
		TokenHash: hash,
		ExpiresAt: time.Now().Add(refreshTokenTTL),
	}); err != nil {
		return "", fmt.Errorf("CreateRefreshToken: %w", err)
	}
	return token, nil
}

func (s *RefreshStore) Rotate(ctx context.Context, token string) (uuid.UUID, string, error) {
	hash := hashToken(token)
	row, err := s.q.GetRefreshToken(ctx, hash)
	if err != nil {
		return uuid.Nil, "", ErrRefreshTokenInvalid
	}

	if err := s.q.RevokeRefreshToken(ctx, hash); err != nil {
		return uuid.Nil, "", fmt.Errorf("revoke old token: %w", err)
	}

	newToken, err := s.Issue(ctx, row.UserID)
	if err != nil {
		return uuid.Nil, "", err
	}
	return row.UserID, newToken, nil
}

func (s *RefreshStore) Revoke(ctx context.Context, token string) error {
	return s.q.RevokeRefreshToken(ctx, hashToken(token))
}

func (s *RefreshStore) RevokeAll(ctx context.Context, userID uuid.UUID) error {
	return s.q.RevokeAllUserRefreshTokens(ctx, userID)
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}
