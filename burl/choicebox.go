package burl

import (
	"math/rand"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	VERTICAL int = iota
	HORIZONTAL
)

//Choicebox is a textbox wherein once can cycle through some predefined choices. Only one choice
//is shown at a time. When focused, arrows appear to indicate which direction on the keyboard 
//cycles the choices (vertical/horizontal).
//NOTE: This is very similar to a text-only list. Someday maybe these should be combined?
type ChoiceBox struct {
	Textbox

	choices   []string
	curChoice int
	direction int //see direction consts above
}

//NewChoiceBox. Textbox parameters as normal. dir parameter is either HORIZONTAL or
//VERTICAL. Defaults to horizontal. Then accepts any number of strings as choices.
func NewChoiceBox(w, h, x, y, z int, bord bool, dir int, choices ...string) (cb *ChoiceBox) {
	cb = new(ChoiceBox)
	cb.Textbox = *NewTextbox(w, h, x, y, z, bord, true, "")
	cb.choices = choices

	if dir == VERTICAL {
		cb.direction = VERTICAL
	} else {
		cb.direction = HORIZONTAL
	}

	if len(cb.choices) == 0 {
		cb.ChangeText("---")
	} else {
		cb.ChangeText(cb.choices[cb.curChoice])
	}

	return
}

func (cb *ChoiceBox) AddChoice(c string) {
	cb.choices = append(cb.choices, c)
}

func (cb *ChoiceBox) Next() {
	cb.curChoice, _ = ModularClamp(cb.curChoice+1, 0, len(cb.choices)-1)
	cb.ChangeText(cb.choices[cb.curChoice])
	PushEvent(NewUIEvent(EV_LIST_CYCLE, "+", cb))
}

func (cb *ChoiceBox) Prev() {
	cb.curChoice, _ = ModularClamp(cb.curChoice-1, 0, len(cb.choices)-1)
	cb.ChangeText(cb.choices[cb.curChoice])
	PushEvent(NewUIEvent(EV_LIST_CYCLE, "-", cb))
}

func (cb *ChoiceBox) GetChoice() int {
	return cb.curChoice
}

func (cb *ChoiceBox) RandomizeChoice() {
	cb.curChoice = rand.Intn(len(cb.choices))
	cb.ChangeText(cb.choices[cb.curChoice])
}

func (cb *ChoiceBox) HandleKeypress(key sdl.Keycode) {
	if cb.direction == HORIZONTAL {
		switch key {
		case sdl.K_LEFT:
			cb.Prev()
		case sdl.K_RIGHT:
			cb.Next()
		}
	} else {
		switch key {
		case sdl.K_UP:
			cb.Prev()
		case sdl.K_DOWN:
			cb.Next()
		}
	}
}

func (cb *ChoiceBox) Render() {
	if cb.visible {
		cb.Textbox.Render()

		//draw choice cycling triangles
		if cb.direction == HORIZONTAL {
			console.ChangeCell(cb.x, cb.y+cb.height/2, cb.z, GLYPH_TRIANGLE_LEFT, COL_WHITE, COL_BLACK)
			console.ChangeCell(cb.x+cb.width-1, cb.y+cb.height/2, cb.z, GLYPH_TRIANGLE_RIGHT, COL_WHITE, COL_BLACK)
		} else {
			console.ChangeCell(cb.x+cb.width/2, cb.y-1, cb.z, GLYPH_TRIANGLE_UP, COL_WHITE, COL_BLACK)
			console.ChangeCell(cb.x+cb.width/2, cb.y+cb.height, cb.z, GLYPH_TRIANGLE_DOWN, COL_WHITE, COL_BLACK)
		}
	}
}
