package config

import (
	"github.com/spf13/viper"
)

func Load() {
	viper.SetDefault("slack.auth.token", "dummy")
	viper.SetDefault("slack.app.token", "dummy")
	//viper.BindEnv("slack.channel", "#websocket-bot")

}
