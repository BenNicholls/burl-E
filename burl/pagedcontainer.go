package burl

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

//Adds a page to the PagedContainer, returning a pointer to the page itself
func (p *PagedContainer) AddPage(title string) *Container {
	offX := 1
	for _, e := range p.pages {
		offX += e.title.width + 1
	}
	titleBox := NewTextbox(len(title)/2+2, 1, offX, 1, 1, false, true, title)

	newPage := new(Page)
	newPage.title = titleBox
	newPage.page = NewContainer(p.width-2, p.height-4, 1, 3, 1, true)
	p.pages = append(p.pages, newPage)
	p.setActivePage()

	return newPage.page
}

func (p PagedContainer) GetPageDims() (int, int) {
	return p.width - 2, p.height - 4
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

func (p PagedContainer) Render(offset ...int) {
	if p.visible {
		offX, offY, offZ := processOffset(offset)

		p.UIElement.Render(offX, offY, offZ)

		if len(p.pages) > 0 {

			p.pages[p.curPage].page.Render(p.x+offX, p.y+offY, p.z+offZ)

			if p.redrawTitles {
				//draw over page title area
				for i := 0; i < p.width*2; i++ {
					console.ChangeCell(p.x+offX+(i%p.width), p.y+offY+(i/p.width), p.z+offZ, GLYPH_NONE, COL_BLACK, COL_BLACK)
				}

				//draw titles
				for i, page := range p.pages {
					page.title.Render(p.x+offX, p.y+offY, p.z+offZ)
					if i == p.curPage {
						//remove border below title of selected page
						console.Clear(p.x+offX+page.title.x, p.y+offY+page.title.y+1, page.title.width, 1)
					}
				}
				p.redrawTitles = false
			}
		}
	}
}
