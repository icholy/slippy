// Package tiles is a collection of conversion utilities to go between geo/pixel/tile/quadkey space
// This package uses WGS84 coordinates and a mercator projection
// There is also a TileIndex which can be used to store data in a single place and aggregate when needed
package tiles

// Tile is a simple struct for holding the XYZ coordinates for use in mapping
type Tile struct {
	X, Y, Z int
}

// Pixel return the NW pixel of this tile
func (t Tile) Pixel() Pixel {
	return Pixel{
		X: t.X * TileSize,
		Y: t.Y * TileSize,
		Z: t.Z,
	}
}

// ToPixelWithOffset returns a pixel at the origin with an offset added. Useful for getting the center pixel of a tile or another non-origin pixel.
func (t Tile) ToPixelWithOffset(offset Pixel) (pixel Pixel) {
	pixel = t.Pixel()
	pixel.X += offset.X
	pixel.Y += offset.Y
	return
}

// FromCoordinate take float lat/lons and a zoom and return a tile
// Clips the coordinates if they are outside of Min/MaxLat/Lon
func FromCoordinate(lat, lon float64, zoom int) Tile {
	c := ClippedCoords(lat, lon)
	p := c.Pixel(zoom)
	t, _ := p.Tile()
	return t
}
