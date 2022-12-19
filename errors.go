package metrics

import "errors"

var (
	ErrNilSender         = errors.New("sender is nil")
	ErrTransportInactive = errors.New("transport is not running")
	ErrCollectorInactive = errors.New("collector is not running")
)
