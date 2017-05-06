package main   

import "github.com/bennicholls/burl/console"
import "github.com/bennicholls/burl/ui"
import "runtime"
import "math/rand"
import "github.com/veandco/go-sdl2/sdl"
import "time"
import "fmt"

var container *ui.Container
var text *ui.Textbox

func main() {

	runtime.LockOSThread()
	rand.Seed(time.Now().UTC().UnixNano())

	var event sdl.Event

	err := console.Setup(80, 45, "res/curses.bmp", "res/DelveFont8x16.bmp", "UI Test")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer console.Cleanup()

	SetupUI()

	running := true

	for running {
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
			// case *sdl.MouseMotionEvent:
			// 	fmt.Printf("[%d ms] MouseMotion\ttype:%d\tid:%d\tx:%d\ty:%d\txrel:%d\tyrel:%d\n",
			// 		t.Timestamp, t.Type, t.Which, t.X, t.Y, t.XRel, t.YRel)
			// case *sdl.MouseButtonEvent:
			// 	fmt.Printf("[%d ms] MouseButton\ttype:%d\tid:%d\tx:%d\ty:%d\tbutton:%d\tstate:%d\n",
			// 		t.Timestamp, t.Type, t.Which, t.X, t.Y, t.Button, t.State)
			// case *sdl.MouseWheelEvent:
			// 	fmt.Printf("[%d ms] MouseWheel\ttype:%d\tid:%d\tx:%d\ty:%d\n",
			// 		t.Timestamp, t.Type, t.Which, t.X, t.Y)
			case *sdl.KeyUpEvent:
				//fmt.Printf("[%d ms] Keyboard\ttype:%d\tsym:%c\tmodifiers:%d\tstate:%d\trepeat:%d\n",
				//	t.Timestamp, t.Type, t.Keysym.Sym, t.Keysym.Mod, t.State, t.Repeat)
				running = false
			}
		}

		Render()
		console.Render()
	}
}

func SetupUI() {
	container = ui.NewContainer(40, 40, 1, 1, 0, true)
	textbox := ui.NewTextbox(30,20,2,2,0,true,false, "TESTING")
	textbox.ChangeText("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec vitae nibh risus. Quisque consectetur lacus eu velit viverra convallis. In at mattis orci. Suspendisse rhoncus lacinia elit ac ullamcorper. Donec id mattis velit, in condimentum massa. Nam non dui eu urna lacinia varius ut nec justo. Suspendisse consequat ornare neque, sit amet cursus enim volutpat in. Proin nibh ante, tempus in laoreet luctus, tempus in eros.")
	textbox.SetTitle("YAY")
	container.Add(textbox)
	container.SetTitle("FANCYTIMES")

}

func Render() {
	container.Render()
}