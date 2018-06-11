package amqpx_test

import (
	"math/rand"
	"os"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/ulule/amqpx"
)

var (
	invalidBrokerURI = "amqp://lionel:richie@127.0.0.1:5000/amqpx"
	brokerURI        = "amqp://guest:guest@127.0.0.1:5672/amqpx"
	brokerURIs       = []string{
		"amqp://guest:guest@127.0.0.1:5672/amqpx",
		"amqp://guest:guest@127.0.0.1:5673/amqpx",
		"amqp://guest:guest@127.0.0.1:5674/amqpx",
	}
)

const (
	goroutineNumber = 2000
	sigkillSleep    = 2 * time.Second
	dialerTimeout   = 15 * time.Second
	dialerHeartbeat = 15 * time.Second
)

func IsClusterMode() bool {
	switch os.Getenv("AMQPX_CLUSTER_MODE") {
	case "y", "Y", "yes", "true":
		return true
	default:
		return false
	}
}

func NewClient(options ...amqpx.ClientOption) (amqpx.Client, error) {
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

func (e *Runner) Contains(s interface{}, contains interface{}, msgAndArgs ...interface{}) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.require.Contains(s, contains, msgAndArgs...)
}

func TestClusterMode(t *testing.T) {
	if !IsClusterMode() {
		t.Skip()
	}

	is := NewRunner(t)

	dialer, err := amqpx.ClusterDialer(brokerURIs)
	is.NoError(err)
	is.NotNil(dialer)

	client, err := amqpx.New(dialer)
	is.NoError(err)
	is.NotNil(client)
	defer func() {
		is.NoError(client.Close())
	}()

	testClientExchange(is, client, "random.cluster")
}

func TestSimpleMode(t *testing.T) {
	is := NewRunner(t)

	dialer, err := amqpx.SimpleDialer(brokerURI)
	is.NoError(err)
	is.NotNil(dialer)

	client, err := amqpx.New(dialer)
	is.NoError(err)
	is.NotNil(client)
	defer func() {
		is.NoError(client.Close())
	}()

	testClientExchange(is, client, "random.simple")
}

func testClientExchange(is *Runner, client amqpx.Client, topic string) {
	messages := GenerateMessages()

	producer := &Producer{runner: is, messages: messages}
	producer.Start()
	for i := 0; i < 16; i++ {
		producer.NewEmitter(client, topic)
	}

	consumer := &Consumer{runner: is}
	consumer.Start()
	for i := 0; i < 16; i++ {
		consumer.NewReceiver(client, topic)
	}

	producer.Wait()
	consumer.Wait()

	is.Equal(len(messages), len(consumer.messages))
	sort.Strings(messages)
	sort.Strings(consumer.messages)
	for i := range messages {
		is.Equal(messages[i], consumer.messages[i])
	}
}

func testClientConcurrentAccess(is *Runner, client amqpx.Client, wg *sync.WaitGroup) {
	time.Sleep(time.Duration(rand.Intn(4000)) * time.Millisecond)

	channel, err := client.Channel()
	wg.Done()

	defer func() {
		if channel != nil {
			thr := channel.Close()
			_ = thr
		}
	}()

	if err != nil && errors.Cause(err) != amqpx.ErrClientClosed {
		is.NoError(err)
	}
}
