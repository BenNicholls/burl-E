package ui

import "github.com/bennicholls/burl/console"
import "github.com/bennicholls/burl/util"

//ProgressBar is a textbox whose background can be a progress bar. Yes.
type ProgressBar struct {
	Textbox
	progresscolour uint32
	progress       int //percentage
}

func NewProgressBar(w, h, x, y, z int, bord, cent bool, txt string, c uint32) *ProgressBar {
	return &ProgressBar{*NewTextbox(w, h, x, y, z, bord, cent, txt), c, 0}
}

//Takes a percentage value, clamped to 0 <= i <= 100
func (pb *ProgressBar) SetProgress(i int) {
	pb.progress = util.Clamp(i, 0, 100)
}

func (pb *ProgressBar) ChangeProgress(d int) {
	pb.SetProgress(pb.progress + d)
}

func (pb ProgressBar) GetProgress() int {
	return pb.progress
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
		} else if pb.progress == 100 {
			barWidth = pb.width
		}

		for i := 0; i < pb.width; i++ {
			for j := 0; j < pb.height; j++ {
				if i < barWidth {
					console.ChangeBackColour(i+offX+pb.x, j+offY+pb.y, pb.z+offZ, pb.progresscolour)
				} else {
					//set to black for now (bgcolor support coming later I assume)
					console.ChangeBackColour(i+offX+pb.x, j+offY+pb.y, pb.z+offZ, 0xFF000000)
				}
			}
		}
	}
}
