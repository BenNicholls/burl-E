package ui

import "github.com/bennicholls/burl/console"

//UI Element that acts as a way to group other elements. Allows for nesting of elements, etc.
type Container struct {
	width, height int
	x, y, z       int
	bordered      bool
	title         string
	visible       bool
	focused       bool
	redraw        bool

	Elements []UIElem
}

func NewContainer(w, h, x, y, z int, bord bool) *Container {
	return &Container{w, h, x, y, z, bord, "", true, false, true, make([]UIElem, 0, 20)}
}

//Adds any number of UIElem to the container.
func (c *Container) Add(elems ...UIElem) {
	for _, e := range elems {
		c.Elements = append(c.Elements, e)
	}
}

//Deletes all UIElem from the container.
func (c *Container) ClearElements() {
	c.Elements = make([]UIElem, 0, 20)
	c.redraw = true
}

func (c *Container) SetTitle(s string) {
	c.title = s
}

//Offets (x,y,z, all optional) are passed through to the nested elements.
func (c *Container) Render(offset ...int) {
	if c.visible {
		offX, offY, offZ := processOffset(offset)

		//previously was just clearing on special occasions, but it has come to my attention that
		//containers floating above other content need to blank the screen or else they appear transparent.
		//in the future, maybe we can have a transparency setting? would that be useful? who knows.
		// if c.redraw {
		// 	console.Clear(c.x+offX, c.y+offY, c.width, c.height)
		// 	c.redraw = false
		// }

		console.Clear(c.x+offX, c.y+offY, c.width, c.height)

		for i := 0; i < len(c.Elements); i++ {
			c.Elements[i].Render(c.x+offX, c.y+offY, c.z+offZ)
		}

		if c.bordered {
			console.DrawBorder(c.x+offX, c.y+offY, c.z+offZ, c.width, c.height, c.title, c.focused)
		}
	}
}

func (c Container) Dims() (int, int) {
	return c.width, c.height
}

func (c Container) Pos() (int, int, int) {
	return c.x, c.y, c.z
}

func (c *Container) ToggleVisible() {
	c.visible = !c.visible
	console.Clear()
}

func (c *Container) SetVisibility(v bool) {
	c.visible = v
	console.Clear()
}

func (c *Container) ToggleFocus() {
	c.focused = !c.focused
}

func (c *Container) MoveTo(x, y, z int) {
	c.x = x
	c.y = y
	c.z = z
}

func (c Container) IsVisible() bool {
	return c.visible
}
