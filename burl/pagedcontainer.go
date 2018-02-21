package burl

type PagedContainer struct {
	Container
	page         int
	pageTitles   []*Textbox
	redrawTitles bool
}

func NewPagedContainer(w, h, x, y, z int, bord bool) *PagedContainer {
	p := new(PagedContainer)
	p.Container = *NewContainer(w, h, x, y, z, bord)
	p.page = 3
	p.pageTitles = make([]*Textbox, 0, 0)
	p.redraw = true

	return p
}

func (p *PagedContainer) AddPage(title string, page UIElem) {
	p.Add(page)
	offX := 1
	for _, t := range p.pageTitles {
		offX += t.width + 1
	}
	titleBox := NewTextbox(len(title)/2+2, 1, offX, 1, 1, false, true, title)
	p.pageTitles = append(p.pageTitles, titleBox)
	p.setActivePage()
}

func (p PagedContainer) GetPageDims() (int, int) {
	return p.width, p.height - 2
}

func (p *PagedContainer) NextPage() {
	p.page, _ = ModularClamp(p.page+1, 0, len(p.pageTitles)-1)
	p.setActivePage()
}

func (p *PagedContainer) PrevPage() {
	p.page, _ = ModularClamp(p.page-1, 0, len(p.pageTitles)-1)
	p.setActivePage()
}

//Finds the active page and fixes up visibilities, borders, etc.
func (p *PagedContainer) setActivePage() {
	for i := 0; i < len(p.pageTitles); i++ {
		if i == p.page {
			p.pageTitles[i].bordered = true
		} else {
			p.pageTitles[i].bordered = false
		}
	}
	p.redrawTitles = true
	console.Clear() //TOTAL HACK. Note: should probably not be using Container as the embedded type here. Maybe make a Page struct??
}

func (p PagedContainer) Render(offset ...int) {
	if p.visible {
		offX, offY, offZ := processOffset(offset)

		p.Container.UIElement.Render(offX, offY, offZ)
		p.Container.Elements[p.page].Render(p.x+offX, p.y+offY, p.z+offZ)

		if p.redrawTitles {

			//draw over page title area
			for i := 0; i < p.width*2; i++ {
				console.ChangeCell(p.x+offX+(i%p.width), p.y+offY+(i/p.width), p.z+offZ, GLYPH_NONE, COL_BLACK, COL_BLACK)
			}

			//draw titles
			for i, titleBox := range p.pageTitles {
				titleBox.Render(p.x+offX, p.y+offY, p.z+offZ)
				if i == p.page {
					//remove border below title of selected page
					console.Clear(p.x+offX+titleBox.x, p.y+offY+titleBox.y+1, titleBox.width, 1)
				}
			}
			p.redrawTitles = false
		}
	}
}
