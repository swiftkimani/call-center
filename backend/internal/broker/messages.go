package broker

const (
	ExchangeDirect = ""

	QueueRecordingReady = "recording.ready"
	QueueCampaignDial   = "campaign.dial"
)

type RecordingReadyMsg struct {
	CallID             string `json:"call_id"`
	ProviderRecordingURL string `json:"provider_recording_url"`
}

type CampaignDialMsg struct {
	CampaignID        string `json:"campaign_id"`
	ContactID         string `json:"contact_id"`
	CustomerID        string `json:"customer_id"`
	CustomerPhone     string `json:"customer_phone"`
}
