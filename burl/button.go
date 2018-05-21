package burl

import (
	"github.com/veandco/go-sdl2/sdl"
)

//Buttons are textboxes that can fire an event when "pressed". Event goes into the EventStream
type Button struct {
	Textbox
	press      *Event
	PressPulse *PulseAnimation //animation that plays when pressed. TODO: make this modifiable. not always going to want a pulseanimation.
}

//Creates a new button. Defaults to non-focused state.
func NewButton(w, h, x, y, z int, bord, cent bool, txt string) *Button {
	p := NewPulseAnimation(0, 0, z, w, h, 20, 1, false)
	return &Button{*NewTextbox(w, h, x, y, z, bord, cent, txt), nil, p}
}

//register an event to fire when the button is pressed. can be anything. use
//registerPressEvent() to register a normal press event. If no event is regsitered,
//a default press event is sent on press.
func (b *Button) Register(e *Event) {
	b.press = e
}

//quick way to register a press event
func (b *Button) RegisterPressEvent(m string) {
	b.Register(NewUIEvent(EV_BUTTON_PRESS, m, b))
}

//fires the registered event, plays press animation.
func (b *Button) Press() {
	b.PressPulse.Activate()
	if b.press != nil {
		PushEvent(b.press)
	} else {
		PushEvent(NewUIEvent(EV_BUTTON_PRESS, "", b))
	}
}

func (b *Button) HandleKeypress(key sdl.Keycode) {
	if key == sdl.K_RETURN {
		b.Press()
	}
}

func (b *Button) Render() {
	if b.visible {
		b.Textbox.Render()
		if b.PressPulse.enabled {
			b.PressPulse.Tick()
			b.PressPulse.Render(b.x, b.y, b.z)
			if b.PressPulse.IsFinished() {
				if b.press != nil {
					PushEvent(NewUIEvent(EV_ANIMATION_DONE, b.press.Message, b))
				} else {
					PushEvent(NewUIEvent(EV_ANIMATION_DONE, "", b))
				}
			}
		}
	}
}
