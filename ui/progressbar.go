package ui

import "github.com/bennicholls/delvetown/console"

//ProgressBar is a textbox whose background can be a progress bar. Yes.
type ProgressBar struct {
	Textbox
	progresscolour uint32
	progress       int //percentage
}

func NewProgressBar(w, h, x, y, z int, bord, cent bool, txt string, c uint32) *ProgressBar {
	return &ProgressBar{*NewTextbox(w, h, x, y, z, bord, cent, txt), c, 0}
}

//Takes a percentage value, 0 <= i <= 100. Ignores value otherwise.
func (pb *ProgressBar) SetProgress(i int) {
	if i >= 0 && i <= 100 {
		pb.progress = i
	}
}

func (pb *ProgressBar) SetProgressColour(c uint32) {
	pb.progresscolour = c
}

func (pb ProgressBar) Render(offset ...int) {
	if pb.visible {
		offX, offY, offZ := processOffset(offset)

		pb.Textbox.Render(offX, offY, offZ)

		barWidth := int(float32(pb.progress) * float32(pb.width) / float32(100))
		if barWidth == 0 && pb.progress != 0 {
			barWidth = 1
		}

		for i := 0; i < barWidth; i++ {
			for j := 0; j < pb.height; j++ {
				console.ChangeBackColour(i+offX, j+offY, pb.progresscolour)
			}
		}
	}
}
