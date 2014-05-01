package events

// Basic events
const (
	RequestReceived = iota
	BeforeHandler
	AfterHandler
	AfterResponse
)

type Emitter struct {
	listeners map[int][]func(interface{},func())
}
func (e *Emitter) On(event int, handler func(interface{},func())) {
	if e.listeners == nil {
		e.listeners = make(map[int][]func(interface{},func()))
	}
	_, ok := e.listeners[event]
	if !ok {
		e.listeners[event] = make([]func(interface{},func()), 0)
	}
	e.listeners[event] = append(e.listeners[event], handler)
}
func (e *Emitter) Emit(session interface{}, event int, next func()) {
	if e.listeners == nil { 
		next()
		return 
	}
	c := make(chan int)
	n := 0
	if l, ok := e.listeners[event]; ok {
		if len(l) > 0 {
			for _, h := range l {
				n++
				go h(session, func() {
					c <- 1
				})
			}
		}
	}
	for n > 0 {
		i := <- c
		if i == 1 {
			n--
		}
	}
	next()
}
