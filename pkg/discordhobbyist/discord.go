package discordhobbyist

import (
	"bytes"
	"encoding/json"
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type DiscordAlertWebhook struct {
	Content DiscordCustomContent `json:"content"`
	Embeds  []struct {
		Title  string `json:"title"`
		Color  int    `json:"color"`
		Type   string `json:"type"`
		Footer struct {
			Text    string `json:"text"`
			IconURL string `json:"icon_url"`
		} `json:"footer"`
		URL   string `json:"url"`
		Image struct {
			URL string `json:"url"`
		} `json:"image"`
	}
	Username string `json:"username"`
}

func (d *DiscordAlertWebhook) GetImage() (string, error) {
	for _, embed := range d.Embeds {
		if embed.Image.URL != "" {
			return embed.Image.URL, nil
		}
	}

	return "", errors.New("no image found")
}

type DiscordCustomContent struct {
	Status   string          `json:"status"`
	Receiver string          `json:"receiver"`
	Alerts   []*DiscordAlert `json:"alerts"`
}

type DiscordAlert struct {
	Status       string            `json:"status"`
	Labels       map[string]string `json:"labels"`
	Annotations  map[string]string `json:"annotations"`
	StartsAt     string            `json:"starts_at"`
	EndsAt       string            `json:"ends_at"`
	GeneratorURL string            `json:"generator_url"`
	SilenceURL   string            `json:"silence_url"`
	DashboardURL string            `json:"dashboard_url"`
	PanelURL     string            `json:"panel_url"`
	ValueString  string            `json:"value_string"`
	Fingerprint  string            `json:"fingerprint"`
}

func (d *DiscordCustomContent) UnmarshalJSON(b []byte) error {
	type Alias DiscordCustomContent

	var alias Alias

	re := regexp.MustCompile(`,\s*([\]}])`)

	stripped := re.ReplaceAllString(string(b), "$1")

	buffer := new(bytes.Buffer)

	escaped, err := strconv.Unquote(stripped)
	if err != nil {
		return err
	}

	b = []byte(escaped)

	// Parse the JSON string into a map
	var data map[string]interface{}

	err = json.Unmarshal(b, &data)
	if err != nil {
		return err
	}

	// Print the modified JSON data
	output, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}

	if err := json.Compact(buffer, output); err != nil {
		return err
	}

	if err := json.Unmarshal(buffer.Bytes(), &alias); err != nil {
		return err
	}

	*d = DiscordCustomContent(alias)

	return nil
}

func ParseAsDiscordWebhook(body []byte) (*DiscordAlertWebhook, error) {
	var webhook DiscordAlertWebhook

	err := json.Unmarshal(body, &webhook)
	if err != nil {
		return nil, err
	}

	return &webhook, nil
}

func CreateNewMessageFromGrafanaWebhookAlert(webhook *DiscordAlertWebhook, alert *DiscordAlert) *discordgo.MessageSend {
	color := 5763719
	if alert.Status == "firing" {
		color = 15548997
	}

	links := []discordgo.MessageComponent{}
	if alert.SilenceURL != "" {
		links = append(links, discordgo.Button{
			Label: "Silence",
			URL:   alert.SilenceURL,
			Style: discordgo.LinkButton,
		})
	}

	if alert.DashboardURL != "" {
		links = append(links, discordgo.Button{
			Label: "Dashboard",
			URL:   alert.DashboardURL,
			Style: discordgo.LinkButton,
		})
	}

	if alert.PanelURL != "" {
		links = append(links, discordgo.Button{
			Label: "Panel",
			URL:   alert.PanelURL,
			Style: discordgo.LinkButton,
		})
	}

	if alert.GeneratorURL != "" {
		links = append(links, discordgo.Button{
			Label: "Alert",
			URL:   alert.GeneratorURL,
			Style: discordgo.LinkButton,
		})
	}

	fields := []*discordgo.MessageEmbedField{}
	for key, value := range alert.Labels {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   key,
			Value:  value,
			Inline: true,
		})
	}

	message := &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       alert.Labels["alertname"] + " - (" + strings.ToUpper(alert.Status) + ")",
				Description: alert.Annotations["description"],
				Color:       color,
				Footer: &discordgo.MessageEmbedFooter{
					Text: alert.ValueString,
				},
				Fields: fields,
			},
		},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: links,
			},
		},
	}

	image, err := webhook.GetImage()
	if err == nil {
		message.Embeds[0].Image = &discordgo.MessageEmbedImage{
			URL: image,
		}
	}

	return message
}
