package burl

type TileMap struct {
	Width, Height int
	Tiles         []Tile
}

func NewMap(w, h int) *TileMap {
	return &TileMap{Width: w, Height: h, Tiles: make([]Tile, w*h)}
}

func (m TileMap) Dims() (int, int) {
	return m.Width, m.Height
}

func (m *TileMap) ChangeTileType(x, y, tile int) {
	if CheckBounds(x, y, m.Width, m.Height) {
		m.Tiles[y*m.Width+x].TileType = tile
	}
}

func (m TileMap) GetTileType(x, y int) int {
	if CheckBounds(x, y, m.Width, m.Height) {
		return m.Tiles[y*m.Width+x].TileType
	} else {
		return 0
	}
}

func (m TileMap) GetTile(x, y int) Tile {
	if CheckBounds(x, y, m.Width, m.Height) {
		return m.Tiles[y*m.Width+x]
	} else {
		return Tile{}
	}
}

func (m *TileMap) SetTile(x, y int, t Tile) {
	if CheckBounds(x, y, m.Width, m.Height) {
		m.Tiles[x+y*m.Width] = t
	}
}

func (m *TileMap) AddEntity(x, y int, e Entity) {
	if CheckBounds(x, y, m.Width, m.Height) {
		m.Tiles[x+y*m.Width].entity = e
		m.ShadowCast(x, y, e.GetLight().Strength, Lighten)
	}
}

func (m *TileMap) RemoveEntity(x, y int) {
	if CheckBounds(x, y, m.Width, m.Height) {
		m.ShadowCast(x, y, m.Tiles[x+y*m.Width].entity.GetLight().Strength, Darken)
		m.Tiles[x+y*m.Width].entity = nil
	}
}

func (m *TileMap) MoveEntity(x, y, dx, dy int) {
	e := m.Tiles[x+y*m.Width].entity
	if e != nil {
		m.RemoveEntity(x, y)
		m.AddEntity(x+dx, y+dy, e)

	}
}

func (m TileMap) GetEntity(x, y int) Entity {
	if CheckBounds(x, y, m.Width, m.Height) {
		return m.Tiles[x+y*m.Width].entity
	} else {
		return nil
	}
}

// func (m *TileMap) AddItem(x, y int, i *Item) {
// 	if CheckBounds(x, y, m.Width, m.Height) && i != nil {
// 		m.Tiles[x+y*m.Width].Item = i
// 	}
// }

// func (m *TileMap) RemoveItem(x, y int) {
// 	if CheckBounds(x, y, m.Width, m.Height) {
// 		m.Tiles[x+y*m.Width].Item = nil
// 	}
// }

// func (m *TileMap) GetItem(x, y int) *Item {
// 	if CheckBounds(x, y, m.Width, m.Height) {
// 		return m.Tiles[x+y*m.Width].Item
// 	} else {
// 		return nil
// 	}
// }

//For testing purposes.
func (m *TileMap) ChangeTileColour(x, y int, c uint32) {
	if CheckBounds(x, y, m.Width, m.Height) {
		m.Tiles[x+y*m.Width].Light.Colour = c
	}
}

func (m TileMap) LastVisible(x, y int) int {
	if CheckBounds(x, y, m.Width, m.Height) {
		return m.Tiles[x+y*m.Width].LastVisible
	} else {
		return 0
	}
}

//NOTE: Consider renaming this.
func (m *TileMap) SetVisible(x, y, tick int) {
	if CheckBounds(x, y, m.Width, m.Height) {
		m.Tiles[x+y*m.Width].LastVisible = tick
	}
}

func (m *TileMap) ClearLights() {
	for i, _ := range m.Tiles {
		m.Tiles[i].Light.Bright = 0
	}
}

//Basic unit for the world. Holds a type (grass, wall, etc), a list of contained items (dropped weapons),
//and a pointer to an Entity if one is standing there. Eventually will hold pathfinding information too.
type Tile struct {
	TileType, Variant int //
	entity            Entity
	LastVisible       int // Records the last tick that this tile was seen
	Light             TileLight
	//Item              *Item
}

func (t Tile) Passable() bool {
	return IsPassable(t.TileType) && t.entity == nil
}

func (t Tile) Transparent() bool {
	return IsTransparent(t.TileType)
}

func (t Tile) Empty() bool {
	//return t.Entity == nil && t.Item == nil && IsPassable(t.tileType)
	return t.entity == nil && IsPassable(t.TileType)
}

func (t Tile) GetVisuals() Visuals {
	return tiledata[t.TileType].vis
}

//Light characteristics for each tile.
type TileLight struct {
	Colour uint32
	Bright int //Brightness level 0-255
}
