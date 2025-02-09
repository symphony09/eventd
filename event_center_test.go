package eventd_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/symphony09/eventd"
)

func TestEventCenter(t *testing.T) {
	center := new(eventd.EventCenter)

	bus1 := new(eventd.EventBus[string])
	bus2 := new(eventd.EventBus[int])

	center.AddListener(bus1, bus2)

	bus1.Subscribe(func(event, obj string) bool {
		if obj == "" {
			fmt.Println("got empty string")
		} else {
			fmt.Printf("got string: %s\n", obj)
		}

		return true
	}, eventd.On("test"))

	bus2.Subscribe(func(event string, obj int) bool {
		if obj == 0 {
			fmt.Println("got zero int")
		} else {
			fmt.Printf("got int: %d\n", obj)
		}

		return true
	}, eventd.On("test"))

	center.Broadcast("test", nil, 1, "a")

	time.Sleep(time.Millisecond * 10)
}
