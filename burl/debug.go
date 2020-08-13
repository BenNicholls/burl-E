package burl

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

var debugger debugWindow
var debug bool
var debugWatches []watch
var debugCommands map[string]func ()

//Activate debugging capabilities. F10 will bring up the debug menu.
func Debug() {
	debug = true
	debugWatches = make([]watch, 0, 20)
	debugCommands = make(map[string]func ())

	if console != nil {
		initDebugWindow()
	}

	RegisterDebugCommand("fullscreen", renderer.ToggleFullscreen)
}

func initDebugWindow() {
	cw, ch := console.Dims()
	debugger = debugWindow{}
	debugger.PagedContainer = *NewPagedContainer(cw/2, ch/2, cw/4, ch/4, 50000, true)
	debugger.SetVisibility(false)

	debugger.logPage = debugger.AddPage("Logs")
	pw, ph := debugger.logPage.Dims()
	debugger.logList = NewList(pw, ph-2, 0, 0, 0, true, "No Log Entries")
	debugger.logList.ToggleHighlight()

	for i := range logger {
		debugger.logList.Append(logger[i].String())
	}
	debugger.logList.ScrollToBottom()

	debugger.logInput = NewInputbox(pw, 1, 0, ph-1, 0, false)
	debugger.logInput.cursorAnimation.Toggle()

	debugger.logPage.Add(debugger.logList, debugger.logInput)

	debugger.watchPage = debugger.AddPage("Watches")
	debugger.watchList = NewList(pw, ph, 0, 0, 0, false, "No Watched Variables")
	debugger.watchList.ToggleHighlight()

	for i := range debugWatches {
		debugger.watchList.Append(debugWatches[i].String())
	}

	debugger.watchPage.Add(debugger.watchList)

	debugger.flagsPage = debugger.AddPage("Flags")
	debugger.flagsList = NewList(pw, ph, 10, 0, 0, false, "No flags")

	debugger.fpsChoice = NewChoiceBox(4, 1, 0, 0, 0, false, HORIZONTAL, "OFF", "ON")
	debugger.changesChoice = NewChoiceBox(4, 1, 0, 0, 0, false, HORIZONTAL, "OFF", "ON")

	debugger.flagsPage.Add(NewTextbox(10, 1, 0, 0, 0, false, false, "FPS Counter"))
	debugger.flagsPage.Add(NewTextbox(10, 1, 0, 1, 0, false, false, "Show Renders"))

	debugger.flagsList.Add(debugger.fpsChoice, debugger.changesChoice)
	debugger.flagsPage.Add(debugger.flagsList)

	return
}

type debugWindow struct {
	PagedContainer

	logPage  *Container
	logList  *List
	logInput *Inputbox

	watchPage *Container
	watchList *List

	flagsPage     *Container
	flagsList     *List
	fpsChoice     *ChoiceBox
	changesChoice *ChoiceBox
}

func (dw *debugWindow) Update() {
	for i := range debugWatches {
		dw.watchList.Change(i, debugWatches[i].String())
	}
}

func (dw *debugWindow) HandleKeypress(key sdl.Keycode) {
	switch key {
	case sdl.K_F10:
		dw.ToggleVisible()
	case sdl.K_TAB:
		dw.NextPage()
	default:
		switch dw.CurrentIndex() {
		case 0: //logs
			if key == sdl.K_PAGEUP || key == sdl.K_PAGEDOWN {
				dw.logList.HandleKeypress(key)
			} else if key == sdl.K_RETURN {
				executeDebugCommand(dw.logInput.GetText())
				dw.logInput.Reset()
			} else {
				dw.logInput.HandleKeypress(key)
			}
		case 2: //flags
			dw.flagsList.HandleKeypress(key)
			if key == sdl.K_LEFT || key == sdl.K_RIGHT {
				switch dw.flagsList.GetSelection() {
				case 0: //fps
					renderer.ToggleDebugMode("fps")
				case 1: //render changes
					renderer.ToggleDebugMode("changes")	
				}
			}
		}
	}
}

//A watch is a variable that we're going to keep an eye on. The debug menu will display the current
//value of the variable as it changes.
type watch struct {
	label string
	value interface{}
}

func (w *watch) String() string {
	if w.value != nil {
		switch v := w.value.(type) {
		case *int:
			return fmt.Sprint(w.label, ": ", *v)
		case *float64:
			return fmt.Sprint(w.label, ": ", *v)
		case *string:
			return fmt.Sprint(w.label, ": ", *v)
		case fmt.Stringer:
			return fmt.Sprint(w.label, ": ", v.String())
		}
	}

	return fmt.Sprint(w.label, ": No value.")
}

//Register a watched variable. If debug mode is on, the value passed here will be available to display
//in the debug menu (F10). REMEMBER: the value MUST be a pointer!! Valid types that can be watched are
//int, float32/64, string, and anything with a String() method. ALSO REMEMBER: watches contain references,
//so anything with a watch won't be garbage collected.
func RegisterWatch(label string, val interface{}) {
	if debug {
		if val != nil {
			debugWatches = append(debugWatches, watch{label, val})
			debugger.watchList.Append(debugWatches[len(debugWatches)-1].String())
			LogInfo("Registered watch: ", label)
		} else {
			LogError("Bad Watch Register: ", label)
		}
	}
}

//Register a function to be called when you invoke the command in the debugger. The function must have no
//arguments and return no value.
func RegisterDebugCommand(command string, action func () ) {
	if debug {
		debugCommands[command] = action
		LogInfo("Added debug command: ", command)
	}
}

func executeDebugCommand(command string) {
	if debug {
		if action, ok := debugCommands[command]; ok {
			LogInfo("Execute debug command: ", command)
			action()
		} else {
			LogError("Bad command: ", command)
		}
	}
}