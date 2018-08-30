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
	Sprite *pixel.Sprite
	Offset pixel.Vec
	Frame  pixel.Rect
}

func NewImageTile(t tiles.Tile, bounds pixel.Rect) ImageTile {
	frame := t.Rect().Intersect(bounds)
	return ImageTile{
		Tile:   t,
		Sprite: pixel.NewSprite(Placeholder, frame),
		Offset: bounds.Min.To(frame.Min),
		Frame:  frame,
	}
}

func (t ImageTile) Fetch() error {
	pic, err := TilePictureData(t.Tile)
	if err != nil {
		return err
	}
	t.Sprite = pixel.NewSprite(pic, pic.Bounds())
	return nil
}

func (t ImageTile) Draw(tg pixel.Target, m pixel.Matrix) {
	half := float64(tiles.TileSize) / 2
	v := t.Vec().Add(pixel.V(half, half))
	t.Sprite.Draw(tg, pixel.IM.Moved(v).Chained(m))
}

// RectTiles returns a slice of tiles requires to fully cover the rect
func RectTiles(bounds pixel.Rect, zoom int) []ImageTile {
	var (
		min = tiles.FromVec(bounds.Min, zoom)
		max = tiles.FromVec(bounds.Max, zoom)
		tt  []ImageTile
	)
	for x := min.X; x <= max.X; x++ {
		for y := max.Y; y <= min.Y; y++ {
			t := tiles.Tile{
				X: x,
				Y: y,
				Z: zoom,
			}
			tt = append(tt, NewImageTile(t, bounds))
		}
	}
	return tt
}
