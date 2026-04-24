package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourorg/callcenter/internal/auth"
	"github.com/yourorg/callcenter/internal/config"
	"github.com/yourorg/callcenter/internal/db/queries"
)

type AuthHandler struct {
	cfg          *config.Config
	q            *queries.Queries
	refreshStore *auth.RefreshStore
}

func NewAuthHandler(cfg *config.Config, pool *pgxpool.Pool) *AuthHandler {
	return &AuthHandler{
		cfg:          cfg,
		q:            queries.New(pool),
		refreshStore: auth.NewRefreshStore(pool),
	}
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Role         string `json:"role"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fail(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Email == "" || req.Password == "" {
		fail(w, http.StatusBadRequest, "email and password required")
		return
	}

	user, err := h.q.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		fail(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	match, err := auth.VerifyPassword(user.PasswordHash, req.Password)
	if err != nil || !match {
		fail(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	var agentID *uuid.UUID
	if agent, err := h.q.GetAgentByUserID(r.Context(), user.ID); err == nil {
		agentID = &agent.ID
	}

	accessToken, err := auth.IssueAccessToken(h.cfg, user.ID, agentID, user.Role)
	if err != nil {
		fail(w, http.StatusInternalServerError, "could not issue token")
		return
	}

	refreshToken, err := h.refreshStore.Issue(r.Context(), user.ID)
	if err != nil {
		fail(w, http.StatusInternalServerError, "could not issue refresh token")
		return
	}

	ok(w, tokenResponse{AccessToken: accessToken, RefreshToken: refreshToken, Role: user.Role})
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.RefreshToken == "" {
		fail(w, http.StatusBadRequest, "refresh_token required")
		return
	}

	userID, newRefresh, err := h.refreshStore.Rotate(r.Context(), req.RefreshToken)
	if err != nil {
		fail(w, http.StatusUnauthorized, "invalid or expired refresh token")
		return
	}

	user, err := h.q.GetUserByID(r.Context(), userID)
	if err != nil {
		fail(w, http.StatusUnauthorized, "user not found")
		return
	}

	accessToken, err := auth.IssueAccessToken(h.cfg, user.ID, nil, user.Role)
	if err != nil {
		fail(w, http.StatusInternalServerError, "could not issue token")
		return
	}

	ok(w, tokenResponse{AccessToken: accessToken, RefreshToken: newRefresh, Role: user.Role})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err == nil && req.RefreshToken != "" {
		h.refreshStore.Revoke(r.Context(), req.RefreshToken)
	}
	ok(w, nil)
}
