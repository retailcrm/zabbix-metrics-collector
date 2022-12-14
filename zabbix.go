package metrics

import (
	"runtime"
	"sync"
	"time"

	"github.com/blacked/go-zabbix"
)

type zabbixTransport struct {
	logger      ErrorLogger
	sender      *zabbix.Sender
	done        chan bool
	metricsHost string
	collectors  []Collector
	interval    uint64
	run         sync.Once
	runMutex    sync.Mutex
}

// NewZabbix creates new metrics transport for zabbix server.
func NewZabbix(
	collectors []Collector, sender *zabbix.Sender, metricsHost string, interval uint64, logger ErrorLogger) Transport {
	t := &zabbixTransport{
		sender:      sender,
		metricsHost: metricsHost,
		collectors:  collectors,
		interval:    interval,
		logger:      logger,
	}
	runtime.SetFinalizer(t, StoppableFinalizer)
	return t
}

// WithCollector adds collector to the list of collectors that will be used for metrics collection.
func (t *zabbixTransport) WithCollector(col Collector) Transport {
	if t.collectors == nil {
		t.collectors = []Collector{col}
		return t
	}
	t.collectors = append(t.collectors, col)
	return t
}

// Run zabbix transport.
func (t *zabbixTransport) Run() {
	t.run.Do(func() {
		t.done = make(chan bool, 1)
		go func() {
			for {
				select {
				case <-t.done:
					t.runMutex.Lock()
					defer t.runMutex.Unlock()
					close(t.done)
					t.done = nil
					t.run = sync.Once{}
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

// Send metrics to zabbix.
func (t *zabbixTransport) Send() error {
	var total int
	metrics := make([][]Metric, len(t.collectors))

	for i, col := range t.collectors {
		rec := col.Metrics()
		total += len(rec)
		metrics[i] = rec
	}

	i := 0
	flattened := make([]*zabbix.Metric, total)
	namesMap := make(map[string]struct{}, total)
	for _, metricList := range metrics {
		for _, item := range metricList {
			if _, ok := namesMap[item.Name]; ok {
				continue
			}
			flattened[i] = zabbix.NewMetric(t.metricsHost, item.Name, item.Value)
			namesMap[item.Name] = struct{}{}
			i++
		}
	}

	_, err := t.sender.Send(zabbix.NewPacket(flattened))
	return err
}

// Stop zabbix transport.
func (t *zabbixTransport) Stop() error {
	t.runMutex.Lock()
	defer t.runMutex.Unlock()
	if t.done == nil {
		return ErrTransportInactive
	}
	t.done <- true
	return nil
}
