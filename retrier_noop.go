package amqpx

type noopRetrier struct{}

func newNoopRetrier() retrier {
	return &noopRetrier{}
}

func (r noopRetrier) retry(handler func() error) error {
	return handler()
}

var _ retrier = (*noopRetrier)(nil)
