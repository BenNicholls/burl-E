package ui

import "github.com/bennicholls/delvetown/util"

//UIElem is the basic definition for all UI elements.
type UIElem interface {
	Render(offset ...int)
	Dims() (w int, h int)
	Pos() (x int, y int, z int)
	SetTitle(title string)
	ToggleVisible()
	SetVisibility(v bool)
	MoveTo(x, y, z int)
}

//Helper funtion for unpacking optional offsets passed to UI render functions. Required to allow for nesting of elements.
func processOffset(offset []int) (x, y, z int) {
	x, y, z = 0, 0, 0
	if len(offset) >= 2 {
		x, y = offset[0], offset[1]
		if len(offset) == 3 {
			z = offset[2]
		}
	}
	return
}

//event IDs
const (
	NONE     int = iota
	ACTIVATE     //used for buttons I guess?
	CHANGE       //used when a UIelem is changed
)

//Event is the basic unit of the EventStream, which is used to record input and interaction events.
type Event struct {
	Caller  UIElem
	ID      int
	Message string
}

//EventStream is the queue of UI Events to be (optionally) consumed by the application.
var EventStream chan *Event

//Allocate event buffer. TODO: is 100 overkill? Not enough? Test this once it is used ever.
func init() {
	EventStream = make(chan *Event, 100)
}

//Emits an event into the EventStream. If the stream is full we flush the whole buffer.
//TODO: Is flushing the buffer a little barbaric? We could maybe just consume half of them or something.
func PushEvent(c UIElem, id int, m string) {
	if len(EventStream) == cap(EventStream) {
		ClearEvents()
		util.LogError("UI Eventstream limit reached! FLUSHY FLUSHY.")
	}

	EventStream <- &Event{c, id, m}
}

//Reallocates the eventstream.
func ClearEvents() {
	EventStream = make(chan *Event, 100)
}

//Grabs a UI event from the stream for consumption purposes.
func PopEvent() *Event {
	if len(EventStream) > 0 {
		e := <-EventStream
		return e
	} else {
		return nil
	}
}
