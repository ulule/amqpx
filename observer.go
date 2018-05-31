package amqpx

// Observer is an event collector.
type Observer interface {
	// OnError is called when an error occurs.
	OnError(err error)

	// OnClose is called when a close/shutdown error occurs.
	OnClose(err error)
}

// A defaultObserver is a no-op implementation of Observer interface.
type defaultObserver struct{}

func (defaultObserver) OnError(err error) {}
func (defaultObserver) OnClose(err error) {}

var _ Observer = (*defaultObserver)(nil)
