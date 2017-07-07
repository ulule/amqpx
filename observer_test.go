package amqpx_test

import (
	"fmt"
)

type TestObserver struct{}

func (TestObserver) OnError(err error) {
	fmt.Println(err)
}

func (TestObserver) OnClose(err error) {
	fmt.Println(err)
}
