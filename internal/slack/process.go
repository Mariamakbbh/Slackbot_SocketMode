package slack

import (
	"sync"
	"time"

	"github.com/slack-go/slack/socketmode"
	"go.uber.org/zap"
)

type SocketChannelConsumer struct {
	// goroutine control
	wg sync.WaitGroup

	// generic config
	workers     int
	groupID     string
	readTimeout time.Duration

	messageCh  chan socketmode.Event
	shutdownCh chan bool

	// external dependencies
	svc     ServiceInterface
	monitor MonitorInterface
	client  *socketmode.Client

	logger *zap.SugaredLogger
}

// consume takes care of running an infinite loop and run the service.Flow method
func (c *SocketChannelConsumer) consume(workerID int) {
	metricLabels := map[string]string{
		"service": "slack-socket",
	}

	c.logger.Infof("starting worker %v", workerID)
	for {
		select {
		case <-c.shutdownCh:
			c.logger.Infof("shutting down worker %v", workerID)
			c.wg.Done()
			return
		case msg := <-c.messageCh:
			err := c.svc.Flow(msg)

			if err != nil {
				c.monitor.Incr("events_errors", metricLabels)
				c.logger.Warn("Moving on...")
				continue
			}

			c.monitor.Incr("events_processed", metricLabels)
		}
	}
}

// read is the infinite loop for the kafka reading part of the consumer
func (c *SocketChannelConsumer) read() {
	metricLabels := map[string]string{
		"service": "slack-socket",
	}

	for {
		select {
		case <-c.shutdownCh:
			c.logger.Info("shutting down reader")
			c.wg.Done()
			return
		case msg := <-c.client.Events:
			c.messageCh <- msg
			c.monitor.Incr("events_received", metricLabels)
		default:
			continue
		}
	}
}

// Stop makes sure goroutines for read and consume are being gracefully stopped
func (c *SocketChannelConsumer) Stop() {
	// send reader shutdown
	c.shutdownCh <- true

	// send #workers shutdowns
	for i := 0; i < c.workers; i++ {
		c.shutdownCh <- true
	}

	c.wg.Wait()
	defer c.logger.Desugar().Sync()
}

// Start creates all the goroutines for read and consume
func (c *SocketChannelConsumer) Start() {

	for i := 0; i < c.workers; i++ {
		c.wg.Add(1)
		go c.consume(i + 1)
	}
	c.wg.Add(1)
	go c.read()
	go c.client.Run()
}

func NewSocketChannelConsumer(cfg ConfigInterface, svc ServiceInterface, monitor MonitorInterface) ConsumerInterface {
	consumer := new(SocketChannelConsumer)

	consumer.svc = svc
	consumer.monitor = monitor

	consumer.client = cfg.GetSocketClient()
	consumer.workers = cfg.Workers()
	consumer.logger = cfg.GetLogger()

	consumer.messageCh = make(chan socketmode.Event, consumer.workers)
	consumer.shutdownCh = make(chan bool)

	consumer.logger.Debug("config:", cfg)

	return consumer
}
