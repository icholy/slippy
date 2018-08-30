package convert

import (
	"github.com/faiface/pixel"

	"github.com/icholy/slippy/tiles"
)

// CoordinateVec converst a coordinate into a vec
func CoordinateVec(c tiles.Coordinate, zoom int) pixel.Vec {
	return PixelVec(c.Pixel(zoom))
}

// VecCoordinate converts a vec into a coordinate
func VecCoordinate(v pixel.Vec, zoom int) tiles.Coordinate {
	return VecPixel(v, zoom).Coords()
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
	t, _ := VecPixel(v, zoom).Tile()
	return t
}

// TileVec returns a vector for the bottom left corner of the tile
func TileVec(t tiles.Tile) pixel.Vec {
	p := t.Pixel()
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

// RectTiles returns a slice of tiles requires to fully cover the rect
func RectTiles(r pixel.Rect, zoom int) []tiles.Tile {
	var (
		min = VecTile(r.Min, zoom)
		max = VecTile(r.Max, zoom)
		tt  []tiles.Tile
	)
	for x := min.X; x <= max.X; x++ {
		for y := max.Y; y <= min.Y; y++ {
			tt = append(tt, tiles.Tile{
				X: x,
				Y: y,
				Z: zoom,
			})
		}
	}
	return tt
}
