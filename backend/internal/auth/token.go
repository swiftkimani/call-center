package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/yourorg/callcenter/internal/config"
)

const accessTokenTTL = 15 * time.Minute

var ErrInvalidToken = errors.New("invalid or expired token")

type Claims struct {
	jwt.RegisteredClaims
	UserID  uuid.UUID `json:"uid"`
	AgentID *uuid.UUID `json:"aid,omitempty"`
	Role    string    `json:"role"`
}

func IssueAccessToken(cfg *config.Config, userID uuid.UUID, agentID *uuid.UUID, role string) (string, error) {
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID:  userID,
		AgentID: agentID,
		Role:    role,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}
	return signed, nil
}

func VerifyAccessToken(cfg *config.Config, tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(cfg.JWTSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
