package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/yourorg/callcenter/internal/hub"
	"github.com/yourorg/callcenter/internal/queues"
)

type SupervisorHandler struct {
	hub      *hub.Hub
	queueSvc *queues.Service
}

func NewSupervisorHandler(h *hub.Hub, queueSvc *queues.Service) *SupervisorHandler {
	return &SupervisorHandler{hub: h, queueSvc: queueSvc}
}

func (h *SupervisorHandler) QueueLive(w http.ResponseWriter, r *http.Request) {
	queueID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		fail(w, http.StatusBadRequest, "invalid queue id")
		return
	}

	snapshot, err := h.queueSvc.Snapshot(r.Context(), queueID)
	if err != nil {
		fail(w, http.StatusInternalServerError, "could not fetch queue snapshot")
		return
	}

	ok(w, snapshot)
}

func (h *SupervisorHandler) Whisper(w http.ResponseWriter, r *http.Request) {
	callID := chi.URLParam(r, "call_id")
	// In a real implementation this would issue a provider API call to add
	// the supervisor as a muted conference participant.
	h.hub.BroadcastSupervisors(hub.NewEvent(hub.EventWhisperStarted, map[string]string{
		"call_id": callID,
	}))
	ok(w, nil)
}

func (h *SupervisorHandler) Barge(w http.ResponseWriter, r *http.Request) {
	callID := chi.URLParam(r, "call_id")
	_ = callID
	// Provider barge implementation goes here.
	ok(w, map[string]string{"status": "not_implemented"})
}
