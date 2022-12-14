package metrics

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// MemoryCollector collects information about memory usage and goroutines count.
type MemoryCollector struct {
	done     chan bool
	memStat  runtime.MemStats
	interval uint64
	run      sync.Once
	memMutex sync.Mutex
}

// NewMemoryCollector creates new instance of AppCollector.
func NewMemoryCollector(interval uint64) *MemoryCollector {
	m := &MemoryCollector{interval: interval}
	runtime.SetFinalizer(m, StoppableFinalizer)
	return m
}

// Metrics fetches and clears metrics.
func (c *MemoryCollector) Metrics() []Metric {
	c.memMutex.Lock()
	memAlloc := c.memStat.Alloc
	memSys := c.memStat.Sys
	c.memMutex.Unlock()

	return []Metric{
		{
			Name:  "app.goroutine.count",
			Value: fmt.Sprintf("%d", runtime.NumGoroutine()),
		},
		{
			Name:  "app.mem.alloc",
			Value: fmt.Sprintf("%d", memAlloc),
		},
		{
			Name:  "app.mem.sys",
			Value: fmt.Sprintf("%d", memSys),
		},
	}
}

// Run collection process.
func (c *MemoryCollector) Run() {
	c.run.Do(func() {
		c.done = make(chan bool, 1)
		go func() {
			for {
				if c.done == nil {
					return
				}

				select {
				case <-c.done:
					return
				default:
					c.memMutex.Lock()
					runtime.ReadMemStats(&c.memStat)
					c.memMutex.Unlock()
					time.Sleep(time.Duration(c.interval) * time.Second)
				}
			}
		}()
	})
}

// Stop collection process.
func (c *MemoryCollector) Stop() error {
	if c.done == nil {
		return ErrCollectorInactive
	}
	c.done <- true
	close(c.done)
	c.done = nil
	return nil
}
