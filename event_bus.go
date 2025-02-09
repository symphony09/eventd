package eventd

import (
	"regexp"
	"sync"
)

type EventBus[T any] struct {
	eventRegexp map[string]*regexp.Regexp

	callbackChain map[Trigger]*Chain[CallBack[T]]

	mu sync.RWMutex
}

type Trigger struct {
	Event string
	Async bool
}

type CallBack[T any] func(event string, obj T) bool

func (bus *EventBus[T]) Subscribe(callback CallBack[T], ops ...Op) (cancel func(), err error) {
	bus.mu.Lock()
	defer bus.mu.Unlock()

	if bus.eventRegexp == nil {
		bus.eventRegexp = make(map[string]*regexp.Regexp)
	}

	if bus.callbackChain == nil {
		bus.callbackChain = make(map[Trigger]*Chain[CallBack[T]])
	}

	option := new(SubscribeOption)
	for _, op := range ops {
		op(option)
	}

	callbackNodes := make(map[Trigger]*ChainNode[CallBack[T]])

	for _, event := range option.Events {
		if bus.eventRegexp[event] == nil {
			if r, err := regexp.Compile(event); err != nil {
				return nil, err
			} else {
				bus.eventRegexp[event] = r
			}
		}

		trigger := Trigger{
			Event: event,
			Async: option.Async,
		}

		if bus.callbackChain[trigger] == nil {
			bus.callbackChain[trigger] = new(Chain[CallBack[T]])
		}

		callbackNodes[trigger] = bus.callbackChain[trigger].AddElem(callback, option.Weight)
	}

	return func() {
		for trigger, callbackNode := range callbackNodes {
			if chain := bus.callbackChain[trigger]; chain != nil {
				chain.RemoveNode(callbackNode)
			}
		}
	}, nil
}

func (bus *EventBus[T]) Emit(event string, object T) {
	bus.mu.RLock()
	defer bus.mu.RUnlock()

	if bus.eventRegexp == nil || bus.callbackChain == nil {
		return
	}

	var wg sync.WaitGroup

	for targetEvent, r := range bus.eventRegexp {
		wg.Add(1)

		go func(targetEvent string, r *regexp.Regexp) {
			defer wg.Done()

			if r.MatchString(event) {
				asyncTrigger := Trigger{
					Event: targetEvent,
					Async: true,
				}

				if chain := bus.callbackChain[asyncTrigger]; chain != nil {
					go func(chain *Chain[CallBack[T]]) {
						iterator := chain.Iteration()

						for iterator.Next() {
							callback := iterator.Get()

							callback(event, object)
						}
					}(chain)
				}

				syncTrigger := Trigger{
					Event: targetEvent,
					Async: false,
				}

				if chain := bus.callbackChain[syncTrigger]; chain != nil {
					iterator := chain.Iteration()

					for iterator.Next() {
						callback := iterator.Get()
						if !callback(event, object) {
							break
						}
					}
				}
			}
		}(targetEvent, r)
	}

	wg.Wait()
}

func (bus *EventBus[T]) Notify(event string, obj any) {
	if obj == nil {
		bus.Emit(event, *new(T))
	} else if v, ok := obj.(T); ok {
		bus.Emit(event, v)
	}
}
