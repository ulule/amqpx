package amqpx

import (
	"github.com/cenk/backoff"
)

type exponentialRetrier struct {
	backoff backoff.BackOff
}

func newExponentialRetrier(opts retryOptions) retrier {
	bo := backoff.NewExponentialBackOff()
	bo.InitialInterval = defaultRetryInitialInterval
	bo.MaxInterval = defaultRetryMaxInterval
	bo.MaxElapsedTime = defaultRetryMaxElapsedTime

	if opts.exponential.initialInterval != 0 {
		bo.InitialInterval = opts.exponential.initialInterval
	}

	if opts.exponential.maxInterval != 0 {
		bo.MaxInterval = opts.exponential.maxInterval
	}

	if opts.exponential.maxElapsedTime != 0 {
		bo.MaxElapsedTime = opts.exponential.maxElapsedTime
	}

	return &exponentialRetrier{backoff: bo}
}

func (r exponentialRetrier) retry(handler func() error) error {
	r.backoff.Reset()
	return backoff.Retry(handler, r.backoff)
}

var _ retrier = (*exponentialRetrier)(nil)