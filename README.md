## Zabbix Metrics Collector
[![Build Status](https://github.com/retailcrm/zabbix-metrics-collector/workflows/ci/badge.svg)](https://github.com/retailcrm/zabbix-metrics-collector/actions?query=workflow%3Aci)
[![Coverage](https://codecov.io/gh/retailcrm/zabbix-metrics-collector/branch/master/graph/badge.svg?logo=codecov&logoColor=white)](https://codecov.io/gh/retailcrm/zabbix-metrics-collector)
[![GitHub release](https://img.shields.io/github/release/retailcrm/zabbix-metrics-collector.svg?logo=github&logoColor=white)](https://github.com/retailcrm/zabbix-metrics-collector/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/retailcrm/zabbix-metrics-collector)](https://goreportcard.com/report/github.com/retailcrm/zabbix-metrics-collector)
[![GoLang version](https://img.shields.io/badge/go->=1.16-blue.svg?logo=go&logoColor=white)](https://golang.org/dl/)
[![pkg.go.dev](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/github.com/retailcrm/zabbix-metrics-collector/core)

This library provides easy-to-use wrapper for sending metrics to the Zabbix server.

Usage:
```go
package main

import (
	"github.com/blacked/go-zabbix"
	"github.com/retailcrm/zabbix-metrics-collector"
)

func main() {
    go metrics("app", "zabbix_server_docker", 10051)
    // some other logic here
}

func metrics(appHostZabbix, host string, port, intervalSeconds int) {
	sender := zabbix.NewSender(host, port)
	collector := metrics.NewMemoryCollector(intervalSeconds)
	proc := metrics.NewZabbix([]metrics.Collector{collector}, sender, appHostZabbix, intervalSeconds, metrics.DefaultLogger)
	proc.Run()
}
```

To collect other metrics you must implement `metrics.Collector` interface and collect relevant data there. The collector 
must be goroutine-safe.

You also can implement your own transport by using `metrics.Transport` interface. It may be useful in some cases. 
The included Zabbix transport will work just fine for the majority of use-cases.
