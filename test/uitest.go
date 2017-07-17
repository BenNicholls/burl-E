package main

import "github.com/bennicholls/burl"
import "github.com/bennicholls/burl/console"
import "github.com/bennicholls/burl/ui"
import "math/rand"

import "time"
import "fmt"

func main() {

	rand.Seed(time.Now().UTC().UnixNano())

	err := console.Setup(80, 45, "../res/curses.bmp", "../res/DelveFont8x16.bmp", "UI Test")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer console.Cleanup()

	t := new(TestUI)
	t.SetupUI()
	console.ToggleFPS()
	burl.InitState(t)
	burl.GameLoop()

}

type TestUI struct {
	burl.BurlState
	container *ui.Container
	tiles     *ui.TileView
}

func (t *TestUI) SetupUI() {
	t.container = ui.NewContainer(40, 40, 1, 1, 0, true)
	textbox := ui.NewTextbox(30, 20, 2, 2, 0, true, false, "TESTING")
	textbox.ChangeText("Loremipsumdolorsitamet,consecteturadipiscingelit.Donecvitaenibhrisus. Quisque consectetur lacus eu velit viverra convallis. In at mattis orci. Suspendisse rhoncus lacinia elit ac ullamcorper. Donec id mattis velit, in condimentum massa. Nam non dui eu urna lacinia varius ut nec justo. Suspendisse consequat ornare neque, sit amet cursus enim volutpat in. Proin nibh ante, tempus in laoreet luctus, tempus in eros.")
	textbox.SetTitle("YAY")
	t.container.Add(textbox)
	t.container.SetTitle("FANCYTIMES")

	t.tiles = ui.NewTileView(40, 40, 41, 1, 0, true)
	t.tiles.DrawCircle(20, 20, 15, 0x32, 0xFFFFFFFF, 0)
}

func (t *TestUI) Render() {
	t.container.Render()
	t.tiles.Render()
}
