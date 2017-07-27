package core

var tiledata []tileTypeData

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

//Inits the tile data repository, which for now is just a slice of datas.
//Also loads a NOTHING entry,
func init() {
	//tiledata[TILETYPE]
	tiledata = make([]tileTypeData, 1)
	LoadTileData("Nothing", false, true, 0, 0xFF000000)
}

//Adds a new entry to the tile data respoitory. Returns the index for the data in the repo.
//TODO: load from file.
func LoadTileData(name string, pass, trans bool, glyph int, c uint32) int {
	tiledata = append(tiledata, tileTypeData{name, pass, trans, Visuals{glyph, c}})
	return len(tiledata) - 1
}

func GetName(t int) string {
	if t < len(tiledata) {
		return tiledata[t].name
	} else {
		return "no tile"
	}
}

func IsPassable(t int) bool {
	if t < len(tiledata) {
		return tiledata[t].passable
	} else {
		return false
	}
}

func IsTransparent(t int) bool {
	if t < len(tiledata) {
		return tiledata[t].transparent
	} else {
		return false
	}
}

func GetTileVisuals(t int) Visuals {
	if t < len(tiledata) {
		return tiledata[t].vis
	} else {
		return Visuals{0, 0xFF000000}
	}
}
