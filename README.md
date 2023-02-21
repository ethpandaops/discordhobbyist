# DiscordHobbyist

> A Grafana -> Discord webhook that enables rich notifications and a more complex set of features than the default Discord integration.

## Usage

```bash
A discord bot that simply forwards messages to a discord channel

Usage:
  discordhobbyist [flags]

Flags:
  -h, --help     help for discordhobbyist
  -t, --toggle   Help message for toggle
```

## Configuration

All configuration is passed via environment variables.

| Variable | Description | Default | Required | Example |
| --- | --- | --- | --- |
| GUILD_ID | The ID of the guild to send messages to |  | Yes | 1234567890 |
| BOT_TOKEN | The bot token to use to authenticate with Discord |  | Yes | abcdefghijklmnopqrstuvwxyz1234567890 |
| APP_ID | The application ID of the bot |  | Yes | 1234567890 |
| INFO_CHANNEL_KEY | The key of the channel to send info messages to |  | Yes | /group_channel_name_here/child_channel_name_here |
| HTTP_ADDR | The address to listen on for HTTP requests | | Yes | :8080 |

## Getting Started

### Docker

Available as a docker image at `ethpandaops/discordhobbyist`

#### Images

- `latest` - distroless, multiarch
- `latest` - debian - debian, multiarch
- `$version` - distroless, multiarch, pinned to a release (i.e. 0.4.0)
- `$version-debian` - debian, multiarch, pinned to a release (i.e. 0.4.0-debian)

### Quick start

```bash
docker run -e GUILD_ID -e BOT_TOKEN -e APP_ID -e INFO_CHANNEL_KEY -e HTTP_ADDR -d --name discordhobbyist -p 8080:8080 -it ethpandaops/discordhobbyist
```

### Standalone

**Downloading a release**
Available [here](https://github.com/ethpandaops/discordhobbyist/releases)

## Contributing

Contributions are greatly appreciated! Pull requests will be reviewed and merged promptly if you're interested in improving the hobbyist!

1. Fork the project
2. Create your feature branch:
    - `git checkout -b feat/new-feature-a`
3. Commit your changes:
    - `git commit -m 'feat: Implement feature a`
4. Push to the branch:
    -`git push origin feat/feature-a`
5. Open a pull request

## Contact

Sam - [@samcmau](https://twitter.com/samcmau)
