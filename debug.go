package burl

import "github.com/bennicholls/burl/console"
import "github.com/bennicholls/burl/util"

func DebugToggleRenderFPS() {
	if console.Ready {
		console.ToggleFPS()
	} else {
		util.LogError("Burl Debugging: Could not toggle FPS, console not set up!")
	}
}

func DebugToggleRenderChangeView() {
	if console.Ready {
		console.ToggleChanges()
		console.Clear()
	} else {
		util.LogError("Burl Debugging: Could not toggle render change viewing, console not set up!")
	}
}
