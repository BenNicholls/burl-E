package delveengine

import "github.com/veandco/go-sdl2/sdl"

type GameModer interface {
	Update() (error, GameModer)
	Render()
	HandleKeypress(sdl.Keycode)
}
