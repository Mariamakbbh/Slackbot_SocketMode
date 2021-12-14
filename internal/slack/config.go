package slack

import (
	"go.uber.org/zap"

	"github.com/slack-go/slack/socketmode"
)

type Config struct {
	logger  *zap.SugaredLogger
	workers int
	client  *socketmode.Client
}

func (c *Config) Workers() int {
	return c.workers
}

func (c *Config) GetSocketClient() *socketmode.Client {
	return c.client
}

// GetLogger returns internal logger used for operations
func (c *Config) GetLogger() *zap.SugaredLogger {
	return c.logger
}

// NewConfig returns an object implementing ConfigInterface
func NewConfig(workers int, client *socketmode.Client, logger *zap.SugaredLogger) ConfigInterface {
	cfg := new(Config)

	if cfg.workers = workers; cfg.workers < 1 {
		cfg.workers = 5
	}

	cfg.client = client
	cfg.logger = logger

	if logger == nil {
		l, _ := zap.NewDevelopment()
		cfg.logger = l.Sugar()
	}

	return cfg
}
