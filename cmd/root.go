package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/ethpandaops/discordhobbyist/pkg/discordhobbyist"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "discordhobbyist",
	Short: "A discord bot that simply forwards messages to a discord channel",
	Run: func(cmd *cobra.Command, args []string) {
		log := logrus.New()

		config := &discordhobbyist.Config{
			GuildID:        os.Getenv("GUILD_ID"),
			BotToken:       os.Getenv("BOT_TOKEN"),
			AppID:          os.Getenv("APP_ID"),
			InfoChannelKey: os.Getenv("INFO_CHANNEL_KEY"),
			HTTPAddr:       os.Getenv("HTTP_ADDR"),
		}

		hobbyist := discordhobbyist.NewDiscordBot(log, config)

		if err := hobbyist.Start(); err != nil {
			log.WithError(err).Fatal("error starting bot")
		}

		log.Info("Bot is now running.  Press CTRL-C to exit.")

		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-sc

		// Cleanly close down the Discord session.
		if err := hobbyist.Stop(); err != nil {
			log.WithError(err).Fatal("error stopping bot")
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
