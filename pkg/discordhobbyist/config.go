package discordhobbyist

import "errors"

type Config struct {
	GuildID  string
	BotToken string
	AppID    string

	InfoChannelKey string

	HTTPAddr string
}

func (c *Config) Validate() error {
	if c.GuildID == "" {
		return errors.New("guild id is required")
	}

	if c.BotToken == "" {
		return errors.New("bot token is required")
	}

	if c.AppID == "" {
		return errors.New("app id is required")
	}

	if c.InfoChannelKey == "" {
		return errors.New("info channel key is required")
	}

	if c.HTTPAddr == "" {
		return errors.New("http addr is required")
	}

	return nil
}
