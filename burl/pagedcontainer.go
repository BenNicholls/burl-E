package burl

import (
	"github.com/veandco/go-sdl2/sdl"
)

//PagedContainer is a container for pages, with little tabs at the top. You can cycle through the pages
//like you'd expect. Cool, right?
type PagedContainer struct {
	UIElement
	titleBorder UIElement
	curPage     int
	pages       []*Page
}

//Page is a single page in the PagedContainer. Might need to think of a better name for this.
type Page struct {
	page  *Container
	title *Textbox
}

func NewPagedContainer(w, h, x, y, z int, bord bool) *PagedContainer {
	p := new(PagedContainer)
	p.UIElement = NewUIElement(w, h, x, y, z, bord)
	p.titleBorder = NewUIElement(w, 2, x, y, z, true)
	p.curPage = 0
	p.pages = make([]*Page, 0, 0)

	return p
}

func (p *PagedContainer) Move(dx, dy, dz int) {
	p.UIElement.Move(dx, dy, dz)
	for i := range p.pages {
		p.pages[i].page.Move(dx, dy, dz)
		p.pages[i].title.Move(dx, dy, dz)
	}
}

//Adds a page to the PagedContainer, returning a pointer to the page itself
func (p *PagedContainer) AddPage(title string) *Container {
	offX := 1
	for _, e := range p.pages {
		offX += e.title.width + 1
	}
	paddedTitle := title
	if len(title)%2 == 1 {
		paddedTitle = paddedTitle + " "
	}
	titleBox := NewTextbox(len(paddedTitle)/2, 1, p.x+offX, p.y+1, p.z+2, false, false, paddedTitle)

	newPage := new(Page)
	newPage.title = titleBox
	newPage.page = NewContainer(p.width, p.height-3, p.x, p.y+3, p.z, false)
	p.pages = append(p.pages, newPage)
	p.setActivePage()

	return newPage.page
}

func (p PagedContainer) GetPageDims() (int, int) {
	return p.width - 2, p.height - 4
}

func (p PagedContainer) CurrentIndex() int {
	return p.curPage
}

func (p PagedContainer) CurrentPage() *Container {
	return p.pages[p.curPage].page
}

func (p *PagedContainer) NextPage() {
	p.Redraw()
	p.pages[p.curPage].page.ToggleVisible()
	p.curPage, _ = ModularClamp(p.curPage+1, 0, len(p.pages)-1)
	p.pages[p.curPage].page.ToggleVisible()
	p.setActivePage()
}

func (p *PagedContainer) PrevPage() {
	p.Redraw()
	p.pages[p.curPage].page.ToggleVisible()
	p.curPage, _ = ModularClamp(p.curPage-1, 0, len(p.pages)-1)
	p.pages[p.curPage].page.ToggleVisible()
	p.setActivePage()
}

//Finds the active page and fixes up visibilities, borders, etc.
func (p *PagedContainer) setActivePage() {
	if len(p.pages) > 0 {
		p.pages[p.curPage].page.SetVisibility(true)
		for i := 0; i < len(p.pages); i++ {
			if p.pages[i].title.bordered != (i == p.curPage) {
				p.pages[i].title.bordered = (i == p.curPage)
				p.dirty = true
			}
		}
	}
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

		if len(p.pages) > 0 {
			p.pages[p.curPage].page.Render()

			if p.dirty {
				//draw over page title area
				console.Clear(p.width, 2, p.x, p.y, p.z)
				p.titleBorder.Render()

				//draw titles
				for _, page := range p.pages {
					page.title.Render()
				}
				p.dirty = false
			}

			//remove border below title of selected page
			curTitle := p.pages[p.curPage].title

			console.Clear(curTitle.width+2, 1, curTitle.x-1, curTitle.y+1, p.z)
			if p.curPage == 0 {
				console.Clear(1, 1, curTitle.x-2, curTitle.y+1, p.z)
				console.ChangeCell(curTitle.x-2, curTitle.y+1, p.z+1, GLYPH_BORDER_UDR, console.BorderColour(p.IsFocused()), COL_BLACK)
			}
			console.ChangeCell(curTitle.x-1, curTitle.y+1, p.z, GLYPH_BORDER_UL, console.BorderColour(p.IsFocused()), COL_BLACK)
			console.ChangeCell(curTitle.x+curTitle.width, curTitle.y+1, p.z, GLYPH_BORDER_UR, console.BorderColour(p.IsFocused()), COL_BLACK)
		}
	}
}
