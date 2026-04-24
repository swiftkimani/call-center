package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/yourorg/callcenter/internal/auth"
	"github.com/yourorg/callcenter/internal/calls"
	"github.com/yourorg/callcenter/internal/customers"
	"github.com/yourorg/callcenter/internal/storage"
)

type CallsHandler struct {
	callsSvc     *calls.Service
	customersSvc *customers.Service
	store        *storage.Client
}

func NewCallsHandler(callsSvc *calls.Service, customersSvc *customers.Service, store *storage.Client) *CallsHandler {
	return &CallsHandler{callsSvc: callsSvc, customersSvc: customersSvc, store: store}
}

func (h *CallsHandler) List(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 25
	}
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	rows, err := h.callsSvc.List(r.Context(), limit, offset)
	if err != nil {
		fail(w, http.StatusInternalServerError, "could not fetch calls")
		return
	}
	ok(w, rows)
}

func (h *CallsHandler) Disposition(w http.ResponseWriter, r *http.Request) {
	callID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		fail(w, http.StatusBadRequest, "invalid call id")
		return
	}

	claims := auth.ClaimsFromContext(r.Context())
	if claims.AgentID == nil {
		fail(w, http.StatusForbidden, "agents only")
		return
	}

	var req struct {
		Category string `json:"category"`
		Notes    string `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Category == "" {
		fail(w, http.StatusBadRequest, "category required")
		return
	}

	if err := h.callsSvc.SaveDisposition(r.Context(), callID, *claims.AgentID, req.Category, req.Notes); err != nil {
		fail(w, http.StatusInternalServerError, "could not save disposition")
		return
	}

	ok(w, nil)
}

func (h *CallsHandler) Recording(w http.ResponseWriter, r *http.Request) {
	callID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		fail(w, http.StatusBadRequest, "invalid call id")
		return
	}

	call, err := h.callsSvc.GetByID(r.Context(), callID)
	if err != nil {
		fail(w, http.StatusNotFound, "call not found")
		return
	}

	if call.RecordingUrl == nil || *call.RecordingUrl == "" {
		fail(w, http.StatusNotFound, "no recording available")
		return
	}

	url, err := h.store.PresignedGetURL(r.Context(), *call.RecordingUrl, 15*time.Minute)
	if err != nil {
		fail(w, http.StatusInternalServerError, "could not generate recording URL")
		return
	}

	ok(w, map[string]string{"url": url})
}

func (h *CallsHandler) InitiateOutbound(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromContext(r.Context())
	if claims.AgentID == nil {
		fail(w, http.StatusForbidden, "agents only")
		return
	}

	var req struct {
		CustomerID    string `json:"customer_id"`
		CustomerPhone string `json:"customer_phone"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fail(w, http.StatusBadRequest, "invalid request body")
		return
	}

	customerID, err := uuid.Parse(req.CustomerID)
	if err != nil {
		fail(w, http.StatusBadRequest, "invalid customer_id")
		return
	}

	call, err := h.callsSvc.InitiateOutbound(r.Context(), *claims.AgentID, customerID, req.CustomerPhone)
	if err != nil {
		fail(w, http.StatusBadRequest, err.Error())
		return
	}

	ok(w, call)
}
