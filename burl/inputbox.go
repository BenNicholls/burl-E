package burl

import (
	"strconv"

	"github.com/veandco/go-sdl2/sdl"
)

//Inputbox is a textbox designed for user input of text, complete with mighty blinking cursor.
//TODO: String longer than size of Textbox. How hard could that be????
type Inputbox struct {
	Textbox
	cursorAnimation *BlinkCharAnimation
}

func NewInputbox(w, h, x, y, z int, bord bool) *Inputbox {
	ib := &Inputbox{*NewTextbox(w, h, x, y, z, bord, false, ""), NewBlinkCharAnimation(0, 0, 0, 20)}
	return ib
}

//Inserts a character/string s.
func (ib *Inputbox) Insert(s string) {
	if len(ib.text)+len(s) > ib.width*ib.height*2 {
		return
	}
	ib.ChangeString(ib.text + s)
}

//Actually more of a backspace action.
func (ib *Inputbox) Delete() {
	switch len(ib.text) {
	case 0:
		return
	default:
		ib.ChangeString(ib.text[:len(ib.text)-1])
	}
}

func (ib *Inputbox) Reset() {
	ib.ChangeString("")
}

//takes a key representing a letter and inserts.
func (ib *Inputbox) InsertText(key rune) {
	if !ValidText(key) {
		return
	}
	s := strconv.QuoteRuneToASCII(key)
	s, _ = strconv.Unquote(s)
	ib.Insert(s)
}

func (ib Inputbox) GetText() string {
	return ib.text
}

func (ib *Inputbox) ToggleFocus() {
	ib.focused = !ib.focused
	ib.cursorAnimation.Toggle()
}

func (ib *Inputbox) HandleKeypress(key sdl.Keycode) {
	switch key {
	case sdl.K_BACKSPACE:
		ib.Delete()
	case sdl.K_SPACE:
		ib.Insert(" ")
	default:
		ib.InsertText(rune(key))
	}
}

//TODO: Fix cursor for if inputbox has centered text. For now, just don't do that (looks silly anyways)
func (ib *Inputbox) Render() {
	if ib.visible {
		ib.Textbox.Render()
		ib.cursorAnimation.Tick()
		ib.cursorAnimation.x = len(ib.text) / 2
		ib.cursorAnimation.SetCharNum(len(ib.text) % 2)
		ib.cursorAnimation.Render(ib.x, ib.y, ib.z)
	}
}
