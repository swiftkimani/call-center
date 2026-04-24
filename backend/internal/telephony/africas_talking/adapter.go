package africas_talking

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/yourorg/callcenter/internal/telephony"
)

const atVoiceURL = "https://voice.africastalking.com/call"

var _ telephony.Adapter = (*Adapter)(nil)

type Adapter struct {
	apiKey        string
	username      string
	webhookSecret string
	httpClient    *http.Client
}

func NewAdapter(apiKey, username, webhookSecret string) *Adapter {
	return &Adapter{
		apiKey:        apiKey,
		username:      username,
		webhookSecret: webhookSecret,
		httpClient:    &http.Client{},
	}
}

func (a *Adapter) DialAgent(ctx context.Context, p telephony.DialParams) error {
	form := url.Values{
		"username":    {a.username},
		"from":        {p.FromNumber},
		"to":          {p.ToNumber},
		"callbackUrl": {p.CallbackURL},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, atVoiceURL,
		bytes.NewBufferString(form.Encode()))
	if err != nil {
		return fmt.Errorf("NewRequest: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("apiKey", a.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return fmt.Errorf("AT dial failed (%d): %s", resp.StatusCode, body)
	}

	var result struct {
		Entries []struct {
			Status      string `json:"status"`
			PhoneNumber string `json:"phoneNumber"`
		} `json:"entries"`
	}
	if err := json.Unmarshal(body, &result); err == nil {
		for _, e := range result.Entries {
			if e.Status != "Success" {
				return fmt.Errorf("AT dial entry failed: %s %s", e.PhoneNumber, e.Status)
			}
		}
	}
	return nil
}
