package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/veandco/go-sdl2/sdl"

	"github.com/bennicholls/burl-E/burl"
)

func main() {

	rand.Seed(time.Now().UTC().UnixNano())

	_, err := burl.InitConsole(80, 45)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = burl.InitRenderer("res/curses.bmp", "res/DelveFont8x16.bmp", "UI Test")

	//burl.ToggleDebugMode("changes")

	t := new(TestUI)
	t.SetupUI()
	burl.InitState(t)
	burl.Debug()
	burl.GameLoop()
}

type TestUI struct {
	burl.StatePrototype
	tiles   *burl.TileView
	palette *burl.TileView
	paged   *burl.PagedContainer
}

func (t *TestUI) SetupUI() {
	t.InitWindow(false)

	textbox := burl.NewTextbox(30, 20, 10, 10, 0, true, false, "")
	textbox.ChangeString("Loremipsumdolorsitamet,consecteturadipiscingelit.Donecvitaenibhrisus. Quisque consectetur lacus eu velit viverra convallis. In at mattis orci. Suspendisse rhoncus lacinia elit ac ullamcorper. Donec id mattis velit, in condimentum massa. Nam non dui eu urna lacinia varius ut nec justo. Suspendisse consequat ornare neque, sit amet cursus enim volutpat in. Proin nibh ante, tempus in laoreet luctus, tempus in eros.")
	textbox.GetBorder().SetTitle("YAY")
	textbox.GetBorder().SetHint("hint goes here")
	t.Window.AddChild(textbox)

	t.tiles = burl.NewTileView(48, 15, 10, 1, 1, true)
	t.tiles.CenterInConsole()
	t.tiles.GetBorder().SetTitle("Whatever")
	t.tiles.LoadImageFromXP("res/anomaly.xp")
	t.Window.AddChild(t.tiles)
	p := burl.GeneratePalette(20, burl.COL_LIME, burl.COL_BLUE)
	p.Add(burl.GeneratePalette(20, burl.COL_BLUE, burl.COL_RED))

	t.palette = burl.NewTileView(40, 1, 4, 31, 0, true)
	t.palette.DrawPalette(0, 0, p, burl.HORIZONTAL)

	t.Window.AddChild(t.palette)

	t.paged = burl.NewPagedContainer(50, 30, 11, 11, 5, true)
	t.paged.AddPage("test")
	t.paged.AddPage("test2")
	t.paged.AddPage("test3")
	t.Window.AddChild(t.paged)
}

func (t *TestUI) HandleKeypress(key sdl.Keycode) {
	t.paged.HandleKeypress(key)
}
