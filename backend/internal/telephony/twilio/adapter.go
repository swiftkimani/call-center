package twilio

import (
	"context"
	"crypto/hmac"
	"crypto/sha1" //nolint:gosec // Twilio mandates SHA-1
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/yourorg/callcenter/internal/telephony"
)

var ErrInvalidSignature = errors.New("invalid twilio signature")

var _ telephony.Adapter = (*Adapter)(nil)

type Adapter struct {
	accountSID  string
	authToken   string
	webhookSecret string
}

func NewAdapter(accountSID, authToken, webhookSecret string) *Adapter {
	return &Adapter{accountSID: accountSID, authToken: authToken, webhookSecret: webhookSecret}
}

func (a *Adapter) VerifyWebhookSignature(r *http.Request) error {
	if a.authToken == "" {
		return nil
	}
	sig := r.Header.Get("X-Twilio-Signature")
	if sig == "" {
		return ErrInvalidSignature
	}

	if err := r.ParseForm(); err != nil {
		return fmt.Errorf("ParseForm: %w", err)
	}

	fullURL := "https://" + r.Host + r.RequestURI
	var keys []string
	for k := range r.Form {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var sb strings.Builder
	sb.WriteString(fullURL)
	for _, k := range keys {
		sb.WriteString(k)
		sb.WriteString(r.FormValue(k))
	}

	mac := hmac.New(sha1.New, []byte(a.authToken)) //nolint:gosec
	mac.Write([]byte(sb.String()))
	expected := base64.StdEncoding.EncodeToString(mac.Sum(nil))

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
		ProviderSID: r.FormValue("CallSid"),
		FromNumber:  r.FormValue("From"),
		ToNumber:    r.FormValue("To"),
		CallStatus:  r.FormValue("CallStatus"),
	}, nil
}

func (a *Adapter) ParseStatusWebhook(r *http.Request) (*telephony.StatusEvent, error) {
	if err := r.ParseForm(); err != nil {
		return nil, fmt.Errorf("ParseForm: %w", err)
	}
	var dur int
	if d := r.FormValue("CallDuration"); d != "" {
		dur, _ = strconv.Atoi(d)
	}
	ev := &telephony.StatusEvent{
		ProviderSID: r.FormValue("CallSid"),
		CallStatus:  r.FormValue("CallStatus"),
		Duration:    dur,
	}
	if price := r.FormValue("Price"); price != "" {
		f, err := strconv.ParseFloat(price, 64)
		if err == nil {
			cents := int32(f * 100)
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
		ProviderSID:  r.FormValue("CallSid"),
		RecordingURL: r.FormValue("RecordingUrl"),
	}, nil
}

func (a *Adapter) BuildInboundResponse(agentPhone, callbackURL string) ([]byte, error) {
	xml := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<Response>
  <Dial action="%s" record="record-from-ringing">
    <Number>%s</Number>
  </Dial>
</Response>`, callbackURL, agentPhone)
	return []byte(xml), nil
}

func (a *Adapter) DialAgent(_ context.Context, p telephony.DialParams) error {
	// Twilio REST API call would go here
	return fmt.Errorf("twilio DialAgent not yet implemented")
}
