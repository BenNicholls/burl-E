package burl

//ProgressBar is a textbox whose background can be a progress bar. Yes.
type ProgressBar struct {
	Textbox
	progresscolour uint32
	progress       int //percentage
	barWidth       int
}

func NewProgressBar(w, h, x, y, z int, bord, cent bool, txt string, c uint32) *ProgressBar {
	return &ProgressBar{*NewTextbox(w, h, x, y, z, bord, cent, txt), c, 0, 0}
}

//Takes a percentage value, clamped to 0 <= i <= 100
func (pb *ProgressBar) SetProgress(i int) {
	pb.progress = Clamp(i, 0, 100)
	pb.CalcBarWidth()
}

func (pb *ProgressBar) ChangeProgress(d int) {
	pb.SetProgress(pb.progress + d)
	pb.CalcBarWidth()
}

func (pb *ProgressBar) CalcBarWidth() {
	pb.barWidth = int(float32(pb.progress) * float32(pb.width) / float32(100))
	if pb.barWidth == 0 && pb.progress != 0 {
		pb.barWidth = 1
	} else if pb.progress == 100 {
		pb.barWidth = pb.width
	}
}

func (pb *ProgressBar) GetProgress() int {
	return pb.progress
}

func (pb *ProgressBar) SetProgressColour(c uint32) {
	pb.progresscolour = c
}

func (pb *ProgressBar) Render() {
	if pb.visible {

		//need to set bgcolor to COL_NONE so we can draw text over the progressbar without borking it
		bg := pb.backColour
		pb.SetBackColour(COL_NONE)
		pb.Textbox.Render()
		pb.SetBackColour(bg)

		for i := 0; i < pb.width; i++ {
			for j := 0; j < pb.height; j++ {
				if i < pb.barWidth {
					console.ChangeBackColour(i+pb.x, j+pb.y, pb.z, pb.progresscolour)
				} else {
					console.ChangeBackColour(i+pb.x, j+pb.y, pb.z, pb.backColour)
				}
			}
		}
	}
}
