package africas_talking

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/yourorg/callcenter/internal/telephony"
)

var ErrInvalidSignature = errors.New("invalid webhook signature")

func (a *Adapter) VerifyWebhookSignature(r *http.Request) error {
	if a.webhookSecret == "" {
		return nil // signature checking disabled in dev
	}
	sig := r.Header.Get("X-AT-Signature")
	if sig == "" {
		return ErrInvalidSignature
	}
	if err := r.ParseForm(); err != nil {
		return fmt.Errorf("ParseForm: %w", err)
	}
	mac := hmac.New(sha256.New, []byte(a.webhookSecret))
	mac.Write([]byte(r.Form.Encode()))
	expected := hex.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(sig), []byte(expected)) {
		return ErrInvalidSignature
	}
	return nil
}

func (a *Adapter) ParseInboundWebhook(r *http.Request) (*telephony.InboundEvent, error) {
	if err := r.ParseForm(); err != nil {
		return nil, fmt.Errorf("ParseForm: %w", err)
	}
	return &telephony.InboundEvent{
		ProviderSID: r.FormValue("sessionId"),
		FromNumber:  r.FormValue("callerNumber"),
		ToNumber:    r.FormValue("destinationNumber"),
		CallStatus:  r.FormValue("callSessionState"),
	}, nil
}

func (a *Adapter) ParseStatusWebhook(r *http.Request) (*telephony.StatusEvent, error) {
	if err := r.ParseForm(); err != nil {
		return nil, fmt.Errorf("ParseForm: %w", err)
	}

	var dur int
	if d := r.FormValue("durationInSeconds"); d != "" {
		dur, _ = strconv.Atoi(d)
	}

	ev := &telephony.StatusEvent{
		ProviderSID:  r.FormValue("sessionId"),
		CallStatus:   r.FormValue("callSessionState"),
		Duration:     dur,
		RecordingURL: r.FormValue("recordingUrl"),
	}

	if cost := r.FormValue("totalCost"); cost != "" {
		c, err := strconv.ParseFloat(cost, 64)
		if err == nil {
			cents := int32(c * 100)
			ev.CostCents = &cents
		}
	}
	return ev, nil
}

func (a *Adapter) ParseRecordingWebhook(r *http.Request) (*telephony.RecordingEvent, error) {
	if err := r.ParseForm(); err != nil {
		return nil, fmt.Errorf("ParseForm: %w", err)
	}
	return &telephony.RecordingEvent{
		ProviderSID:  r.FormValue("sessionId"),
		RecordingURL: r.FormValue("recordingUrl"),
	}, nil
}
