package burl

import (
	"github.com/veandco/go-sdl2/sdl"
)

//PagedContainer is a container for pages, with little tabs at the top. You can cycle through the pages
//like you'd expect. Cool, right?
type PagedContainer struct {
	UIElement

	pages   UIElement
	curPage int
	titles  UIElement
}

func NewPagedContainer(w, h, x, y, z int, bord bool) *PagedContainer {
	p := new(PagedContainer)
	p.UIElement = NewUIElement(w, h, x, y, z, bord)
	p.curPage = 0
	p.titles = NewUIElement(w, 2, 0, 0, 0, true)
	p.pages = NewUIElement(w, h-2, 0, 3, 0, false)

	p.AddChild(&p.titles, &p.pages)

	return p
}

//Adds a page to the PagedContainer, returning a pointer to the page itself
func (p *PagedContainer) AddPage(title string) *UIElement {
	newPage := NewUIElement(p.width, p.height-2, 0, 0, 0, false)
	p.pages.AddChild(&newPage)

	offX := 1
	for _, e := range p.titles.children {
		w, _ := e.Dims()
		offX += w + 1
	}
	paddedTitle := title
	if len(title)%2 == 1 {
		paddedTitle = paddedTitle + " "
	}

	newTitle := NewTextbox(len(paddedTitle)/2, 1, offX, 1, 1, false, false, paddedTitle)
	p.titles.AddChild(newTitle)

	p.setActivePage()

	return &newPage
}

func (p *PagedContainer) GetPageDims() (int, int) {
	return p.pages.Dims()
}

func (p *PagedContainer) CurrentIndex() int {
	return p.curPage
}

func (p *PagedContainer) CurrentPage() *UIElement {
	return p.pages.children[p.curPage].(*UIElement)
}

func (p *PagedContainer) NextPage() {
	if len(p.pages.children) <= 1 {
		return
	}

	p.curPage, _ = ModularClamp(p.curPage+1, 0, len(p.pages.children)-1)
	p.setActivePage()
}

func (p *PagedContainer) PrevPage() {
	if len(p.pages.children) <= 1 {
		return
	}

	p.curPage, _ = ModularClamp(p.curPage-1, 0, len(p.pages.children)-1)
	p.setActivePage()
}

//Finds the active page and fixes up visibilities, borders, etc.
func (p *PagedContainer) setActivePage() {
	if len(p.pages.children) == 0 {
		return
	}

	for i := range p.pages.children {
		p.pages.children[i].SetVisibility(i == p.curPage)
		p.titles.children[i].GetBorder().Set(i == p.curPage)
	}

	p.redraw = true
	p.titles.redraw = true
	p.pages.redraw = true
}

func (p *PagedContainer) HandleKeypress(key sdl.Keycode) {
	switch key {
	case sdl.K_TAB:
		p.NextPage()
	}
}

func (p *PagedContainer) Render() {
	if p.visible {
		p.UIElement.Render()

		// if len(p.pages) > 0 {
		// 	p.pages[p.curPage].Render()

		// 	if p.dirty {
		// 		//draw over page title area
		// 		p.Clear(p.width, 2, 0, 0, 0)

		// 		//draw titles
		// 		for _, page := range p.pages {
		// 			page.title.Render()
		// 		}
		// 		p.dirty = false
		// 	}

		// 	//remove border below title of selected page
		// 	curTitle := p.pages[p.curPage].title
		// 	p.Fill(curTitle.x-1, curTitle.y+1, 1, curTitle.width+2, 1, GLYPH_NONE, COL_BLACK, COL_BLACK)
		// 	if p.curPage == 0 {
		// 		p.ChangeCell(curTitle.x-2, curTitle.y+1, 1, GLYPH_BORDER_UDR, p.BorderColour(p.IsFocused()), COL_BLACK)
		// 	}
		// 	p.ChangeCell(curTitle.x-1, curTitle.y+1, 1, GLYPH_BORDER_UL, p.BorderColour(p.IsFocused()), COL_BLACK)
		// 	p.ChangeCell(curTitle.x+curTitle.width, curTitle.y+1, 1, GLYPH_BORDER_UR, p.BorderColour(p.IsFocused()), COL_BLACK)
		// }
	}
}
