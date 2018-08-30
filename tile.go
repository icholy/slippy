// Package tiles is a collection of conversion utilities to go between geo/pixel/tile/quadkey space
// This package uses WGS84 coordinates and a mercator projection
// There is also a TileIndex which can be used to store data in a single place and aggregate when needed
package slippy

import (
	"fmt"

	"github.com/faiface/pixel"
)

// Tile is a simple struct for holding the XYZ coordinates for use in mapping
type Tile struct {
	X, Y, Z int
}

// String returns the string representation of the tile
func (t Tile) String() string {
	return fmt.Sprintf("Tile(%d, %d, %d)", t.X, t.Y, t.Z)
}

// fromVec converts a vec into a tile that contains that vec
func fromVec(v pixel.Vec, zoom int) Tile {
	return Tile{
		X: int(v.X) / TileSize,
		Y: (-int(v.Y)) / TileSize,
		Z: zoom,
	}
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
