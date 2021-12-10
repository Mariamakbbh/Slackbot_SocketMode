package config

import (
	"github.com/spf13/viper"
)

func Load() {
	viper.BindEnv("slack.auth.token", "SLACK_AUTH_TOKEN")
	viper.BindEnv("slack.app.token", "SLACK_APP_TOKEN")
	viper.SetDefault("slack.auth.token", "dummy")
	viper.SetDefault("slack.app.token", "dummy")
	//viper.BindEnv("slack.channel", "#websocket-bot")

}
