package discordhobbyist

import "encoding/json"

type SlackWebhook struct {
	Text        string `json:"text"`
	Username    string `json:"username"`
	Attachments []struct {
		Title      string `json:"title"`
		TitleLink  string `json:"title_link"`
		Text       string `json:"text"`
		Fallback   string `json:"fallback"`
		Color      string `json:"color"`
		Footer     string `json:"footer"`
		FooterIcon string `json:"footer_icon"`
		TS         int    `json:"ts"`
	} `json:"attachments"`
}

func ParseAsSlackWebhook(body []byte) (*SlackWebhook, error) {
	var webhook SlackWebhook

	err := json.Unmarshal(body, &webhook)
	if err != nil {
		return nil, err
	}

	return &webhook, nil
}
