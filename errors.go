package metrics

import "errors"

var (
	ErrTransportInactive = errors.New("transport is not running")
	ErrCollectorInactive = errors.New("collector is not running")
)
