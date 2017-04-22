package ui

import "github.com/bennicholls/delvetown/console"

//Animations! Animations have a initalized state, then they evolve for a time, and then are garbaged.
//Some animations will naturally be specific to some ui elements and not others, those limitations
//will have to be listed here as best we can.
//Remember that animations always start DISABLED (enabled = false) and must be activated manually.
type Animator interface {
	Tick()
	Render(offset ...int)
	Toggle()
}

//BlinkAnimation inverts the fore/back colours of a single cell. Speed controls frequency.
type BlinkAnimation struct {
	tick    int
	speed   int //number of frames between blinks
	x, y    int //position (possibly relative to element or container)
	enabled bool
}

func NewBlinkAnimation(x, y, speed int) *BlinkAnimation {
	return &BlinkAnimation{0, speed, x, y, false}
}

func (ba *BlinkAnimation) Toggle() {
	ba.enabled = !ba.enabled
	ba.tick = 0
}

func (ba *BlinkAnimation) Tick() {
	if ba.enabled {
		ba.tick++
	}
}

func (ba *BlinkAnimation) Render(offset ...int) {
	if ba.enabled {
		if ba.tick%(ba.speed*2) < ba.speed {
			offX, offY, offZ := processOffset(offset)
			console.Invert(ba.x+offX, ba.y+offY, offZ)
		}
	}
}

//PulseAnimation make a rectangular area pulse with colour.
//TODO: Support for fun colours!!!!!!!!!!YES!!!!
type PulseAnimation struct {
	tick       int
	dur        int //duration of a pulse
	num        int //number of pulses to do
	x, y, w, h int
	enabled    bool
	repeat     bool
	done       bool
}

func NewPulseAnimation(x, y, w, h, dur, num int, repeat bool) *PulseAnimation {
	return &PulseAnimation{0, dur, num, x, y, w, h, false, repeat, false}
}

func (pa *PulseAnimation) Toggle() {
	pa.enabled = !pa.enabled
	if pa.repeat {
		pa.tick = 0
	}
}

func (pa *PulseAnimation) Tick() {
	if pa.enabled {
		pa.tick++
		if !pa.repeat {
			if pa.tick == pa.dur*pa.num {
				pa.done = true
			}
		}
	}
}

func (pa *PulseAnimation) Render(offset ...int) {
	if pa.enabled {
		//interpolate for correct pulse colour

		n := pa.tick % pa.dur
		if n > pa.dur/2 {
			n = pa.dur - n
		}
		c := int(255 * (float32(n) / float32((pa.dur / 2))))

		offX, offY, _ := processOffset(offset)

		for i := 0; i < pa.w*pa.h; i++ {
			console.ChangeBackColour(pa.x+offX+i%pa.w, pa.y+offY+i/pa.w, console.MakeColour(c, c, c))
		}
	}
}
