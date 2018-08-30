package slippy

import (
	"fmt"
	"image"
	"math/rand"

	"github.com/faiface/pixel"

	"github.com/icholy/slippy/tiles"
	"github.com/icholy/slippy/util"
)

var (
	Placeholder = pixel.PictureDataFromImage(
		image.NewRGBA(image.Rect(0, 0, tiles.TileSize, tiles.TileSize)),
	)
)

func TilePictureData(t tiles.Tile) (*pixel.PictureData, error) {
	url := URL(t)
	img, err := util.FetchImage(url)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

type ImageTile struct {
	tiles.Tile
	Sprite *pixel.Sprite
	Frame  pixel.Rect
	Offset pixel.Vec
}

func URL(t tiles.Tile) string {
	shards := []string{"a", "b", "c"}
	return fmt.Sprintf(
		"http://%[1]s.tile.openstreetmap.org/%[2]d/%[3]d/%[4]d.png",
		shards[rand.Intn(len(shards))], t.Z, t.X, t.Y,
	)
}

func (t *ImageTile) Fetch() error {
	pic, err := TilePictureData(t.Tile)
	if err != nil {
		return err
	}
	t.Sprite = pixel.NewSprite(pic, pic.Bounds())
	return nil
}

func (t ImageTile) Draw(tg pixel.Target) {
	if t.Sprite == nil {
		return
	}
	m := float64(tiles.TileSize) / 2
	v := t.Vec().Add(pixel.V(m, m))
	t.Sprite.Draw(tg, pixel.IM.Moved(v))
}

func NewImageTile(t tiles.Tile, bounds pixel.Rect) ImageTile {
	var (
		pic     = Placeholder
		rect    = t.Rect()
		overlap = rect.Intersect(bounds)
		frame   = overlap.Moved(pixel.ZV.Sub(rect.Min))
	)

	return ImageTile{
		Tile:   t,
		Sprite: pixel.NewSprite(pic, pic.Bounds()),
		Frame:  frame,
	}
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
