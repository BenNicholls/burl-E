package ui

//Buttons are textboxes that can fire an event when "pressed". Event goes into the ui.EventStream
//Currently the only button in the game uses a hardcoded PulseAnimation, this should probably be modifiable.
type Button struct {
	Textbox
	press      *Event
	focusPulse *PulseAnimation
}

//Creates a new button. Defaults to non-focussed state.
func NewButton(w, h, x, y, z int, bord, cent bool, txt string) *Button {
	p := NewPulseAnimation(0, 0, w, h, 20, 0, true)
	return &Button{*NewTextbox(w, h, x, y, z, bord, cent, txt), nil, p}
}

//register an event to fire when the button is pressed
func (b *Button) Register(e *Event) {
	b.press = e
}

//fires the registered event
func (b Button) Press() {
	if b.press != nil {
		EventStream <- b.press
	}
}

func (b *Button) ToggleFocus() {
	b.focused = !b.focused
	b.focusPulse.Toggle()
}

func (b Button) Render(offset ...int) {
	if b.visible {
		offX, offY, offZ := processOffset(offset)

		b.Textbox.Render(offX, offY, offZ)
		b.focusPulse.Tick()
		b.focusPulse.Render(b.x+offX, b.y+offY, b.z+offZ)
	}
}
