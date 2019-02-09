package main

import "github.com/bennicholls/burl-E/burl"
import "math/rand"
import "time"
import "fmt"

func main() {

	rand.Seed(time.Now().UTC().UnixNano())

	console, err := burl.InitConsole(80, 45, "res/curses.bmp", "res/DelveFont8x16.bmp", "UI Test")
	if err != nil {
		fmt.Println(err)
		return
	}

	//console.ToggleChanges()

	t := new(TestUI)
	t.SetupUI()
	console.ToggleFPS()
	burl.InitState(t)
	burl.GameLoop()
}

type TestUI struct {
	burl.StatePrototype
	tiles   *burl.TileView
	palette *burl.TileView

	yes bool
}

func (t *TestUI) SetupUI() {
	t.InitWindow(false)

	textbox := burl.NewTextbox(30, 20, 2, 2, 0, true, false, "")
	textbox.ChangeText("Loremipsumdolorsitamet,consecteturadipiscingelit.Donecvitaenibhrisus. Quisque consectetur lacus eu velit viverra convallis. In at mattis orci. Suspendisse rhoncus lacinia elit ac ullamcorper. Donec id mattis velit, in condimentum massa. Nam non dui eu urna lacinia varius ut nec justo. Suspendisse consequat ornare neque, sit amet cursus enim volutpat in. Proin nibh ante, tempus in laoreet luctus, tempus in eros.")
	textbox.SetTitle("YAY")
	t.Window.Add(textbox)

	t.tiles = burl.NewTileView(48, 15, 10, 1, 1, true)
	t.tiles.CenterInConsole()
	t.tiles.SetTitle("Whatever")
	t.tiles.LoadImageFromXP("res/anomaly.xp")

	p := burl.GeneratePalette(20, burl.COL_LIME, burl.COL_BLUE)
	p.Add(burl.GeneratePalette(20, burl.COL_BLUE, burl.COL_RED))
	burl.LogInfo(len(p))

	t.palette = burl.NewTileView(40, 1, 4, 36, 3, true)
	t.palette.DrawPalette(0, 0, p)

	t.Window.Add(t.palette)

	t.yes = false
}

func (t *TestUI) Render() {
	if !t.yes {
		t.tiles.Render()
		t.yes = true
	}
}
