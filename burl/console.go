package burl

//Console is the framebuffer which composites the UI of the current active State. This happens once
//per frame before being sent to the Renderer to be displayed.
type Console struct {
	Canvas
	//TODO sync.RWmutex?? i think??
}

//Setup the game window, renderer, etc
func (c *Console) Setup(w, h int) (err error) {
	c.Init(w, h)
	c.Clear()

	return nil
}
