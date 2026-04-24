package handlers

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourorg/callcenter/internal/auth"
	"github.com/yourorg/callcenter/internal/db/queries"
)

type MeHandler struct {
	q *queries.Queries
}

func NewMeHandler(pool *pgxpool.Pool) *MeHandler {
	return &MeHandler{q: queries.New(pool)}
}

func (h *MeHandler) Get(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromContext(r.Context())
	if claims == nil {
		fail(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	user, err := h.q.GetUserByID(r.Context(), claims.UserID)
	if err != nil {
		fail(w, http.StatusNotFound, "user not found")
		return
	}

	result := map[string]any{
		"id":        user.ID,
		"email":     user.Email,
		"full_name": user.FullName,
		"role":      user.Role,
	}

	if claims.AgentID != nil {
		agent, err := h.q.GetAgentByUserID(r.Context(), user.ID)
		if err == nil {
			result["agent"] = map[string]any{
				"id":        agent.ID,
				"extension": agent.Extension,
				"skills":    agent.Skills,
				"status":    agent.Status,
			}
		}
	}

	ok(w, result)
}
