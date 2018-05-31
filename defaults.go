package amqpx

import "time"

// Retry
const (
	defaultRetryInitialInterval = 100 * time.Millisecond
	defaultRetryMaxInterval     = 32 * time.Second
	defaultRetryMaxElapsedTime  = 7 * time.Minute
)

// Dialer
var (
	defaultDialerTimeout   = 30 * time.Second
	defaultDialerHeartbeat = 10 * time.Second
)

// Pooler
const (
	// DefaultConnectionsCapacity is default connections pool capacity.
	DefaultConnectionsCapacity = 10
)
