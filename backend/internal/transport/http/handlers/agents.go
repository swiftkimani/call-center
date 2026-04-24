package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/yourorg/callcenter/internal/agents"
	"github.com/yourorg/callcenter/internal/auth"
)

type AgentsHandler struct {
	svc *agents.Service
}

func NewAgentsHandler(svc *agents.Service) *AgentsHandler {
	return &AgentsHandler{svc: svc}
}

func (h *AgentsHandler) SetStatus(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromContext(r.Context())

	// Agents may only update their own status; admins/supervisors can update anyone's.
	targetID := chi.URLParam(r, "id")
	if claims.Role == "agent" {
		if claims.AgentID == nil || claims.AgentID.String() != targetID {
			fail(w, http.StatusForbidden, "forbidden")
			return
		}
	}

	agentID, err := uuid.Parse(targetID)
	if err != nil {
		fail(w, http.StatusBadRequest, "invalid agent id")
		return
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Status == "" {
		fail(w, http.StatusBadRequest, "status required")
		return
	}

	if err := h.svc.SetStatus(r.Context(), agentID, req.Status); err != nil {
		fail(w, http.StatusBadRequest, err.Error())
		return
	}

	ok(w, map[string]string{"status": req.Status})
}

func (h *AgentsHandler) Heartbeat(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromContext(r.Context())
	if claims.AgentID == nil {
		fail(w, http.StatusBadRequest, "not an agent")
		return
	}
	h.svc.Heartbeat(r.Context(), *claims.AgentID)
	w.WriteHeader(http.StatusNoContent)
}
