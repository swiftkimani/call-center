package hub

import (
	"time"

	"github.com/google/uuid"
)

const (
	EventIncomingCall    = "incoming_call"
	EventCallEnded       = "call_ended"
	EventStatusChanged   = "status_changed"
	EventWhisperStarted  = "whisper_started"
	EventQueueUpdate     = "queue_update"
	EventAgentStateChange = "agent_state_change"
	EventSLABreach       = "sla_breach"
	EventBroadcast       = "broadcast"
)

type Event struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	Data      any       `json:"data"`
}

func NewEvent(eventType string, data any) Event {
	return Event{
		ID:        uuid.New().String(),
		Type:      eventType,
		Timestamp: time.Now().UTC(),
		Data:      data,
	}
}
