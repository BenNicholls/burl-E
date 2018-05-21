package burl

//Animations! Animations have a initalized state, then they evolve for a time, and then are garbaged.
//Some animations will naturally be specific to some ui elements and not others, those limitations
//will have to be listed here as best we can.
//Remember that animations always start DISABLED (enabled = false) and must be activated manually.
type Animator interface {
	Tick()
	Render(offX, offY, offZ int)
	Toggle()
	Activate() //Activates the animation. If it's already running, restarts it.
	IsFinished() bool
	Move(dx, dy, dz int)
	MoveTo(x, y int)
}

type Animation struct {
	x, y, z int
	tick    int
	enabled bool
	repeat  bool
	done    bool
}

func NewAnimation(x, y, z int, repeat bool) Animation {
	return Animation{x, y, z, 0, false, repeat, false}
}

func (a *Animation) Tick() {
	if a.enabled {
		a.tick++
	}
}

func (a Animation) Render(offX, offY, offZ int) {

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

func (a *Animation) Move(dx, dy, dz int) {
	a.x += dx
	a.y += dy
	a.z += dz
}

func (a *Animation) MoveTo(x, y int) {
	a.x = x
	a.y = y
}

//BlinkCharAnimation draws a blinking cursor character. Speed controls frequency.
type BlinkCharAnimation struct {
	Animation
	speed        int  //number of frames between blinks
	state        bool //cursor shown or not shown
	startCharNum int  //charNum to start drawing from
}

func NewBlinkCharAnimation(x, y, z, speed int) *BlinkCharAnimation {
	return &BlinkCharAnimation{NewAnimation(x, y, z, true), speed, true, 0}
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
		ba.state = true
	}

	ba.Animation.Activate()
}

//charnum: 0 = left, 1 = right char
func (ba *BlinkCharAnimation) SetCharNum(num int) {
	ba.startCharNum = num % 2
}

func (ba *BlinkCharAnimation) Render(offX, offY, offZ int) {
	if ba.enabled {
		if ba.state {
			console.ChangeForeColour(ba.x+offX, ba.y+offY, ba.z+offZ, COL_WHITE)
			console.ChangeChar(ba.x+offX, ba.y+offY, ba.z+offZ, 31, ba.startCharNum)
		} else {
			console.ChangeChar(ba.x+offX, ba.y+offY, ba.z+offZ, 32, ba.startCharNum)
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

func NewPulseAnimation(x, y, z, w, h, dur, num int, repeat bool) *PulseAnimation {
	return &PulseAnimation{NewAnimation(x, y, z, repeat), dur, num, w, h}
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

func (pa *PulseAnimation) Render(offX, offY, offZ int) {
	if pa.enabled {
		//interpolate for correct pulse colour
		n := pa.tick % pa.dur
		if n > pa.dur/2 {
			n = pa.dur - n
		}

		c := Lerp(0, 255, n, pa.dur/2)
		col := MakeOpaqueColour(c, c, c)

		for i := 0; i < pa.w*pa.h; i++ {
			console.ChangeBackColour(pa.x+i%pa.w+offX, pa.y+i/pa.w+offY, pa.z+offZ, col)
		}
	}
}
