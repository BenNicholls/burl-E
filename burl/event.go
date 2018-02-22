package burl

//event IDs used by burl internally
const (
	NONE            int = iota
	UPDATE_UI_EVENT     //signifies some UI needs to be updated
	MAX_EVENTS
)

//Event is the basic unit of the EventStream.
type Event struct {
	ID      int
	Message string
}

//EventStream is the queue of UI Events to be (optionally) consumed by the application.
var EventStream chan *Event

//Sometimes we end up putting like 10000 UI_UPDATEs in the queue and break it, so we keep track of
//which UPDATE_UI_EVENTs we psuh into the queue to ensure we don't duplicate them. They are removed from
//the map once popped. There never needs to be more than one since the events SHOULD be being consumed
//once per frame.
var UIEvents map[string]bool

//Allocate event buffer. 1000 events should be enough, right???
func init() {
	EventStream = make(chan *Event, 1000)
	UIEvents = make(map[string]bool, 100)
}

//Emits an event into the EventStream. If the stream is full we flush the whole buffer.
//TODO: Is flushing the buffer a little barbaric? We could maybe just consume half of them or something.
func PushEvent(id int, m string) {
	if len(EventStream) == cap(EventStream) {
		ClearEvents()
		LogError("UI Eventstream limit reached! FLUSHY FLUSHY.")
	}

	if id == UPDATE_UI_EVENT {
		if UIEvents[m] {
			return
		}
		UIEvents[m] = true
	}

	EventStream <- &Event{id, m}
}

//Reallocates the eventstream.
func ClearEvents() {
	EventStream = make(chan *Event, 1000)
}

//Grabs a UI event from the stream for consumption purposes.
func PopEvent() *Event {
	if len(EventStream) > 0 {
		e := <-EventStream

		if e.ID == UPDATE_UI_EVENT {
			UIEvents[e.Message] = false
		}

		return e
	} else {
		return nil
	}
}
