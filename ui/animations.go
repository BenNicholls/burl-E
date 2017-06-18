package ui

import "github.com/bennicholls/burl/console"

//Animations! Animations have a initalized state, then they evolve for a time, and then are garbaged.
//Some animations will naturally be specific to some ui elements and not others, those limitations
//will have to be listed here as best we can.
//Remember that animations always start DISABLED (enabled = false) and must be activated manually.
type Animator interface {
	Tick()
	Render(offset ...int)
	Toggle()
	Activate() //Activates the animation. If it's already running, restarts it.
	IsFinished() bool
	Move(dx, dy int)
	MoveTo(x, y int)
}

type Animation struct {
	x, y    int
	tick    int
	enabled bool
	repeat  bool
	done    bool
}

func NewAnimation(x, y int, repeat bool) Animation {
	return Animation{x, y, 0, false, repeat, false}
}

func (a *Animation) Tick() {
	if a.enabled {
		a.tick++
	}
}

func (a Animation) Render() {

}

//Turns on animation. Does not restart animation (can be used as pause button).
func (a *Animation) Toggle() {
	a.enabled = !a.enabled
}

//plays the animation from the beginning. if it's already going, restarts it.
func (a *Animation) Activate() {
	a.tick = 0
	a.enabled = true
}

func (a Animation) IsFinished() bool {
	return a.done
}

func (a *Animation) Move(dx, dy int) {
	a.x += dx
	a.y += dy
}

func (a *Animation) MoveTo(x, y int) {
	a.x = x
	a.y = y
}

//BlinkCharAnimation draws a blinking cursor character. Speed controls frequency.
type BlinkCharAnimation struct {
	Animation
	speed int  //number of frames between blinks
	state bool //cursor shown or not shown
}

func NewBlinkCharAnimation(x, y, speed int) *BlinkCharAnimation {
	return &BlinkCharAnimation{NewAnimation(x, y, true), speed, true}
}

func (ba *BlinkCharAnimation) Tick() {
	if ba.enabled {
		ba.Animation.Tick()
		if ba.tick%ba.speed == 0 {
			ba.state = !ba.state
		}
	}
}

func (ba *BlinkCharAnimation) Activate() {
	if ba.enabled {
		ba.state = false
	}

	ba.Animation.Activate()
}

//charnum: 0 = left, 1 = right char
func (ba *BlinkCharAnimation) Render(charNum int, offset ...int) {
	if ba.enabled {
		offX, offY, offZ := processOffset(offset)
		if ba.state {
			console.ChangeChar(ba.x+offX, ba.y+offY, offZ, 31, charNum)
		} else {
			console.ChangeChar(ba.x+offX, ba.y+offY, offZ, 32, charNum)
		}
	}
}

//PulseAnimation make a rectangular area pulse with colour.
//TODO: Support for fun colours!!!!!!!!!!YES!!!!
type PulseAnimation struct {
	Animation
	dur  int //duration of a pulse
	num  int //number of pulses to do
	w, h int
}

func NewPulseAnimation(x, y, w, h, dur, num int, repeat bool) *PulseAnimation {
	return &PulseAnimation{NewAnimation(x, y, repeat), dur, num, w, h}
}

func (pa *PulseAnimation) Tick() {
	if pa.enabled {
		pa.Animation.Tick()
		if !pa.repeat {
			if pa.tick == pa.dur*pa.num+1 {
				pa.enabled = false
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

		offX, offY, offZ := processOffset(offset)

		for i := 0; i < pa.w*pa.h; i++ {
			console.ChangeBackColour(pa.x+offX+i%pa.w, pa.y+offY+i/pa.w, offZ, console.MakeColour(c, c, c))
		}
	}
}
