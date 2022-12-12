package metrics

import (
	"runtime"
	"sync"
	"time"

	"github.com/blacked/go-zabbix"
)

type zabbixTransport struct {
	sender      *zabbix.Sender
	logger      ErrorLogger
	metricsHost string
	collector   Collector
	interval    uint64
	run         sync.Once
	done        chan bool
}

// NewZabbix creates new metrics transport for zabbix server.
func NewZabbix(collector Collector, sender *zabbix.Sender, metricsHost string, interval uint64, logger ErrorLogger) Transport {
	t := &zabbixTransport{
		sender:      sender,
		metricsHost: metricsHost,
		collector:   collector,
		interval:    interval,
		logger:      logger,
	}
	runtime.SetFinalizer(t, StoppableFinalizer)
	return t
}

// Run zabbix transport
func (t *zabbixTransport) Run() {
	t.run.Do(func() {
		t.done = make(chan bool, 1)
		go func() {
			for {
				select {
				case <-t.done:
					return
				case <-time.After(time.Duration(t.interval) * time.Second):
					if err := t.Send(); err != nil {
						t.logger.Errorf("cannot send metrics to Zabbix: %v", err)
					}
				}
			}
		}()
	})
}

// Send metrics to zabbix
func (t *zabbixTransport) Send() error {
	metrics := t.collector.Metrics()

	data := make([]*zabbix.Metric, len(metrics))
	for i, m := range metrics {
		data[i] = zabbix.NewMetric(t.metricsHost, m.Name, m.Value)
	}
	_, err := t.sender.Send(zabbix.NewPacket(data))

	return err
}

// Stop zabbix transport
func (t *zabbixTransport) Stop() error {
	if t.done == nil {
		return ErrTransportInactive
	}
	t.done <- true
	close(t.done)
	return nil
}
