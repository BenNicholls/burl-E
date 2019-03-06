package burl

import (
	"github.com/veandco/go-sdl2/sdl"
)

//UIElem is the basic definition for all UI elements.
type UIElem interface {
	Render()
	Redraw()
	Dims() (w int, h int)
	Pos() (x int, y int, z int)
	Bounds() Rect
	MoveTo(x, y, z int)
	Move(dx, dy, dz int)
	SetTitle(title string)
	SetHint(hint string)
	ToggleVisible()
	SetVisibility(v bool)
	IsVisible() bool
	SetBackColour(colour uint32)
	ToggleFocus()
	IsFocused() bool
	CenterInConsole()
	SetTabID(id int)
	TabID() int
	HandleKeypress(key sdl.Keycode)
}

type UIElement struct {
	x, y, z       int
	width, height int
	bordered      bool
	title         string
	hint          string
	visible       bool
	focused       bool
	tabID         int    //for to tab between elements in a container
	dirty         bool   //only used for some elements. could be used all around probably??
	backColour    uint32 //defaults to COL_BLACK. forecolour is controlled by the specific element type.

	anims []Animator
}

func NewUIElement(w, h, x, y, z int, bord bool) UIElement {
	return UIElement{
		x:        x,
		y:        y,
		z:        z,
		width:    w,
		height:   h,
		bordered: bord,
		visible:  true,
		anims:    make([]Animator, 0, 20),
		dirty:    true,
		backColour: COL_BLACK,
	}
}

//basic render function for all elements.
func (u *UIElement) Render() {
	if u.visible {
		if u.bordered {
			console.DrawBorder(u.x, u.y, u.z, u.width, u.height, u.title, u.hint, u.focused)
		}

		for i, _ := range u.anims {
			u.anims[i].Tick()
			u.anims[i].Render(u.x, u.y, u.z)
			//remove animation if it is done
			if u.anims[i].IsFinished() {
				u.anims = append(u.anims[:i], u.anims[i+1:]...)
			}
		}
	}
}

func (u *UIElement) Redraw() {
	if u.visible {
		if u.bordered {
			console.Fill(u.x-1, u.y-1, u.z, u.width+2, u.height+2, GLYPH_NONE, COL_BLACK, u.backColour)
		} else {
			console.Fill(u.x, u.y, u.z, u.width, u.height, GLYPH_NONE, COL_BLACK, u.backColour)
		}
		u.dirty = true
	}
}

func (u *UIElement) Dims() (int, int) {
	return u.width, u.height
}

func (u*UIElement) Pos() (int, int, int) {
	return u.x, u.y, u.z
}

func (u *UIElement) Bounds() Rect {
	return Rect{u.width, u.height, u.x, u.y}
}

func (u *UIElement) Move(dx, dy, dz int) {
	u.x += dx
	u.y += dy
	u.z += dz
}

func (u *UIElement) MoveTo(x, y, z int) {
	u.x = x
	u.y = y
	u.z = z
}

func (u *UIElement) SetTitle(txt string) {
	u.title = txt
}

func (u *UIElement) SetHint(txt string) {
	u.hint = txt
}

func (u *UIElement) SetBackColour(c uint32) {
	u.backColour = c
}

func (u *UIElement) ToggleVisible() {
	u.visible = !u.visible

	if !u.visible {
		if u.bordered {
			console.Clear(u.width+2, u.height+2, u.x-1, u.y-1)
		} else {
			console.Clear(u.width, u.height, u.x, u.y)
		}

		//redraw elements underneath this one
		if gameState != nil && gameState.GetWindow() != nil {
			if gameState.GetDialog() != nil {
				gameState.GetDialog().GetWindow().Redraw()
			}			
			for _, elem := range gameState.GetWindow().Elements {
				eX, eY, _ := elem.Pos()
				eW, eH := elem.Dims()

				eRect := Rect{eW + 2, eH + 2, eX - 1, eY - 1} //bounding rect is bigger to account for borders

				intersection := FindIntersectionRect(u.Bounds(), eRect)

				if intersection.W != 0 && intersection.H != 0 {
					elem.Redraw()
				}
			}
		}
	} else {
		u.Redraw()
	}
}

func (u *UIElement) SetVisibility(v bool) {
	if u.visible != v {
		u.ToggleVisible()
	}
}

func (u *UIElement) IsVisible() bool {
	return u.visible
}

func (u *UIElement) ToggleFocus() {
	u.focused = !u.focused
}

func (u *UIElement) IsFocused() bool {
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
		u.MoveTo((w-u.width)/2, (h-u.height)/2, u.z)
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
	u.MoveTo((w-x-u.width)/2+x, u.y, u.z)
}

//Centers the element vertically within the range defined by (h, y)
func (u *UIElement) CenterY(h, y int) {
	u.MoveTo(u.x, (h-y-u.height)/2+y, u.z)
}

//Sets a Tab number for the element. Elements in the same container can be tabbed back and forth across
//the TabIDs, starting with 1. Default is 0, which is ignored by the tabbing function.
func (u *UIElement) SetTabID(id int) {
	u.tabID = id
}

func (u *UIElement) TabID() int {
	return u.tabID
}

func (u *UIElement) HandleKeypress(key sdl.Keycode) {
	//No-op. Maybe i'll make a default "no action associated with that key"
	//animation later, like maybe it subtly pulses once or something. Might be annoying though.
}
