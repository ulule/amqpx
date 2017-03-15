package amqpx

type singlePool struct {
	dialer Dialer
}

func NewSinglePool(dialer Dialer) (Pooler, error) {
	pool := &singlePool{
		dialer: dialer,
	}

	return pool, nil
}

func (s *singlePool) Get() (Connector, error) {
	return s.dialer()
}

func (s *singlePool) Close() error {
	return nil
}

func (s *singlePool) Length() int {
	return 1
}
