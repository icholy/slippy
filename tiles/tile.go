// Package tiles is a collection of conversion utilities to go between geo/pixel/tile/quadkey space
// This package uses WGS84 coordinates and a mercator projection
// There is also a TileIndex which can be used to store data in a single place and aggregate when needed
package tiles

import (
	"github.com/faiface/pixel"
)

// Tile is a simple struct for holding the XYZ coordinates for use in mapping
type Tile struct {
	X, Y, Z int
}

// pixel return the NW pixel of this tile
func (t Tile) pixel() Pixel {
	return Pixel{
		X: t.X * TileSize,
		Y: t.Y * TileSize,
		Z: t.Z,
	}
}

// ToPixelWithOffset returns a pixel at the origin with an offset added. Useful for getting the center pixel of a tile or another non-origin pixel.
func (t Tile) ToPixelWithOffset(offset Pixel) (pixel Pixel) {
	pixel = t.pixel()
	pixel.X += offset.X
	pixel.Y += offset.Y
	return
}

// FromCoordinate take float lat/lons and a zoom and return a tile
// Clips the coordinates if they are outside of Min/MaxLat/Lon
func FromCoordinate(lat, lon float64, zoom int) Tile {
	c := ClippedCoords(lat, lon)
	p := c.pixel(zoom)
	return p.Tile()
}

// Vec returns a vector for the bottom left corner of the tile
func (t Tile) Vec() pixel.Vec {
	p := t.pixel()
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

// VecTile converts a vec into a tile that contains that vec
func VecTile(v pixel.Vec, zoom int) Tile {
	return VecPixel(v, zoom).Tile()
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
