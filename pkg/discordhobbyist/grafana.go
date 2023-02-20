package discordhobbyist

import "encoding/json"

type GrafanaWebhookAlert struct {
	Status       string `json:"status"`
	Labels       map[string]string
	Annotations  map[string]string
	StartsAt     string `json:"startsAt"`
	EndsAt       string `json:"endsAt"`
	GeneratorURL string `json:"generatorURL"`
	Fingerprint  string `json:"fingerprint"`
	SilenceURL   string `json:"silenceURL"`
	DashboardURL string `json:"dashboardURL"`
	PanelURL     string `json:"panelURL"`
	ValueString  string `json:"valueString"`
}

type GrafanaWebhook struct {
	Receiver          string                `json:"receiver"`
	Status            string                `json:"status"`
	Alerts            []GrafanaWebhookAlert `json:"alerts"`
	GroupLabels       map[string]string     `json:"groupLabels"`
	CommonLabels      map[string]string     `json:"commonLabels"`
	CommonAnnotations map[string]string     `json:"commonAnnotations"`
	ExternalURL       string                `json:"externalURL"`
	Version           string                `json:"version"`
	GroupKey          string                `json:"groupKey"`
	TruncatedAlerts   int                   `json:"truncatedAlerts"`
	OrgId             int                   `json:"orgId"`
	Title             string                `json:"title"`
	State             string                `json:"state"`
	Message           string                `json:"message"`
}

func ParseAsGrafanaWebhook(body []byte) (*GrafanaWebhook, error) {
	var webhook GrafanaWebhook

	err := json.Unmarshal(body, &webhook)
	if err != nil {
		return nil, err
	}

	return &webhook, nil
}
