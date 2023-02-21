package discordhobbyist

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

type DiscordBot struct {
	log     logrus.FieldLogger
	config  *Config
	session *discordgo.Session

	http *HTTPServer
}

func NewDiscordBot(log logrus.FieldLogger, config *Config) *DiscordBot {
	d := &DiscordBot{
		log:    log,
		config: config,
	}

	d.http = NewHTTPServer(log, []func(ctx context.Context, params httprouter.Params, body []byte) error{
		d.handleChannelRequest,
	})

	return d
}

func (d *DiscordBot) Start() error {
	d.log.WithField("config", d.config.BotToken).Info("Starting discordhobbyist bot")

	if err := d.config.Validate(); err != nil {
		return err
	}

	session, err := discordgo.New("Bot " + d.config.BotToken)
	if err != nil {
		return err
	}

	d.session = session

	d.log.Info("Opening discord session")

	err = session.Open()
	if err != nil {
		return err
	}

	d.log.Info("Connected to discord")

	if err := d.http.Start(d.config.HTTPAddr); err != nil {
		return err
	}

	d.log.Info("Started HTTP server")

	return nil
}

func (d *DiscordBot) Channels() (map[string]*discordgo.Channel, error) {
	routes := make(map[string]*discordgo.Channel)

	for _, guild := range d.session.State.Guilds {
		channels, err := d.session.GuildChannels(guild.ID)
		if err != nil {
			return nil, err
		}

		groups := make(map[string]*discordgo.Channel)

		for _, channel := range channels {
			if channel.Type != discordgo.ChannelTypeGuildCategory {
				continue
			}

			groups[channel.ID] = channel
		}

		for _, channel := range channels {
			// Check if channel is a guild text channel and not a voice or DM channel
			if channel.Type != discordgo.ChannelTypeGuildText {
				continue
			}

			// Check if channel is in a group
			group, ok := groups[channel.ParentID]
			if !ok {
				continue
			}

			key := "/" + strings.ToLower(group.Name) + "/" + strings.ToLower(channel.Name)

			routes[key] = channel
		}
	}

	return routes, nil
}

func (d *DiscordBot) Stop() error {
	d.log.Info("Stopping discordhobbyist bot")

	return d.session.Close()
}

func (d *DiscordBot) handleChannelRequest(ctx context.Context, params httprouter.Params, body []byte) error {
	channels, err := d.Channels()
	if err != nil {
		d.log.WithError(err).Error("error fetching channels")

		return err
	}

	key := GetChannelKey(params.ByName("group"), params.ByName("channel"))

	channel, ok := channels[key]
	if !ok {
		d.log.Error("no route found for path")

		err = d.CreateInfoMessage(ctx, fmt.Sprintf("no route found for path %s. Is there a typo in the configured URL?", key))
		if err != nil {
			d.log.WithError(err).Error("error sending info message")
		}

		return errors.New("no route found for path")
	}

	webhook, err := ParseAsDiscordWebhook(body)
	if err != nil {
		d.log.WithError(err).WithField("body", string(body)).Error("error parsing request body as discord webhook")

		err := d.CreateInfoMessage(ctx, "An invalid payload was sent to the webhook. Please check the webhook configuration in Grafana, and logs for Discordhobbyist for the full payload.")
		if err != nil {
			d.log.WithError(err).Error("error sending info message")
		}

		return errors.New("error parsing request body as discord webhook")
	}

	for _, alert := range webhook.Content.Alerts {
		// Send the request body to the channel
		payload := CreateNewMessageFromGrafanaWebhookAlert(webhook, alert)

		_, err := d.session.ChannelMessageSendComplex(channel.ID, payload)
		if err != nil {
			d.log.WithError(err).Error("error sending message to channel")

			return errors.New("error sending message to channel")
		}
	}

	return nil
}

func (d *DiscordBot) CreateInfoMessage(ctx context.Context, msg string) error {
	channels, err := d.Channels()
	if err != nil {
		return err
	}

	channel, ok := channels[d.config.InfoChannelKey]
	if !ok {
		return errors.New("no route found for info channel")
	}

	_, err = d.session.ChannelMessageSend(channel.ID, msg)
	if err != nil {
		return err
	}

	return nil
}
