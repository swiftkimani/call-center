package handlers

import (
	"net/http"
	"time"

	"github.com/yourorg/callcenter/internal/reports"
)

type ReportsHandler struct {
	svc *reports.Service
}

func NewReportsHandler(svc *reports.Service) *ReportsHandler {
	return &ReportsHandler{svc: svc}
}

func (h *ReportsHandler) Daily(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")
	date := time.Now().UTC()
	if dateStr != "" {
		var err error
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			fail(w, http.StatusBadRequest, "date must be YYYY-MM-DD")
			return
		}
	}

	if r.URL.Query().Get("format") == "csv" {
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", `attachment; filename="report.csv"`)
		if err := h.svc.ExportCSV(r.Context(), w, date); err != nil {
			fail(w, http.StatusInternalServerError, "export failed")
		}
		return
	}

	summary, err := h.svc.Daily(r.Context(), date)
	if err != nil {
		fail(w, http.StatusInternalServerError, "report failed")
		return
	}
	ok(w, summary)
}
