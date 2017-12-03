package main

import "github.com/bennicholls/burl-E/burl"
import "math/rand"
import "time"
import "fmt"

func main() {

	rand.Seed(time.Now().UTC().UnixNano())

	console, err := burl.InitConsole(80, 45, "../res/curses.bmp", "../res/DelveFont8x16.bmp", "UI Test")
	if err != nil {
		fmt.Println(err)
		return
	}

	console.ToggleChanges()

	t := new(TestUI)
	t.SetupUI()
	console.ToggleFPS()
	burl.InitState(t)
	burl.GameLoop()
}

type TestUI struct {
	burl.BaseState
	container *burl.Container
	tiles     *burl.TileView
	
	yes bool
}

func (t *TestUI) SetupUI() {
	t.container = burl.NewContainer(40, 20, 2, 2, 0, true)
	t.container.SetTitle("")
	
	textbox := burl.NewTextbox(30, 20, 2, 2, 0, true, false, "")
	textbox.ChangeText("Loremipsumdolorsitamet,consecteturadipiscingelit.Donecvitaenibhrisus. Quisque consectetur lacus eu velit viverra convallis. In at mattis orci. Suspendisse rhoncus lacinia elit ac ullamcorper. Donec id mattis velit, in condimentum massa. Nam non dui eu urna lacinia varius ut nec justo. Suspendisse consequat ornare neque, sit amet cursus enim volutpat in. Proin nibh ante, tempus in laoreet luctus, tempus in eros.")
	textbox.SetTitle("YAY")
	t.container.Add(textbox)
	
	t.tiles = burl.NewTileView(30, 20, 45, 1, 0, true)
	t.tiles.SetTitle("Whatever")
	t.tiles.DrawCircle(15, 10, 2, burl.GLYPH_FACE1, 0xFFFFFFFF, 0xFF000000)
	t.yes = false
}

func (t *TestUI) Render() {
	if !t.yes {
		t.container.Render()
		t.tiles.Render()
		t.yes = true
	}
}
