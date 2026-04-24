package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/yourorg/callcenter/internal/auth"
	"github.com/yourorg/callcenter/internal/config"
	"github.com/yourorg/callcenter/internal/transport/http/handlers"
	"github.com/yourorg/callcenter/internal/transport/http/webhooks"
	"github.com/yourorg/callcenter/internal/transport/ws"
	"github.com/yourorg/callcenter/internal/hub"
)

type Deps struct {
	Cfg         *config.Config
	AuthHandler *handlers.AuthHandler
	MeHandler   *handlers.MeHandler
	AgentsHandler *handlers.AgentsHandler
	CampaignsHandler *handlers.CampaignsHandler
	CallsHandler  *handlers.CallsHandler
	CustomersHandler *handlers.CustomersHandler
	ReportsHandler   *handlers.ReportsHandler
	SupervisorHandler *handlers.SupervisorHandler
	ATWebhookHandler  *webhooks.ATHandler
	Hub               *hub.Hub
}

func NewRouter(d Deps) http.Handler {
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.RealIP)
	r.Use(RequestID)
	r.Use(RequestLogger)
	r.Use(Recoverer)

	// Health check (no auth)
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// Public auth routes
	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/login", d.AuthHandler.Login)
		r.Post("/refresh", d.AuthHandler.Refresh)
		r.Post("/logout", d.AuthHandler.Logout)
	})

	// Authenticated API routes
	r.Group(func(r chi.Router) {
		r.Use(auth.BearerAuth(d.Cfg))

		r.Get("/api/v1/me", d.MeHandler.Get)
		r.Post("/api/v1/agents/{id}/status", d.AgentsHandler.SetStatus)
		r.Post("/api/v1/agents/heartbeat", d.AgentsHandler.Heartbeat)

		r.Post("/api/v1/calls/outbound", d.CallsHandler.InitiateOutbound)
		r.Get("/api/v1/calls", d.CallsHandler.List)
		r.Post("/api/v1/calls/{id}/disposition", d.CallsHandler.Disposition)
		r.Get("/api/v1/calls/{id}/recording", d.CallsHandler.Recording)

		r.Get("/api/v1/campaigns", d.CampaignsHandler.List)
		r.Get("/api/v1/campaigns/{id}", d.CampaignsHandler.Get)

		r.Get("/api/v1/customers", d.CustomersHandler.Search)
		r.Get("/api/v1/customers/{id}", d.CustomersHandler.Get)
		r.Put("/api/v1/customers/{id}", d.CustomersHandler.Update)

		// Supervisor-only routes
		r.Group(func(r chi.Router) {
			r.Use(auth.RoleRequired("supervisor", "admin"))
			r.Get("/api/v1/queues/{id}/live", d.SupervisorHandler.QueueLive)
			r.Post("/api/v1/supervisor/{call_id}/whisper", d.SupervisorHandler.Whisper)
			r.Post("/api/v1/supervisor/{call_id}/barge", d.SupervisorHandler.Barge)
			r.Post("/api/v1/campaigns", d.CampaignsHandler.Create)
			r.Post("/api/v1/campaigns/{id}/contacts", d.CampaignsHandler.ImportContacts)
			r.Post("/api/v1/campaigns/{id}/status", d.CampaignsHandler.UpdateStatus)
			r.Get("/api/v1/reports/daily", d.ReportsHandler.Daily)
		})
	})

	// WebSocket endpoints (auth via ?token= query param)
	r.Get("/ws/agent", ws.AgentHandler(d.Hub, d.Cfg))
	r.Get("/ws/supervisor", ws.SupervisorHandler(d.Hub, d.Cfg))

	// Telephony webhooks (HMAC-authenticated, no JWT)
	r.Route("/webhooks", func(r chi.Router) {
		r.Post("/voice/inbound", d.ATWebhookHandler.Inbound)
		r.Post("/voice/status", d.ATWebhookHandler.Status)
		r.Post("/voice/recording", d.ATWebhookHandler.RecordingReady)
	})

	return r
}
