package burl

//Buttons are textboxes that can fire an event when "pressed". Event goes into the ui.EventStream
type Button struct {
	Textbox
	press      *Event
	PressPulse *PulseAnimation //animation that plays when pressed. TODO: make this modifiable. not always going to want a pulseanimation.
}

//Creates a new button. Defaults to non-focused state.
func NewButton(w, h, x, y, z int, bord, cent bool, txt string) *Button {
	p := NewPulseAnimation(0, 0, w, h, 20, 1, false)
	return &Button{*NewTextbox(w, h, x, y, z, bord, cent, txt), nil, p}
}

//register an event to fire when the button is pressed
func (b *Button) Register(e *Event) {
	b.press = e
}

//fires the registered event, plays press animation.
func (b Button) Press() {
	b.PressPulse.Activate()
	if b.press != nil {
		EventStream <- b.press
	}
}

func (b *Button) ToggleFocus() {
	b.focused = !b.focused
}

func (b Button) Render(offset ...int) {
	if b.visible {
		offX, offY, offZ := processOffset(offset)

		b.Textbox.Render(offX, offY, offZ)
		b.PressPulse.Tick()
		b.PressPulse.Render(b.x+offX, b.y+offY, b.z+offZ)
	}
}
