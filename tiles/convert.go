package tiles

import (
	"github.com/faiface/pixel"
)

// Vec converts a coordinate into a vec
func (c Coordinate) Vec(zoom int) pixel.Vec {
	return c.Pixel(zoom).Vec()
}

// VecCoordinate converts a vec into a coordinate
func VecCoordinate(v pixel.Vec, zoom int) Coordinate {
	return VecPixel(v, zoom).Coords()
}

// VecPixel converts a vec into a pixel
func VecPixel(v pixel.Vec, zoom int) Pixel {
	return Pixel{
		X: int(v.X),
		Y: -int(v.Y),
		Z: zoom,
	}
}

// Vec converts a pixel into a vec
func (p Pixel) Vec() pixel.Vec {
	return pixel.V(
		float64(p.X),
		-float64(p.Y),
	)
}

// VecTile converts a vec into a tile that contains that vec
func VecTile(v pixel.Vec, zoom int) Tile {
	return VecPixel(v, zoom).Tile()
}

// Vec returns a vector for the bottom left corner of the tile
func (t Tile) Vec() pixel.Vec {
	p := t.Pixel()
	return Pixel{
		X: p.X,
		Y: p.Y + TileSize,
	}.Vec()
}

// Rect returns a rectangle of the tile
func (t Tile) Rect() pixel.Rect {
	v := t.Vec()
	return pixel.R(
		v.X,
		v.Y,
		v.X+float64(TileSize),
		v.Y+float64(TileSize),
	)
}

// RectTiles returns a slice of tiles requires to fully cover the rect
func RectTiles(r pixel.Rect, zoom int) []Tile {
	var (
		min = VecTile(r.Min, zoom)
		max = VecTile(r.Max, zoom)
		tt  []Tile
	)
	for x := min.X; x <= max.X; x++ {
		for y := max.Y; y <= min.Y; y++ {
			tt = append(tt, Tile{
				X: x,
				Y: y,
				Z: zoom,
			})
		}
	}
	return tt
}
