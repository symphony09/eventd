package eventd

type EventCenter struct {
	listeners []Listener
}

func (center *EventCenter) AddListerner(listeners ...Listener) {
	center.listeners = append(center.listeners, listeners...)
}

func (center *EventCenter) Broadcast(event string, objects ...any) {
	for _, listener := range center.listeners {
		for _, obj := range objects {
			listener.Notify(event, obj)
		}
	}
}
