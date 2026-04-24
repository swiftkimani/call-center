package webhooks

import (
	"context"
	"fmt"
	"net/http"

	"github.com/yourorg/callcenter/internal/agents"
	"github.com/yourorg/callcenter/internal/broker"
	"github.com/yourorg/callcenter/internal/calls"
	"github.com/yourorg/callcenter/internal/config"
	"github.com/yourorg/callcenter/internal/customers"
	"github.com/yourorg/callcenter/internal/hub"
	"github.com/yourorg/callcenter/internal/queues"
	"github.com/yourorg/callcenter/internal/telephony"
	"github.com/rs/zerolog/log"
)

type ATHandler struct {
	adapter      telephony.Adapter
	cfg          *config.Config
	callsSvc     *calls.Service
	customersSvc *customers.Service
	agentsSvc    *agents.Service
	queueSvc     *queues.Service
	queueRouter  *queues.Router
	hub          *hub.Hub
	publisher    *broker.Publisher
}

func NewATHandler(
	adapter telephony.Adapter,
	cfg *config.Config,
	callsSvc *calls.Service,
	customersSvc *customers.Service,
	agentsSvc *agents.Service,
	queueSvc *queues.Service,
	queueRouter *queues.Router,
	h *hub.Hub,
	publisher *broker.Publisher,
) *ATHandler {
	return &ATHandler{
		adapter:      adapter,
		cfg:          cfg,
		callsSvc:     callsSvc,
		customersSvc: customersSvc,
		agentsSvc:    agentsSvc,
		queueSvc:     queueSvc,
		queueRouter:  queueRouter,
		hub:          h,
		publisher:    publisher,
	}
}

// Inbound handles POST /webhooks/voice/inbound
func (h *ATHandler) Inbound(w http.ResponseWriter, r *http.Request) {
	if err := h.adapter.VerifyWebhookSignature(r); err != nil {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	evt, err := h.adapter.ParseInboundWebhook(r)
	if err != nil {
		log.Error().Err(err).Msg("parse inbound webhook")
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	log.Info().
		Str("sid", evt.ProviderSID).
		Str("from", evt.FromNumber).
		Msg("inbound call")

	ctx := context.Background()

	queue, err := h.queueSvc.DefaultQueue(ctx)
	if err != nil {
		log.Error().Err(err).Msg("resolve default queue")
		http.Error(w, "queue unavailable", http.StatusServiceUnavailable)
		return
	}

	// Find or create customer
	customer, _ := h.customersSvc.FindOrCreate(ctx, evt.FromNumber, "")

	callParams := calls.CreateInboundParams{
		ProviderSID: evt.ProviderSID,
		FromNumber:  evt.FromNumber,
		ToNumber:    evt.ToNumber,
		QueueID:     &queue.ID,
	}
	if customer != nil {
		callParams.CustomerID = &customer.ID
	}

	call, err := h.callsSvc.CreateInbound(ctx, callParams)
	if err != nil {
		log.Error().Err(err).Msg("CreateInbound")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	// Select the best available agent
	agent, err := h.queueRouter.SelectAgent(ctx, queue.SkillsRequired)
	if err != nil {
		if enqueueErr := h.queueSvc.Enqueue(ctx, queue.ID, call.ID); enqueueErr != nil {
			log.Error().Err(enqueueErr).Str("queue_id", queue.ID.String()).Msg("enqueue inbound call")
		} else if snapshot, snapErr := h.queueSvc.Snapshot(ctx, queue.ID); snapErr == nil {
			h.hub.BroadcastSupervisors(hub.NewEvent(hub.EventQueueUpdate, snapshot))
		}

		// No agent available — respond with a hold message
		w.Header().Set("Content-Type", "application/xml")
		fmt.Fprintf(w, `<?xml version="1.0" encoding="UTF-8"?><Response><Say>All agents are busy. Please hold.</Say></Response>`)
		return
	}

	// Mark agent busy
	if err := h.agentsSvc.SetStatus(ctx, agent.ID, "busy"); err != nil {
		log.Warn().Err(err).Str("agent_id", agent.ID.String()).Msg("set agent busy")
	}

	// Answer the call
	if err := h.callsSvc.Answer(ctx, call.ID, agent.ID); err != nil {
		log.Error().Err(err).Msg("Answer call")
	}

	// Notify the agent's browser via WebSocket
	h.hub.PublishToAgent(ctx, agent.ID.String(), hub.NewEvent(hub.EventIncomingCall, map[string]any{
		"call_id":     call.ID,
		"from_number": evt.FromNumber,
		"customer":    customer,
	}))

	// Build XML response to bridge the call
	callbackURL := h.cfg.BaseURL + "/webhooks/voice/status"
	xml, err := h.adapter.BuildInboundResponse(agent.Extension, callbackURL)
	if err != nil {
		log.Error().Err(err).Msg("BuildInboundResponse")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	w.Write(xml)
}

// Status handles POST /webhooks/voice/status
func (h *ATHandler) Status(w http.ResponseWriter, r *http.Request) {
	if err := h.adapter.VerifyWebhookSignature(r); err != nil {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	evt, err := h.adapter.ParseStatusWebhook(r)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	call, err := h.callsSvc.GetByProviderSID(ctx, evt.ProviderSID)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	status := mapATStatus(evt.CallStatus)
	if err := h.callsSvc.End(ctx, call.ID, status, evt.CostCents); err != nil {
		log.Error().Err(err).Msg("End call")
	}

	if call.QueueID != nil {
		if err := h.queueSvc.Remove(ctx, *call.QueueID, call.ID); err != nil {
			log.Warn().Err(err).Str("queue_id", call.QueueID.String()).Str("call_id", call.ID.String()).Msg("remove call from queue")
		} else if snapshot, err := h.queueSvc.Snapshot(ctx, *call.QueueID); err == nil {
			h.hub.BroadcastSupervisors(hub.NewEvent(hub.EventQueueUpdate, snapshot))
		}
	}

	// Move agent back to wrap-up
	if call.AgentID != nil {
		h.agentsSvc.SetStatus(ctx, *call.AgentID, "wrap_up")
		h.hub.PublishToAgent(ctx, call.AgentID.String(), hub.NewEvent(hub.EventCallEnded, map[string]any{
			"call_id": call.ID,
			"status":  status,
		}))
	}

	// Enqueue recording download if URL provided
	if evt.RecordingURL != "" {
		h.publisher.Publish(ctx, broker.QueueRecordingReady, broker.RecordingReadyMsg{
			CallID:               call.ID.String(),
			ProviderRecordingURL: evt.RecordingURL,
		})
	}

	w.WriteHeader(http.StatusOK)
}

// RecordingReady handles POST /webhooks/voice/recording
func (h *ATHandler) RecordingReady(w http.ResponseWriter, r *http.Request) {
	if err := h.adapter.VerifyWebhookSignature(r); err != nil {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	evt, err := h.adapter.ParseRecordingWebhook(r)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	call, err := h.callsSvc.GetByProviderSID(ctx, evt.ProviderSID)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	h.publisher.Publish(ctx, broker.QueueRecordingReady, broker.RecordingReadyMsg{
		CallID:               call.ID.String(),
		ProviderRecordingURL: evt.RecordingURL,
	})

	w.WriteHeader(http.StatusOK)
}

func mapATStatus(atStatus string) string {
	switch atStatus {
	case "Completed":
		return "completed"
	case "NoAnswer":
		return "no_answer"
	case "Failed":
		return "failed"
	case "Abandoned":
		return "abandoned"
	default:
		return "completed"
	}
}
