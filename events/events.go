package events

// Basic events
const (
	RequestReceived = iota
	BeforeHandler
	AfterHandler
)

type Event struct {
	Event int
}
