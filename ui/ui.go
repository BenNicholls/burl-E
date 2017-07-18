package ui

import "github.com/bennicholls/burl/console"

//UIElem is the basic definition for all UI elements.
type UIElem interface {
	Render(offset ...int)
	Dims() (w int, h int)
	Pos() (x int, y int, z int)
	SetTitle(title string)
	ToggleVisible()
	SetVisibility(v bool)
	IsVisible() bool
	IsFocused() bool
	MoveTo(x, y, z int)
	Rect() (int, int, int, int)
}

type UIElement struct {
	x, y, z       int
	width, height int
	bordered      bool
	title         string
	visible       bool
	focused       bool

	anims []Animator
}

func NewUIElement(x, y, z, width, height int, bord bool) UIElement {
	return UIElement{x, y, z, width, height, bord, "", true, false, make([]Animator, 0, 20)}
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
