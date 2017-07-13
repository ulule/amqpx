package amqpx

import (
	"fmt"
)

var (
	// ErrObserverRequired occurs when given observer is empty.
	ErrObserverRequired = fmt.Errorf("an observer instance is required")
)

// Observer is an event collector.
type Observer interface {
	// OnError is called when an error occurs.
	OnError(err error)
	// OnClose is called when a close/shutdown error occurs.
	OnClose(err error)
}

// WithObserver will configure client with given observer.
func WithObserver(observer Observer) Option {
	return option(func(options *clientOptions) error {
		if observer == nil {
			return ErrObserverRequired
		}

		options.observer = observer
		return nil
	})
}

// A defaultObserver is a no-op implementation of Observer interface.
type defaultObserver struct{}

func (defaultObserver) OnError(err error) {}

func (defaultObserver) OnClose(err error) {}
