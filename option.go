package eventd

type SubscribeOption struct {
	Events []string
	Weight int
	Async  bool
}

type Op func(option *SubscribeOption)

var On = func(events ...string) Op {
	return func(option *SubscribeOption) {
		option.Events = append(option.Events, events...)
	}
}

var Weight = func(w int) Op {
	return func(option *SubscribeOption) {
		option.Weight = w
	}
}

var Async Op = func(option *SubscribeOption) {
	option.Async = true
}
