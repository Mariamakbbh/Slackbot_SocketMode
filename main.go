package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mariama/WebSocket_SlackBot/internal/config"
	consumer "github.com/mariama/WebSocket_SlackBot/internal/slack"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {

	// Load Env variables from .dot file
	config.Load()
	token := viper.GetString("slack.auth.token")
	appToken := viper.GetString("slack.app.token")

	l, _ := zap.NewDevelopment()
	logger := l.Sugar()

	// Create a new client to slack by giving token
	// Set debug to true while developing
	// Also add a ApplicationToken option to the client
	client := slack.New(token, slack.OptionDebug(true), slack.OptionAppLevelToken(appToken))
	// go-slack comes with a SocketMode package that we need to use that accepts a Slack client and outputs a Socket mode client instead
	socketClient := socketmode.New(
		client,
		socketmode.OptionDebug(true),
		// Option to set a custom logger
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)

	cfg := consumer.NewConfig(5, socketClient, logger)
	svc := consumer.NewService(client, socketClient, logger)

	sc := consumer.NewSocketChannelConsumer(
		cfg,
		svc,
		consumer.NewNoopMonitor(),
	)

	sc.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)

	// Block until we receive our signal.
	<-c

	sc.Stop()
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	logger.Info("Shutting down")
	os.Exit(0)

}
