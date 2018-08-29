package slippy

import (
	"github.com/buckhx/tiles"
	"github.com/faiface/pixel"
)

// CoordinateVec converst a coordinate into a vec
func CoordinateVec(c tiles.Coordinate, zoom int) pixel.Vec {
	return PixelVec(c.ToPixel(zoom))
}

// VecCoordinate converts a vec into a coordinate
func VecCoordinate(v pixel.Vec, zoom int) tiles.Coordinate {
	return VecPixel(v, zoom).ToCoords()
}

// VecPixel converts a vec into a pixel
func VecPixel(v pixel.Vec, zoom int) tiles.Pixel {
	return tiles.Pixel{
		X: int(v.X),
		Y: -int(v.Y),
		Z: zoom,
	}
}

// PixelVec converts a pixel into a vec
func PixelVec(p tiles.Pixel) pixel.Vec {
	return pixel.V(
		float64(p.X),
		-float64(p.Y),
	)
}

// VecTile converts a vec into a tile that contains that vec
func VecTile(v pixel.Vec, zoom int) tiles.Tile {
	t, _ := VecPixel(v, zoom).ToTile()
	return t
}

// TileVec returns a vector for the bottom left corner of the tile
func TileVec(t tiles.Tile) pixel.Vec {
	p := t.ToPixel()
	return PixelVec(tiles.Pixel{
		X: p.X,
		Y: p.Y + tiles.TileSize,
	})
}

// TileRect returns a rectangle of the tile
func TileRect(t tiles.Tile) pixel.Rect {
	v := TileVec(t)
	return pixel.R(
		v.X,
		v.Y,
		v.X+float64(tiles.TileSize),
		v.Y+float64(tiles.TileSize),
	)
}

func RecTiles(r)
