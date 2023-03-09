package eventd

type Listener interface {
	Notify(event string, obj any)
}
