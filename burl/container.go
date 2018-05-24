package burl

//UI Element that acts as a way to group other elements. Allows for nesting of elements, etc.
type Container struct {
	UIElement
	redraw bool

	Elements []UIElem
}

func NewContainer(w, h, x, y, z int, bord bool) *Container {
	return &Container{NewUIElement(w, h, x, y, z, bord), true, make([]UIElem, 0, 20)}
}

//Adds any number of UIElem to the container.
func (c *Container) Add(elems ...UIElem) {
	for _, e := range elems {
		e.Move(c.x, c.y, c.z)
		c.Elements = append(c.Elements, e)
	}
}

func (c *Container) Move(dx, dy, dz int) {
	c.UIElement.Move(dx, dy, dz)
	for i := range c.Elements {
		c.Elements[i].Move(dx, dy, dz)
	}
}

func (c *Container) MoveTo(x, y, z int) {
	for i := range c.Elements {
		c.Elements[i].Move(x-c.x, y-c.y, z-c.z)
	}

	c.x = x
	c.y = y
	c.z = z
}

//Deletes all UIElem from the container.
func (c *Container) ClearElements() {
	if len(c.Elements) > 0 {
		c.Elements = make([]UIElem, 0, 20)
		c.redraw = true
	}
}

//Finds the next element in the tabbing order (if one is defined) among elements
//in the container. Cycles back to the top of the order once the bottom is reached.
//Returns the original element if no next element is found.
func (c *Container) FindNextTab(e UIElem) UIElem {
	var next UIElem
	var top UIElem

	for i := range c.Elements {
		if c.Elements[i].TabID() > e.TabID() {
			if next == nil || c.Elements[i].TabID() < next.TabID() {
				next = c.Elements[i]
			}
		}

		if c.Elements[i].TabID() > 0 {
			if top == nil || c.Elements[i].TabID() < top.TabID() {
				top = c.Elements[i]
			}
		}
	}

	if next == nil {
		if top == nil {
			return e
		} else {
			return top
		}
	}

	return next
}

//Finds the previous element in the tabbing order (if one is defined) among elements
//in the container. Cycles back to the bottom of the order once the top is reached.
//Returns the original element if no previous element is found.
func (c *Container) FindPrevTab(e UIElem) UIElem {
	var prev UIElem
	var bottom UIElem

	for i := range c.Elements {
		if c.Elements[i].TabID() > 0 && c.Elements[i].TabID() < e.TabID() {
			if prev == nil || c.Elements[i].TabID() > prev.TabID() {
				prev = c.Elements[i]
			}
		}

		if c.Elements[i].TabID() > 0 {
			if bottom == nil || c.Elements[i].TabID() > bottom.TabID() {
				bottom = c.Elements[i]
			}
		}
	}

	if prev == nil {
		if bottom == nil {
			return e
		} else {
			return bottom
		}
	}

	return prev
}

func (c *Container) Render() {
	if c.visible {
		if c.redraw {
			c.Redraw()
			c.redraw = false
		}

		for i := 0; i < len(c.Elements); i++ {
			c.Elements[i].Render()
		}

		c.UIElement.Render()
	}
}
