package burl

import "errors"
import "runtime"
import "github.com/veandco/go-sdl2/sdl"

var console *Console
var gameState State
var nextState State

//Initializes the game State. Call before running the game loop.
func InitState(m State) {
	if gameState == nil {
		gameState = m
	}
}

//Tell burl to change from one state to another. This is done at the end of frame. Only the first
//call to this function will succeed per frame, subsequent calls evoke an error and are ignored.
func ChangeState(m State) {
	if nextState == nil {
		nextState = m
		PushEvent(NewEvent(EV_CHANGE_STATE, ""))
	} else {
		LogError("Multiple state changes detected in one frame!")
	}
}

//Initializes the console. Returns a pointer to the console so the user can
//manipulate it manually if they prefer. Returns nil if there was an error.
func InitConsole(w, h int, glyphPath, fontPath, title string) (*Console, error) {
	console = new(Console)
	err := console.Setup(w, h, glyphPath, fontPath, title)
	if err == nil {
		initDebugger()
		return console, err
	} else {
		return nil, err
	}
}

//Activate debugging capabilities. F10 will bring up the debug menu.
func Debug() {
	debug = true
}

//The Big Enchelada! This is the gameloop that runs everything. Make sure to run
//burl.InitState() and burl.InitConsole before beginning the game!
func GameLoop() error {
	runtime.LockOSThread() //sdl is inherently single-threaded.
	defer outputLogToDisk()

	if !console.Ready {
		return errors.New("Console not set up. Run burl.InitConsole() before starting game loop!")
	}
	defer console.Cleanup()

	if gameState == nil {
		return errors.New("No gameState initialized. Run burl.InitState() before starting game loop!")
	}

	var event sdl.Event
	running := true

	for running {
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				gameState.Shutdown()
				running = false
			case *sdl.WindowEvent:
				if t.Event == sdl.WINDOWEVENT_RESTORED {
					console.ForceRedraw()
				}
			case *sdl.KeyboardEvent:
				if t.Type == sdl.KEYDOWN {
					if debug && debugger.IsVisible() {
						debugger.HandleKeypress(t.Keysym.Sym)
					} else if t.Keysym.Sym == sdl.K_F10 {
						debugger.ToggleVisible()
					} else {
						if d := gameState.GetDialog(); d == nil {
							gameState.HandleKeypress(t.Keysym.Sym)
						} else {
							d.HandleKeypress(t.Keysym.Sym)
						}
					}
				}
			}
		}

		if debug && debugger.IsVisible() {
			debugger.Update()
		}

		if d := gameState.GetDialog(); d == nil {
			gameState.Update()
		} else {
			d.Update()
			if d.Done() {
				gameState.CloseDialog()
			}
		}

		//serve events to application for handling
		for e := PopEvent(); e != nil; e = PopEvent() {
			if d := gameState.GetDialog(); d == nil {
				gameState.HandleEvent(e)
			} else {
				d.HandleEvent(e)
			}
		}

		//TODO: get console.Render() running in another thread (i think this is a good idea... maybe?)
		if d := gameState.GetDialog(); d != nil {
			d.Render()
			if w := d.GetWindow(); w != nil {
				w.Render()
			}
		}
		gameState.Render()
		if w := gameState.GetWindow(); w != nil {
			w.Render()
		}

		if debug {
			debugger.Render()
		}

		console.Render() //should this come after the burl events are processed??

		//process burl-handled events
		for e := popInternalEvent(); e != nil; e = popInternalEvent() {
			switch e.ID {
			case EV_QUIT:
				gameState.Shutdown()
				running = false
			case EV_CHANGE_STATE:
				gameState.Shutdown()
				console.Clear()
				gameState = nextState
				nextState = nil
			}
		}
	}

	return nil
}

//Defines a game state (level, menu, anything that can take input, update itself, render to screen.)
type State interface {
	HandleKeypress(sdl.Keycode)
	Update()
	HandleEvent(*Event) //called for each event in the stream, every frame
	Render()
	GetTick() int
	GetWindow() *Container
	GetDialog() Dialog
	CloseDialog()
	Shutdown() //called on program exit
}

//Dialogs are states that can report when they are done.
type Dialog interface {
	State
	Done() bool
}

//base state object, compose states around this if you want
type StatePrototype struct {
	Tick   int //update ticks since init
	Window *Container
	dialog Dialog
}

func (b StatePrototype) GetTick() int {
	return b.Tick
}

func (b StatePrototype) HandleKeypress(key sdl.Keycode) {

}

func (b *StatePrototype) Update() {
	b.Tick++
}

func (b StatePrototype) Render() {

}

func (b StatePrototype) Shutdown() {

}

func (b StatePrototype) HandleEvent(e *Event) {

}

func (b *StatePrototype) InitWindow(bord bool) {
	w, h := console.Dims()
	x, y := 0, 0
	if bord {
		w, h, x, y = w-2, h-2, 1, 1
	}
	b.Window = NewContainer(w, h, x, y, 0, true)
}

func (b StatePrototype) GetWindow() *Container {
	return b.Window
}

func (b *StatePrototype) OpenDialog(d Dialog) {
	if b.dialog != nil {
		b.CloseDialog()
	}
	b.dialog = d
}

func (b StatePrototype) GetDialog() Dialog {
	return b.dialog
}

func (b *StatePrototype) CloseDialog() {
	b.dialog.GetWindow().ToggleVisible()
	b.dialog = nil
}
