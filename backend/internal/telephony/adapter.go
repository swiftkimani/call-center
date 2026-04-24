package telephony

import (
	"context"
	"net/http"
)

// InboundEvent is parsed from a provider's inbound call webhook.
type InboundEvent struct {
	ProviderSID string
	FromNumber  string
	ToNumber    string
	CallStatus  string
}

// StatusEvent is parsed from a provider's call status webhook.
type StatusEvent struct {
	ProviderSID string
	CallStatus  string
	Duration    int
	CostCents   *int32
	RecordingURL string
}

// RecordingEvent is parsed from a provider's recording-ready webhook.
type RecordingEvent struct {
	ProviderSID  string
	RecordingURL string
}

// DialParams carries everything needed to initiate an outbound call.
type DialParams struct {
	CallID      string
	FromNumber  string
	ToNumber    string
	CallbackURL string // status webhook
}

// Adapter is the interface all telephony providers must satisfy.
type Adapter interface {
	// VerifyWebhookSignature returns nil if the request is genuine.
	VerifyWebhookSignature(r *http.Request) error

	// ParseInboundWebhook deserialises the provider's inbound call POST.
	ParseInboundWebhook(r *http.Request) (*InboundEvent, error)

	// ParseStatusWebhook deserialises call-status update POSTs.
	ParseStatusWebhook(r *http.Request) (*StatusEvent, error)

	// ParseRecordingWebhook deserialises the recording-ready POST.
	ParseRecordingWebhook(r *http.Request) (*RecordingEvent, error)

	// BuildInboundResponse returns the XML bytes the provider expects back.
	BuildInboundResponse(agentPhone, callbackURL string) ([]byte, error)

	// DialAgent places an outbound call to an agent's device / browser.
	DialAgent(ctx context.Context, p DialParams) error
}
