package ui

import "github.com/bennicholls/burl/console"

//UI Element that acts as a way to group other elements. Allows for nesting of elements, etc.
type Container struct {
	UIElement
	redraw bool

	Elements []UIElem
}

func NewContainer(w, h, x, y, z int, bord bool) *Container {
	return &Container{NewUIElement(x, y, z, w, h, bord), true, make([]UIElem, 0, 20)}
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

//Offets (x,y,z, all optional) are passed through to the nested elements.
func (c *Container) Render(offset ...int) {
	if c.visible {
		offX, offY, offZ := processOffset(offset)

		if c.redraw {
			console.Clear(c.x+offX, c.y+offY, c.width, c.height)
			c.redraw = false
		}

		//draw over container, so we don't appear transparent.
		for i := 0; i < c.width*c.height; i++ {
			console.ChangeColours(c.x+offX+(i%c.width), c.y+offY+(i/c.width), c.z+offZ, 0xFF000000, 0xFF000000)
		}

		for i := 0; i < len(c.Elements); i++ {
			c.Elements[i].Render(c.x+offX, c.y+offY, c.z+offZ)
		}

		c.UIElement.Render(offX, offY, offZ)
	}
}
