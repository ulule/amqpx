package amqpx

import (
	"fmt"
)

var (
	// ErrInvalidConnectionsPoolCapacity occurs when the defined connections pool's capacity is invalid .
	ErrInvalidConnectionsPoolCapacity = fmt.Errorf("invalid connections pool capacity")

	// ErrInvalidRetryDuration occurs when the defined retry duration is invalid .
	ErrInvalidRetryDuration = fmt.Errorf("invalid retry duration")

	// ErrInvalidDialerTimeout occurs when the defined dialer timeout is invalid.
	ErrInvalidDialerTimeout = fmt.Errorf("invalid dialer timeout")

	// ErrInvalidDialerHeartbeat occurs when the defined dialer heartbeat is invalid.
	ErrInvalidDialerHeartbeat = fmt.Errorf("invalid dialer heartbeat")

	// ErrNoConnectionAvailable occurs when the connections pool's has no healthy connections.
	ErrNoConnectionAvailable = fmt.Errorf("no connection available")

	// ErrClientClosed occurs when operating on a closed client.
	ErrClientClosed = fmt.Errorf("client is closed")

	// ErrBrokerURIRequired occurs when a dialer has no broker URI.
	ErrBrokerURIRequired = fmt.Errorf("broker URI is required")

	// ErrObserverRequired occurs when given observer is empty.
	ErrObserverRequired = fmt.Errorf("an observer instance is required")

	// ErrLoggerRequired occurs when given logger is not set.
	ErrLoggerRequired = fmt.Errorf("a logger instance is required")
)

// Error Messages
const (
	ErrMessageCannotCreateDialer    = "cannot create a new dialer"
	ErrMessageCannotCreateClient    = "cannot create a new client"
	ErrMessageCannotOpenConnection  = "cannot open a new connection"
	ErrMessageCannotOpenChannel     = "cannot open a new channel"
	ErrMessageCannotCloseConnection = "cannot close connection"
	ErrMessageCannotCloseChannel    = "cannot close channel"
	ErrMessageDialTimeout           = "dialing remote address has timeout"
	ErrMessageReadTimeout           = "reading on socket has timeout"
	ErrMessageWriteTimeout          = "writing on socket has timeout"
)
