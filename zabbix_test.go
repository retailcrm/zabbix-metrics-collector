package metrics

import (
	"testing"

	"github.com/blacked/go-zabbix"
	"github.com/stretchr/testify/assert"
)

func TestZabbixTransport_RunAndStop(t *testing.T) {
	c := NewMemoryCollector(10)
	s := zabbix.NewSender("127.0.0.1", 8081)

	transport := NewZabbix(c, s, "host", 10, DefaultLogger)

	transport.Run()
	transport.Run()

	assert.NoError(t, transport.Stop())
	assert.Error(t, transport.Stop())
}
