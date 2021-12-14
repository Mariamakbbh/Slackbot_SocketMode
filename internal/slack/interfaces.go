package slack

import (
	"github.com/slack-go/slack/socketmode"
	"go.uber.org/zap"
)

type ConsumerInterface interface {
	Start()
	Stop()
}

type ServiceInterface interface {
	Flow(event socketmode.Event) error
}

type MonitorInterface interface {
	Gauge(key string, value int, labels map[string]string)
	Incr(key string, labels map[string]string)
}

// ConfigInterface is the interface for a core config
type ConfigInterface interface {
	GetLogger() *zap.SugaredLogger
	Workers() int
	GetSocketClient() *socketmode.Client
}
