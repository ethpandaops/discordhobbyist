package discordhobbyist

import (
	"context"
	"errors"
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

	channel, ok := channels[GetChannelKey(params.ByName("group"), params.ByName("channel"))]
	if !ok {
		d.log.Error("no route found for path")

		return errors.New("no route found for path")
	}

	webhook, err := ParseAsDiscordWebhook(body)
	if err != nil {
		d.log.WithError(err).Error("error parsing request body as discord webhook")

		return errors.New("error parsing request body as discord webhook")
	}

	for _, alert := range webhook.Content.Alerts {
		// Send the request body to the channel
		payload := CreateNewMessageFromGrafanaWebhookAlert(webhook, alert)

		// Check if the channel already contains a message with the same fingerprint
		if strings.EqualFold(alert.Status, "resolved") {
			firingMessage, err := d.FindFiringAlertMessageWithFingerprint(ctx, channel, alert.Fingerprint)
			if err == nil {
				// Add a reply to the original message
				_, err = d.session.ChannelMessageSendEmbedsReply(channel.ID, payload.Embeds, firingMessage.Reference())
				if err == nil {
					continue
				}

				// This sucks but we'll just send the resolved message without the reference to the original
				d.log.WithError(err).Error("error sending reply resolved message to channel")
			}

			// We'll still send the resolved message just without the reference to the original
			d.log.WithError(err).Warn("error finding original firing message with fingerprint")
		}

		_, err := d.session.ChannelMessageSendComplex(channel.ID, payload)
		if err != nil {
			d.log.WithError(err).Error("error sending message to channel")

			return errors.New("error sending message to channel")
		}
	}

	return nil
}

func (d *DiscordBot) FindFiringAlertMessageWithFingerprint(ctx context.Context, channel *discordgo.Channel, fingerprint string) (*discordgo.Message, error) {
	messages, err := d.session.ChannelMessages(channel.ID, 100, "", "", "")
	if err != nil {
		return nil, err
	}

	for _, message := range messages {
		// Ignore messages that are not from the us
		if message.Author.ID != d.session.State.User.ID {
			continue
		}

		// Ignore messages that are not alerts that are firing
		if !strings.Contains(strings.ToLower(message.Content), "firing") {
			continue
		}

		// Inspect the embeds for the fingerprint
		for _, embed := range message.Embeds {
			for _, field := range embed.Fields {
				if field.Name == "Fingerprint" && field.Value == fingerprint {
					return message, nil
				}
			}
		}
	}

	return nil, errors.New("channel does not contain a valid message with the specified fingerprint")
}
