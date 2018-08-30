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

// FromVec converts a vec into a tile that contains that vec
func FromVec(v pixel.Vec, zoom int) Tile {
	return Tile{
		X: int(v.X) / TileSize,
		Y: (-int(v.Y)) / TileSize,
		Z: zoom,
	}
}

// FromLatLon take float lat/lons and a zoom and return a tile
// Clips the coordinates if they are outside of Min/MaxLat/Lon
func FromLatLon(lat, lon float64, zoom int) Tile {
	c := ClippedCoords(lat, lon)
	return FromVec(c.Vec(zoom), zoom)
}

// Vec returns a vector for the bottom left corner of the tile
func (t Tile) Vec() pixel.Vec {
	return pixel.V(
		float64(t.X*TileSize),
		-float64((t.Y*TileSize)+TileSize),
	)
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
		min = FromVec(r.Min, zoom)
		max = FromVec(r.Max, zoom)
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
