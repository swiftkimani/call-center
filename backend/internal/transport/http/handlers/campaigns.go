package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/yourorg/callcenter/internal/auth"
	"github.com/yourorg/callcenter/internal/campaigns"
)

type CampaignsHandler struct {
	svc *campaigns.Service
}

func NewCampaignsHandler(svc *campaigns.Service) *CampaignsHandler {
	return &CampaignsHandler{svc: svc}
}

func (h *CampaignsHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromContext(r.Context())
	if claims == nil {
		fail(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req struct {
		Name        string `json:"name"`
		ScheduledAt string `json:"scheduled_at"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fail(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" {
		fail(w, http.StatusBadRequest, "name required")
		return
	}

	var scheduledAt *time.Time
	if req.ScheduledAt != "" {
		t, err := time.Parse(time.RFC3339, req.ScheduledAt)
		if err != nil {
			fail(w, http.StatusBadRequest, "scheduled_at must be RFC3339")
			return
		}
		scheduledAt = &t
	}

	campaign, err := h.svc.Create(r.Context(), req.Name, scheduledAt, claims.UserID)
	if err != nil {
		fail(w, http.StatusInternalServerError, "could not create campaign")
		return
	}
	ok(w, campaign)
}

func (h *CampaignsHandler) List(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 25
	}
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	items, err := h.svc.List(r.Context(), limit, offset)
	if err != nil {
		fail(w, http.StatusInternalServerError, "could not list campaigns")
		return
	}
	ok(w, items)
}

func (h *CampaignsHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		fail(w, http.StatusBadRequest, "invalid campaign id")
		return
	}

	item, err := h.svc.Get(r.Context(), id)
	if err != nil {
		fail(w, http.StatusNotFound, "campaign not found")
		return
	}
	ok(w, item)
}

func (h *CampaignsHandler) ImportContacts(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		fail(w, http.StatusBadRequest, "invalid campaign id")
		return
	}

	var req struct {
		Contacts []struct {
			CustomerID string `json:"customer_id"`
		} `json:"contacts"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fail(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if len(req.Contacts) == 0 {
		fail(w, http.StatusBadRequest, "contacts required")
		return
	}

	rows := make([]campaigns.ContactRow, 0, len(req.Contacts))
	for _, contact := range req.Contacts {
		customerID, err := uuid.Parse(contact.CustomerID)
		if err != nil {
			fail(w, http.StatusBadRequest, "invalid customer_id")
			return
		}
		rows = append(rows, campaigns.ContactRow{CustomerID: customerID})
	}

	if err := h.svc.ImportContacts(r.Context(), id, rows); err != nil {
		fail(w, http.StatusInternalServerError, "could not import contacts")
		return
	}

	ok(w, map[string]any{
		"campaign_id": id,
		"imported":    len(rows),
	})
}

func (h *CampaignsHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		fail(w, http.StatusBadRequest, "invalid campaign id")
		return
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Status == "" {
		fail(w, http.StatusBadRequest, "status required")
		return
	}

	if err := h.svc.UpdateStatus(r.Context(), id, req.Status); err != nil {
		fail(w, http.StatusInternalServerError, "could not update campaign status")
		return
	}

	ok(w, map[string]string{"status": req.Status})
}
