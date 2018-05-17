package burl

import (
	"github.com/veandco/go-sdl2/sdl"
)

//PagedContainer is a container for pages, with little tabs at the top. You can cycle through the pages
//like you'd expect. Cool, right?
type PagedContainer struct {
	UIElement
	curPage      int
	pages        []*Page
	redrawTitles bool
}

//Page is a single page in the PagedContainer. Might need to think of a better name for this.
type Page struct {
	page  *Container
	title *Textbox
}

func NewPagedContainer(w, h, x, y, z int, bord bool) *PagedContainer {
	p := new(PagedContainer)
	p.UIElement = NewUIElement(w, h, x, y, z, bord)
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
	newPage.page = NewContainer(p.width, p.height-3, p.x, p.y+3, p.z, true)
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
	p.curPage, _ = ModularClamp(p.curPage+1, 0, len(p.pages)-1)
	p.setActivePage()
}

func (p *PagedContainer) PrevPage() {
	p.curPage, _ = ModularClamp(p.curPage-1, 0, len(p.pages)-1)
	p.setActivePage()
}

//Finds the active page and fixes up visibilities, borders, etc.
func (p *PagedContainer) setActivePage() {
	for i := 0; i < len(p.pages); i++ {
		if i == p.curPage {
			p.pages[i].title.bordered = true
		} else {
			p.pages[i].title.bordered = false
		}
	}
	p.redrawTitles = true
	console.Clear()
}

func (p *PagedContainer) HandleKeypress(key sdl.Keycode) {
	switch key {
	case sdl.K_TAB:
		p.NextPage()
	}
}

func (p PagedContainer) Render() {
	if p.visible {
		p.UIElement.Render()

		if len(p.pages) > 0 {

			p.pages[p.curPage].page.Render()

			if p.redrawTitles {
				//draw over page title area
				console.Fill(p.x, p.y, p.z, p.width, 2, GLYPH_NONE, COL_BLACK, COL_BLACK)

				//draw titles
				for i, page := range p.pages {
					page.title.Render()
					if i == p.curPage {
						//remove border below title of selected page
						console.Clear(page.title.x-1, page.title.y+1, page.title.width+2, 1)
						if i == 0 {
							console.Clear(page.title.x-2, page.title.y+1, 1, 1)
							console.ChangeCell(page.title.x-2, page.title.y+1, p.z, GLYPH_BORDER_UDR, console.BorderColour(p.IsFocused()), COL_BLACK)
						}
						console.ChangeCell(page.title.x-1, page.title.y+1, p.z, GLYPH_BORDER_UL, console.BorderColour(p.IsFocused()), COL_BLACK)
						console.ChangeCell(page.title.x+page.title.width, page.title.y+1, p.z, GLYPH_BORDER_UR, console.BorderColour(p.IsFocused()), COL_BLACK)
					}
				}
				p.redrawTitles = false
			}
		}
	}
}
