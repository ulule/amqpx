package amqpx_test

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ulule/amqpx"
)

var (
	brokerURI = "amqp://guest:guest@127.0.0.1:5672/amqpx"
)

func NewClient(options ...amqpx.Option) (amqpx.Client, error) {
	dialer, err := amqpx.SimpleDialer(brokerURI)
	if err != nil {
		return nil, err
	}

	return amqpx.New(dialer, options...)
}

type Runner struct {
	mutex   sync.Mutex
	require *require.Assertions
}

func NewRunner(t *testing.T) *Runner {
	return &Runner{
		require: require.New(t),
	}
}

func (e *Runner) NoError(err error, msgAndArgs ...interface{}) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.require.NoError(err, msgAndArgs...)
}

func (e *Runner) Error(err error, msgAndArgs ...interface{}) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.require.Error(err, msgAndArgs...)
}

func (e *Runner) Nil(object interface{}, msgAndArgs ...interface{}) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.require.Nil(object, msgAndArgs...)
}

func (e *Runner) NotNil(object interface{}, msgAndArgs ...interface{}) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.require.NotNil(object, msgAndArgs...)
}

func (e *Runner) IsType(expectedType interface{}, object interface{}, msgAndArgs ...interface{}) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.require.IsType(expectedType, object, msgAndArgs...)
}

func (e *Runner) Equal(expected interface{}, actual interface{}, msgAndArgs ...interface{}) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.require.Equal(expected, actual, msgAndArgs...)
}

func (e *Runner) True(value bool, msgAndArgs ...interface{}) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.require.True(value, msgAndArgs...)
}

func (e *Runner) False(value bool, msgAndArgs ...interface{}) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.require.False(value, msgAndArgs...)
}
