package amqpx

import (
	"errors"
	"fmt"
)

var (
	// ErrInvalidConnectionsPoolCapacity occurs when the defined connections pool's capacity is invalid .
	ErrInvalidConnectionsPoolCapacity = fmt.Errorf("invalid connections pool capacity")

	// ErrNoConnectionAvailable occurs when the connections pool's has no healthy connections.
	ErrNoConnectionAvailable = fmt.Errorf("no connection available")

	// ErrClientClosed occurs when operating on a closed client.
	ErrClientClosed = fmt.Errorf("client is closed")

	// ErrBrokerURIRequired occurs when a dialer has no broker URI.
	ErrBrokerURIRequired = fmt.Errorf("broker URI is required")

	// ErrObserverRequired occurs when given observer is empty.
	ErrObserverRequired = fmt.Errorf("an observer instance is required")

	// ErrInvalidRetryDuration occurs when the defined retry duration is invalid .
	ErrInvalidRetryDuration = fmt.Errorf("invalid retry duration")

	// ErrOpenConnection occurs when a new connection cannot be open.
	ErrOpenConnection = errors.New("cannot open a new connection")

	// ErrOpenChannel occurs when a new channel cannot be open.
	ErrOpenChannel = errors.New("cannot open a new channel")

	// ErrRetryExceeded occurs when a retrier exceeded its attempts.
	ErrRetryExceeded = errors.New("retry exceeded")
)
