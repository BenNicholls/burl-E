package burl

func DebugToggleRenderFPS() {
	if console.Ready {
		console.ToggleFPS()
	} else {
		LogError("Burl Debugging: Could not toggle FPS, console not set up!")
	}
}

func DebugToggleRenderChangeView() {
	if console.Ready {
		console.ToggleChanges()
		console.Clear()
	} else {
		LogError("Burl Debugging: Could not toggle render change viewing, console not set up!")
	}
}
