package burl

type TileMap struct {
	width, height int
	tiles         []Tile
}

func NewMap(w, h int) *TileMap {
	return &TileMap{width: w, height: h, tiles: make([]Tile, w*h)}
}

func (m TileMap) Dims() (int, int) {
	return m.width, m.height
}

func (m *TileMap) ChangeTileType(x, y, tile int) {
	if CheckBounds(x, y, m.width, m.height) {
		m.tiles[y*m.width+x].tileType = tile
	}
}

func (m TileMap) GetTileType(x, y int) int {
	if CheckBounds(x, y, m.width, m.height) {
		return m.tiles[y*m.width+x].tileType
	} else {
		return 0
	}
}

func (m TileMap) GetTile(x, y int) Tile {
	if CheckBounds(x, y, m.width, m.height) {
		return m.tiles[y*m.width+x]
	} else {
		return Tile{}
	}
}

func (m *TileMap) SetTile(x, y int, t Tile) {
	if CheckBounds(x, y, m.width, m.height) {
		m.tiles[x+y*m.width] = t
	}
}

func (m *TileMap) AddEntity(x, y int, e Entity) {
	if CheckBounds(x, y, m.width, m.height) {
		m.tiles[x+y*m.width].entity = e
		m.ShadowCast(x, y, e.GetLight().Strength, Lighten)
	}
}

func (m *TileMap) RemoveEntity(x, y int) {
	if CheckBounds(x, y, m.width, m.height) {
		m.ShadowCast(x, y, m.tiles[x+y*m.width].entity.GetLight().Strength, Darken)
		m.tiles[x+y*m.width].entity = nil
	}
}

func (m *TileMap) MoveEntity(x, y, dx, dy int) {
	e := m.tiles[x+y*m.width].entity
	if e != nil {
		m.RemoveEntity(x, y)
		m.AddEntity(x+dx, y+dy, e)

	}
}

func (m TileMap) GetEntity(x, y int) Entity {
	if CheckBounds(x, y, m.width, m.height) {
		return m.tiles[x+y*m.width].entity
	} else {
		return nil
	}
}

// func (m *TileMap) AddItem(x, y int, i *Item) {
// 	if CheckBounds(x, y, m.width, m.height) && i != nil {
// 		m.tiles[x+y*m.width].Item = i
// 	}
// }

// func (m *TileMap) RemoveItem(x, y int) {
// 	if CheckBounds(x, y, m.width, m.height) {
// 		m.tiles[x+y*m.width].Item = nil
// 	}
// }

// func (m *TileMap) GetItem(x, y int) *Item {
// 	if CheckBounds(x, y, m.width, m.height) {
// 		return m.tiles[x+y*m.width].Item
// 	} else {
// 		return nil
// 	}
// }

//For testing purposes.
func (m *TileMap) ChangeTileColour(x, y int, c uint32) {
	if CheckBounds(x, y, m.width, m.height) {
		m.tiles[x+y*m.width].Light.Colour = c
	}
}

func (m TileMap) LastVisible(x, y int) int {
	if CheckBounds(x, y, m.width, m.height) {
		return m.tiles[x+y*m.width].LastVisible
	} else {
		return 0
	}
}

//NOTE: Consider renaming this.
func (m *TileMap) SetVisible(x, y, tick int) {
	if CheckBounds(x, y, m.width, m.height) {
		m.tiles[x+y*m.width].LastVisible = tick
	}
}

func (m *TileMap) ClearLights() {
	for i, _ := range m.tiles {
		m.tiles[i].Light.Bright = 0
	}
}

//Basic unit for the world. Holds a type (grass, wall, etc), a list of contained items
//(dropped weapons), and a pointer to an Entity if one is standing there
//Eventually will hold pathfinding information too.
type Tile struct {
	tileType, variant int //
	passable          bool
	entity            Entity
	LastVisible       int // Records the last tick that this tile was seen
	Light             TileLight
	//Item              *Item
}

func (t Tile) TileType() int {
	return t.tileType
}

func (t Tile) Passable() bool {
	return IsPassable(t.tileType) && t.entity == nil
}

func (t Tile) Transparent() bool {
	return IsTransparent(t.tileType)
}

func (t Tile) Empty() bool {
	//return t.Entity == nil && t.Item == nil && IsPassable(t.tileType)
	return t.entity == nil && IsPassable(t.tileType)
}

func (t Tile) GetVisuals() Visuals {
	return tiledata[t.tileType].vis
}

//Light characteristics for each tile.
type TileLight struct {
	Colour uint32
	Bright int //Brightness level 0-255
}
