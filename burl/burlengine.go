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
		PushEvent(NewEvent(CHANGE_STATE, ""))
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
		return console, err
	} else {
		return nil, err
	}

}

//The Big Enchelada! This is the gameloop that runs everything. Make sure to run
//burl.InitState() and burl.InitConsole before beginning the game!
func GameLoop() error {
	//TODO: implement that horrible thread job queue thing from the go-sdl2 package
	runtime.LockOSThread() //fixes some kind of go-sdl2 based thread release bug.

	if !console.Ready {
		return errors.New("Console not set up. Run burl.InitConsole() before starting game loop!")
	}

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
					gameState.HandleKeypress(t.Keysym.Sym)
				}
			}
		}

		gameState.Update()

		//serve events to application for handling
		for e := PopEvent(); e != nil; e = PopEvent() {
			gameState.HandleEvent(e)
		}

		//TODO: get console.Render() running in another thread (i think this is a good idea... maybe?)
		gameState.Render()
		console.Render() //should this come after the burl events are processed??

		//process burl-handled events
		for e := popInternalEvent(); e != nil; e = popInternalEvent() {
			switch e.ID {
			case QUIT_EVENT:
				gameState.Shutdown()
				running = false
			case CHANGE_STATE:
				gameState.Shutdown()
				gameState = nextState
				nextState = nil
			}
		}
	}

	log.Close()

	return nil
}

//Defines a game state (level, menu, anything that can take input, update itself, render to screen.)
type State interface {
	HandleKeypress(sdl.Keycode)
	Update()
	HandleEvent(*Event) //called for each event in the stream, every frame
	Render()
	GetTick() int
	Shutdown() //called on program exit
}

//base state object, compose states around this if you want
type BaseState struct {
	tick int //update ticks since init
}

func (b BaseState) GetTick() int {
	return b.tick
}

func (b BaseState) HandleKeypress(key sdl.Keycode) {

}

func (b *BaseState) Update() {
	b.tick++
}

func (b BaseState) Render() {

}

func (b BaseState) Shutdown() {

}

func (b BaseState) HandleEvent(e *Event) {

}
