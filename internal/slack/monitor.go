package slack

type NoopMonitor struct{}

func (m *NoopMonitor) Gauge(key string, value int, labels map[string]string) {}
func (m *NoopMonitor) Incr(key string, labels map[string]string)             {}

func NewNoopMonitor() MonitorInterface {
	return new(NoopMonitor)
}
