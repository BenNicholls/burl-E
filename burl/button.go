package burl

//Buttons are textboxes that can fire an event when "pressed". Event goes into the EventStream
type Button struct {
	Textbox
	press      *Event
	PressPulse *PulseAnimation //animation that plays when pressed. TODO: make this modifiable. not always going to want a pulseanimation.
}

//Creates a new button. Defaults to non-focused state.
//TODO: some kind of ID system so we can include an ID with pressed events??
func NewButton(w, h, x, y, z int, bord, cent bool, txt string) *Button {
	p := NewPulseAnimation(0, 0, w, h, 20, 1, false)
	return &Button{*NewTextbox(w, h, x, y, z, bord, cent, txt), nil, p}
}

//register an event to fire when the button is pressed
func (b *Button) Register(e *Event) {
	b.press = e
}

//fires the registered event, plays press animation.
func (b *Button) Press() {
	b.PressPulse.Activate()
	if b.press != nil {
		PushEvent(b.press)
	}
}

func (b *Button) Render(offset ...int) {
	if b.visible {
		offX, offY, offZ := processOffset(offset)

		b.Textbox.Render(offX, offY, offZ)
		if b.PressPulse.enabled {
			b.PressPulse.Tick()
			b.PressPulse.Render(b.x+offX, b.y+offY, b.z+offZ)
			if b.PressPulse.IsFinished() && b.press != nil {
				PushEvent(NewEvent(ANIMATION_DONE, b.press.Message))
			}
		}
	}
}
