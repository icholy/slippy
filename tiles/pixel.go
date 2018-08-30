package tiles

import (
	"math"
)

// Pixel in a WGS84 Mercator map projection with a NW origin (0,0) of the projection
type Pixel struct {
	X, Y, Z int
}

func (p Pixel) floatX() float64 {
	return float64(p.X)
}

func (p Pixel) floatY() float64 {
	return float64(p.Y)
}

// Coords converts to WGS84 coordaintes
func (p Pixel) Coords() Coordinate {
	size := float64(mapDimensions(p.Z))
	x := (clip(p.floatX(), 0, size-1) / size) - 0.5
	y := 0.5 - (clip(p.floatY(), 0, size-1) / size)
	lat := 90 - 360*math.Atan(math.Exp(-y*2*math.Pi))/math.Pi
	lon := 360.0 * x
	return ClippedCoords(lat, lon)
}

// Tile gets the tile that contains this pixel
func (p Pixel) Tile() Tile {
	return Tile{
		X: p.X / TileSize,
		Y: p.Y / TileSize,
		Z: p.Z,
	}
}
