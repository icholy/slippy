package slippy

import (
	"image"

	"github.com/faiface/pixel"

	"github.com/icholy/slippy/tiles"
)

var (
	Placeholder = pixel.PictureDataFromImage(
		image.NewRGBA(image.Rect(0, 0, tiles.TileSize, tiles.TileSize)),
	)
)

type ImageTile struct {
	tiles.Tile
	Pic    *pixel.PictureData
	Offset pixel.Vec
}

// RectTiles returns a slice of tiles requires to fully cover the rect
func RectTiles(r pixel.Rect, zoom int) []ImageTile {
	var (
		min = tiles.FromVec(r.Min, zoom)
		max = tiles.FromVec(r.Max, zoom)
		tt  []ImageTile
	)
	for x := min.X; x <= max.X; x++ {
		for y := max.Y; y <= min.Y; y++ {
			t := tiles.Tile{
				X: x,
				Y: y,
				Z: zoom,
			}
			tt = append(tt, ImageTile{
				Tile: t,
				Pic:  Placeholder,
			})
		}
	}
	return tt
}
