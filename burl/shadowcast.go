package burl

//Shadowcatser runs 8 times over different quadrants. rMatrix supplies rotation coefficients.
//Linear algebra to the rescue.
var rMatrix = [8][4]int{{1, 0, 0, 1}, {-1, 0, 0, 1}, {0, -1, 1, 0}, {0, -1, -1, 0}, {-1, 0, 0, -1}, {1, 0, 0, -1}, {0, 1, -1, 0}, {0, 1, 1, 0}}

//THE BIG CHEESE - The one and only Shadowcaster! For all of your FOV needs. fn is a function
//for the shadowcaster to apply to open spaces it finds.
func (m *TileMap) ShadowCast(x, y, radius int, fn Cast) {
	if radius <= 0 {
		return
	}
	fn(m, x, y, 0, radius)
	for i := 0; i < 8; i++ {
		m.scan(x, y, 1, 1.0, 0.0, radius, rMatrix[i], i%2 == 0, fn)
	}
}

//TODO: General cleanup. Direct port from python, not exactly golangish.
//NOTE: The 'cull' bool controls the logic for ensuring the 8 passes don't overlap at the edges.
//It is set to true for the odd-numbered scans. The shadowcaster still visits these squares twice,
//but the function fn is not run twice. Trust me Ben, this was the best way you could think of and
//your other solutions created crazy behaviour. Leave it alone!
func (m *TileMap) scan(x, y, row int, slope1, slope2 float32, radius int, r [4]int, cull bool, fn Cast) {
	if slope1 < slope2 {
		return
	}
	blocked := false

	//scan #radius rows
	for j := row; j < radius+1 && !blocked; j++ {

		//scan row
		for dx, dy, newStart := -j, -j, slope1; dx <= 0; dx++ {
			mx, my := x+dx*r[0]+dy*r[1], y+dx*r[2]+dy*r[3] //map coordinates
			if !CheckBounds(mx, my, m.Width, m.Height) {
				continue
			}
			lSlope, rSlope := (float32(dx)-0.5)/(float32(dy)+0.5), (float32(dx)+0.5)/(float32(dy)-0.5)

			if newStart < rSlope {
				continue
			} else if slope2 > lSlope {
				break
			} else {
				if d := Distance(0, 0, dx, dy); d < radius*radius {
					if !cull || !(dx == 0 || dy == 0 || dx == dy) {
						fn(m, mx, my, d, radius)
					}
				}
				//scanning a block
				if blocked {
					if m.Tiles[mx+my*m.Width].Transparent() {
						blocked = false
						slope1 = newStart
					} else {
						newStart = rSlope
					}
				} else {
					//blocked square, commence child scan
					if !m.Tiles[mx+my*m.Width].Transparent() && j < radius {
						blocked = true
						m.scan(x, y, j+1, newStart, lSlope, radius, r, cull, fn)
						newStart = rSlope
					}
				}
			}
		}
	}
}

//type specifying precisely what you can pass to the shadowcaster. parameters here are the info
//that the shadowcaster will deliver
type Cast func(m *TileMap, x, y, d, r int)

//Run this over a tilemap to light squares. Linearly interpolates from max (255) at center to 0 at r
func Lighten(m *TileMap, x, y, d, r int) {
	m.Tiles[x+y*m.Width].Light.Bright += (255 - int(255*float32(d)/float32(r*r)))
}

//Same as above, but opposite.
func Darken(m *TileMap, x, y, d, r int) {
	m.Tiles[x+y*m.Width].Light.Bright -= (255 - int(255*float32(d)/float32(r*r)))

	if m.Tiles[x+y*m.Width].Light.Bright < 0 {
		m.Tiles[x+y*m.Width].Light.Bright = 0
	}
}

//gets a list of the position of all empty tiles
//TODO: would this be better returning a list of *Tile?
func GetEmptySpacesCast(spaces *[]Coord) Cast {
	return func(m *TileMap, x, y, d, r int) {
		if m.GetTile(x, y).Empty() {
			*spaces = append(*spaces, Coord{x, y})
		}
	}
}
