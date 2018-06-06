package amqpx

import "time"

// Retry
const (
	DefaultRetryInitialInterval = 100 * time.Millisecond
	DefaultRetryMaxInterval     = 32 * time.Second
	DefaultRetryMaxElapsedTime  = 7 * time.Minute
)

// Dialer
var (
	DefaultDialerTimeout   = 30 * time.Second
	DefaultDialerHeartbeat = 10 * time.Second
)

// Pooler
const (
	// DefaultConnectionsCapacity is default connections pool capacity.
	DefaultConnectionsCapacity = 10
)
