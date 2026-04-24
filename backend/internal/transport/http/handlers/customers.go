package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/yourorg/callcenter/internal/customers"
)

type CustomersHandler struct {
	svc *customers.Service
}

func NewCustomersHandler(svc *customers.Service) *CustomersHandler {
	return &CustomersHandler{svc: svc}
}

func (h *CustomersHandler) Search(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 25
	}
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	results, err := h.svc.Search(r.Context(), q, limit, offset)
	if err != nil {
		fail(w, http.StatusInternalServerError, "search failed")
		return
	}
	ok(w, results)
}

func (h *CustomersHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		fail(w, http.StatusBadRequest, "invalid id")
		return
	}
	c, err := h.svc.Get(r.Context(), id)
	if err != nil {
		fail(w, http.StatusNotFound, "customer not found")
		return
	}
	ok(w, c)
}

func (h *CustomersHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		fail(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req struct {
		FullName string   `json:"full_name"`
		Email    string   `json:"email"`
		Tags     []string `json:"tags"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fail(w, http.StatusBadRequest, "invalid request body")
		return
	}

	c, err := h.svc.Update(r.Context(), id, req.FullName, req.Email, req.Tags)
	if err != nil {
		fail(w, http.StatusInternalServerError, "update failed")
		return
	}
	ok(w, c)
}
