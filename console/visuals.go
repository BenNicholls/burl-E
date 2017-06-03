package console

type Visuals struct {
	Glyph      int
	ForeColour uint32
	BackColour uint32
}

func (v *Visuals) ChangeGlyph(g int) {
	v.Glyph = g
}

func (v *Visuals) ChangeForeColour(f uint32) {
	v.ForeColour = f
}

func (v *Visuals) ChangeBackColour(b uint32) {
	v.BackColour = b
}
