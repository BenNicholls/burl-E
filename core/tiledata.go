package core

var tiledata []tileTypeData

//default tiletypes POSSIBLE TODO: dynamic changing tile properties? Think about this.
const (
	TILE_NOTHING = iota
	MAX_TILETYPES
)

type tileTypeData struct {
	name        string
	passable    bool
	transparent bool
	vis         Visuals
}

type Visuals struct {
	Glyph      int
	ForeColour uint32
}

func init() {

	//tiledata[TILETYPE]
	tiledata = make([]tileTypeData, MAX_TILETYPES)

	//tiledata definitions go here. TODO: some kind of data loading function, load from file.
	tiledata[TILE_NOTHING] = tileTypeData{"Nothing", false, true, Visuals{0, 0x000000}}
}

//takes tiletype, returns glyph
func GetName(t int) string {
	return tiledata[t].name
}

func IsPassable(t int) bool {
	return tiledata[t].passable
}

func IsTransparent(t int) bool {
	return tiledata[t].transparent
}

func GetTileVisuals(t int) Visuals {
	return tiledata[t].vis
}
