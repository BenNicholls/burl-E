package burl

import (
	"github.com/veandco/go-sdl2/sdl"
)

//UIElem is the basic definition for all UI elements.
type UIElem interface {
	Render()
	GetCanvas() *Canvas
	Redraw()
	TriggerRedraw()
	Dims() (w int, h int)
	Pos() (x int, y int, z int)
	Bounds() Rect
	Colours() (fore, back uint32)
	MoveTo(x, y, z int)
	Move(dx, dy, dz int)
	GetBorder() *Border
	ToggleVisible()
	SetVisibility(v bool)
	IsVisible() bool
	SetForeColour(colour uint32)
	SetBackColour(colour uint32)
	ToggleFocus()
	IsFocused() bool
	CenterInConsole()
	SetTabID(id int)
	TabID() int
	HandleKeypress(key sdl.Keycode)
	AddChild(...UIElem)
	ClearChildren()
	GetParent() UIElem
	SetParent(UIElem)
}

type UIElement struct {
	Canvas

	parent     UIElem   //the UI Element this one is nested inside. Defaults to nil, updated when added to a container type.
	children   []UIElem //elements nested within this one. children are rendered after the UI Element proper.
	x, y, z    int      //position relative to its parent. UIElements cannot be drawn outside of their parent
	visible    bool
	focused    bool
	tabID      int    //for to tab between elements
	dirty      bool   //if we need to re-render the contents of the element
	redraw     bool   //if this element needs to be redrawn by a parent object
	foreColour uint32 //defaults to COL_WHITE
	backColour uint32 //defaults to COL_BLACK
	border     Border

	anims []Animator
}

func NewUIElement(w, h, x, y, z int, bord bool) UIElement {
	element := UIElement{
		children:   make([]UIElem, 0, 0),
		x:          x,
		y:          y,
		z:          z,
		visible:    true,
		anims:      make([]Animator, 0, 0),
		dirty:      true,
		redraw:     false,
		foreColour: COL_WHITE,
		backColour: COL_BLACK,
	}
	element.Canvas.Init(w, h)
	element.Redraw()
	element.border = console.defaultBorderStyle
	element.border.Set(bord)
	return element
}

//basic render function for all elements.
func (u *UIElement) Render() {
	if !u.visible {
		return
	}

	//update child elements (if any). renders on children can propogate signals (redraw, etc) to this element, so this
	//must be done first
	for _, child := range u.children {
		child.Render()
	}

	if u.redraw {
		u.Redraw()
	}

	for i, _ := range u.anims {
		u.anims[i].Tick()
		//u.anims[i].Render(u.x, u.y, u.z)
		//remove animation if it is done
		if u.anims[i].IsFinished() {
			u.anims = append(u.anims[:i], u.anims[i+1:]...)
		}
	}

	//composite together children elements, if any exist
	for _, child := range u.children { //TODO: do these need to be ordered by Z depth? hmm.
		if child.IsVisible() {
			x, y, z := child.Pos()
			w, h := child.Dims()
			u.CopyFromCanvas(x, y, z, child.GetCanvas())
			u.DrawBorder(x, y, z, w, h, child.GetBorder())
		}
	}

	u.dirty = false
}

func (u *UIElement) GetCanvas() *Canvas {
	return &u.Canvas
}

func (u *UIElement) Redraw() {
	u.Clear()
	u.Fill(0, 0, 0, u.width, u.height, GLYPH_NONE, u.foreColour, u.backColour)
	u.dirty = true
	u.redraw = false
}

func (u *UIElement) TriggerRedraw() {
	u.redraw = true
}

func (u *UIElement) Dims() (int, int) {
	return u.width, u.height
}

func (u *UIElement) Pos() (int, int, int) {
	return u.x, u.y, u.z
}

func (u *UIElement) Bounds() Rect {
	return Rect{u.width, u.height, u.x, u.y}
}

func (u *UIElement) Colours() (uint32, uint32) {
	return u.foreColour, u.backColour
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

func (u *UIElement) SetForeColour(c uint32) {
	u.foreColour = c
}

func (u *UIElement) SetBackColour(c uint32) {
	u.backColour = c
}

func (u *UIElement) ToggleVisible() {
	u.visible = !u.visible
	if u.visible {
		u.redraw = true
		u.border.redraw = true
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
	if u.focused {
		u.border.SetColour(COL_PURPLE)
	} else {
		u.border.SetColour(COL_LIGHTGREY)
	}
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

//Add any number of children to UIElement.
//TODO: loop check? guard against double-adds? z-depth sorting? all good ideas. rethink once UI element IDs or whatever go in.
func (u *UIElement) AddChild(elems ...UIElem) {
	for _, elem := range elems {
		u.children = append(u.children, elem)
		elem.SetParent(u)
	}
}

func (u *UIElement) ClearChildren() {
	//de-link any children, then remake the array
	for _, child := range u.children {
		child.SetParent(nil)
	}

	u.children = make([]UIElem, 0, 0)
	u.Redraw()
}

func (u *UIElement) GetParent() UIElem {
	return u.parent
}

func (u *UIElement) SetParent(elem UIElem) {
	u.parent = elem
}

//Centers the element within the console as a whole. Requires the console to be initialized first.
func (u *UIElement) CenterInConsole() {
	if console != nil {
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

func (u *UIElement) GetBorder() *Border {
	return &u.border
}

//Border data for UI Elements. Will eventually also contain data for styling
type Border struct {
	enabled    bool
	title      string
	hint       string
	redraw     bool //if the border needs to be redrawn by a parent element
	foreColour uint32
	backColour uint32
}

//Enables or disables the Border
func (b *Border) Set(bord bool) {
	if b.enabled != bord {
		b.Toggle()
	}
}

//Toggles the border on/off
func (b *Border) Toggle() {
	b.enabled = !b.enabled
	if b.enabled {
		b.redraw = true
	}
}

func (b *Border) SetTitle(title string) {
	if b.title != title {
		b.title = title
		b.redraw = true
	}
}

func (b *Border) SetHint(hint string) {
	if b.hint != hint {
		b.hint = hint
		b.redraw = true
	}
}

//Sets the colour of the border
//TODO: capability to set the backcolour? reading from border styling data somewhere? iunno
func (b *Border) SetColour(col uint32) {
	if b.foreColour != col {
		b.foreColour = col
		b.redraw = true
	}
}
