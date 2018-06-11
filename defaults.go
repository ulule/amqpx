package amqpx

import "time"

// Dialer default configuration.
var (
	DefaultDialerTimeout   = 30 * time.Second
	DefaultDialerHeartbeat = 10 * time.Second
)

// Pooler default configuration.
const (
	// DefaultConnectionsCapacity is default connections pool capacity.
	DefaultConnectionsCapacity = 10
)
