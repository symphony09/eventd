package eventd_test

import (
	"fmt"
	"math"
	"testing"
	"time"

	. "github.com/symphony09/eventd"
)

func TestEventBus(t *testing.T) {
	bus := new(EventBus[int])

	cancelDivide, err := bus.Subscribe(func(event string, i int) bool {
		fmt.Println("divide", event, i, "to", i/2)
		return true
	}, On("put"))

	if err != nil {
		t.Error(err)
	}

	cancelCheck, _ := bus.Subscribe(func(event string, i int) bool {
		fmt.Println("check", event, i)
		if i != 0 {
			return true
		} else {
			bus.Emit("check failed", i)
			return false
		}
	}, On("put"), Weight(math.MaxInt))

	_, _ = bus.Subscribe(func(event string, i int) bool {
		fmt.Println("log", event, i)
		return true
	}, On(".*"), Async)

	bus.Emit("put", 0)
	bus.Emit("put", 2)

	time.Sleep(time.Millisecond * 10)

	cancelDivide()
	cancelCheck()

	bus.Emit("put", 4)

	time.Sleep(time.Millisecond * 10)
}
