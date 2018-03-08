package burl

import (
	"github.com/veandco/go-sdl2/sdl"
)

//UIElem is the basic definition for all UI elements.
type UIElem interface {
	Render(offset ...int)
	Dims() (w int, h int)
	Pos() (x int, y int, z int)
	SetTitle(title string)
	ToggleVisible()
	ToggleFocus()
	SetVisibility(v bool)
	IsVisible() bool
	IsFocused() bool
	MoveTo(x, y, z int)
	Rect() (int, int, int, int)
	CenterInConsole()
	SetTabID(id int)
	TabID() int
	HandleKeypress(key sdl.Keycode)
}

type UIElement struct {
	width, height int
	x, y, z       int
	bordered      bool
	title         string
	visible       bool
	focused       bool
	tabID         int //for to tab between elements in a container

	anims []Animator
}

func NewUIElement(width, height, x, y, z int, bord bool) UIElement {
	return UIElement{width, height, x, y, z, bord, "", true, false, 0, make([]Animator, 0, 20)}
}

//basic render function for all elements.
func (u *UIElement) Render(offset ...int) {
	if u.visible {
		offX, offY, offZ := processOffset(offset)

		if u.bordered {
			console.DrawBorder(u.x+offX, u.y+offY, u.z+offZ, u.width, u.height, u.title, u.focused)
		}

		for i, _ := range u.anims {
			u.anims[i].Tick()
			u.anims[i].Render(u.x+offX, u.y+offY, u.z+offZ)
			//remove animation if it is done
			if u.anims[i].IsFinished() {
				u.anims = append(u.anims[:i], u.anims[i+1:]...)
			}
		}
	}
}

func (u UIElement) Dims() (int, int) {
	return u.width, u.height
}

func (u UIElement) Pos() (int, int, int) {
	return u.x, u.y, u.z
}

func (u *UIElement) SetTitle(txt string) {
	u.title = txt
}

func (u *UIElement) ToggleVisible() {
	u.visible = !u.visible
	console.Clear()
}

func (u *UIElement) SetVisibility(v bool) {
	if u.visible != v {
		u.ToggleVisible()
	}
}

func (u *UIElement) ToggleFocus() {
	u.focused = !u.focused
}

func (u *UIElement) MoveTo(x, y, z int) {
	u.x = x
	u.y = y
	u.z = z
}

func (u UIElement) Rect() (int, int, int, int) {
	return u.x, u.y, u.width, u.height
}

func (u UIElement) IsVisible() bool {
	return u.visible
}

func (u UIElement) IsFocused() bool {
	return u.focused
}

func (u *UIElement) AddAnimation(a Animator) {
	u.anims = append(u.anims, a)
}

func (u *UIElement) RemoveAnimation(a Animator) {
	for i, anim := range u.anims {
		if anim == a {
			u.anims = append(u.anims[:i], u.anims[i+1:]...)
		}
	}
}

//Centers the element within the console as a whole. Requires the console to be initialized first.
func (u *UIElement) CenterInConsole() {
	if console.Ready {
		w, h := console.Dims()
		u.x, u.y = (w-u.width)/2, (h-u.height)/2
	} else {
		LogError("UI Element cannot center: console not setup.")
	}
}

//Centers the element within the rect defined by (w, h, x, y)
func (u *UIElement) Center(w, h, x, y int) {
	u.CenterX(w, x)
	u.CenterY(h, y)
}

//Centers the element horizontally within the range defined by (w, x)
func (u *UIElement) CenterX(w, x int) {
	u.x = (w-x-u.width)/2 + x
}

//Centers the element vertically within the range defined by (h, y)
func (u *UIElement) CenterY(h, y int) {
	u.y = (h-y-u.height)/2 + y
}

//Sets a Tab number for the element. Elements in the same container can be
//tabbed back and forth across the TabIDs, starting with 1. Default is 0, which
//is ignored by the tabbing function.
func (u *UIElement) SetTabID(id int) {
	u.tabID = id
}

func (u UIElement) TabID() int {
	return u.tabID
}

func (u *UIElement) HandleKeypress(key sdl.Keycode) {
	//No-op. Maybe i'll make a default "no action associated with that key"
	//animation later, like maybe it subtly pulses once or something. Might be annoying though.
}

//Helper funtion for unpacking optional offsets passed to UI render functions. Required to allow for nesting of elements.
func processOffset(offset []int) (x, y, z int) {
	x, y, z = 0, 0, 0
	if len(offset) >= 2 {
		x, y = offset[0], offset[1]
		if len(offset) == 3 {
			z = offset[2]
		}
	}
	return
}
